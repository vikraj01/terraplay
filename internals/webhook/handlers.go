package webhook

import (
	"archive/zip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

type WorkflowRunPayload struct {
	Action      string `json:"action"` // The action for the workflow run event, e.g., "requested", "completed"
	WorkflowRun struct {
		ID         int64  `json:"id"`
		NodeID     string `json:"node_id"`
		Status     string `json:"status"`
		Conclusion string `json:"conclusion"`
		WorkflowID int64  `json:"workflow_id"`
		RunNumber  int64  `json:"run_number"`
		Event      string `json:"event"`
		CreatedAt  string `json:"created_at"`
		UpdatedAt  string `json:"updated_at"`
		URL        string `json:"url"`
		HTMLURL    string `json:"html_url"`
		JobsURL    string `json:"logs_url"` // This is needed to fetch the job logs
		Outputs    struct {
			ServerIP string `json:"server_ip"`
		} `json:"outputs"`
		Inputs struct {
			RunID  string `json:"run_id"`
			UserID string `json:"user_id"`
			Game   string `json:"game"`
		} `json:"inputs"`
	} `json:"workflow_run"`
	Repository struct {
		ID       int64  `json:"id"`
		NodeID   string `json:"node_id"`
		Name     string `json:"name"`
		FullName string `json:"full_name"`
	} `json:"repository"`
	Sender struct {
		Login string `json:"login"`
		ID    int64  `json:"id"`
	} `json:"sender"`
}

func HandleWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		http.Error(w, "Can't read body", http.StatusBadRequest)
		return
	}

	requestID := uuid.New().String()
	timestamp := time.Now().Format("20060102_150405")
	folder := "webhooks_logs"

	err = os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		log.Printf("Failed to create folder for logs: %v", err)
	}

	// rawBodyFile := filepath.Join(folder, fmt.Sprintf("webhook_raw_body_%s_%s.json", timestamp, requestID))

	// err = ioutil.WriteFile(rawBodyFile, body, 0644)
	// if err != nil {
	//     log.Printf("Failed to write raw body to file: %v", err)
	// } else {
	//     log.Printf("Raw body written to: %s", rawBodyFile)
	// }

	// VERIFYING THE SIGNATURE
	signature := r.Header.Get("X-Hub-Signature-256")
	if signature == "" {
		http.Error(w, "Missing signature", http.StatusUnauthorized)
		return
	}

	secret := os.Getenv("GITHUB_WEBHOOK_SECRET")
	if secret == "" {
		log.Println("GITHUB_WEBHOOK_SECRET is not set")
		http.Error(w, "Secret not set on the server", http.StatusInternalServerError)
		return
	}

	if !verifySignature(signature, body, secret) {
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	githubEvent := r.Header.Get("X-GitHub-Event")
	fmt.Println("GITHUB EVENT", githubEvent)

	// Conditional Saving State Based On The Github Event
	switch githubEvent {
	case "workflow_dispatch":
		// handleWorkflowDispatch(body, folder, timestamp, requestID)
	case "workflow_run":
		handleWorkflowRun(body, folder, timestamp, requestID)
	case "workflow_job":
		// HandleWorkflowJob(body, folder, timestamp, requestID)
	}

	var payload WorkflowRunPayload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		log.Printf("Failed to parse webhook payload: %v", err)
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// parsedPayloadFile := filepath.Join(folder, fmt.Sprintf("webhook_parsed_payload_%s_%s.json", timestamp, requestID))

	// payloadJSON, err := json.MarshalIndent(payload, "", "  ")
	// if err != nil {
	//     log.Printf("Failed to marshal JSON payload: %v", err)
	// }

	// err = ioutil.WriteFile(parsedPayloadFile, payloadJSON, 0644)
	// if err != nil {
	//     log.Printf("Failed to write JSON payload to file: %v", err)
	// } else {
	//     log.Printf("Parsed payload written to: %s", parsedPayloadFile)
	// }

	if payload.WorkflowRun.Status == "completed" {
		sendToDiscord(payload.WorkflowRun.Inputs.UserID, payload.WorkflowRun.Inputs.Game, payload.WorkflowRun.Conclusion, payload.WorkflowRun.Outputs.ServerIP, payload.WorkflowRun.Inputs.RunID)
	}

	fmt.Fprint(w, "Webhook received and processed")
}

func verifySignature(signature string, body []byte, secret string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	_, err := mac.Write(body)
	if err != nil {
		log.Println("Error writing body to HMAC:", err)
		return false
	}

	expectedSignature := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}

func sendToDiscord(userID, game, status, serverIP, runID string) {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		log.Fatalf("DISCORD_BOT_TOKEN is not set")
	}

	channelID := os.Getenv("DISCORD_CHANNEL_ID")
	if channelID == "" {
		log.Fatalf("DISCORD_CHANNEL_ID is not set")
	}

	log.Printf("UserID: %s, Game: %s, RunID: %s, Status: %s, ServerIP: %s", userID, game, runID, status, serverIP)

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Printf("Failed to create Discord session: %v", err)
		return
	}
	defer dg.Close()

	if userID == "" {
		log.Printf("UserID is empty, skipping user mention")
	}

	message := fmt.Sprintf("<@%s>, your game deployment for '%s' (Run ID: %s) completed with status: %s", userID, game, runID, status)

	if serverIP != "" {
		message += fmt.Sprintf("\nServer IP: %s", serverIP)
	}

	log.Printf("Sending message to Discord channel %s: %s", channelID, message)

	_, err = dg.ChannelMessageSend(channelID, message)
	if err != nil {
		log.Printf("Failed to send message to Discord: %v", err)
	}
}

func handleWorkflowDispatch(body []byte, folder, timestamp, requestID string) {
	type WorkflowDispatchPayload struct {
		Ref        string            `json:"ref"`
		Workflow   string            `json:"workflow"` // Workflow filename or ID
		Inputs     map[string]string `json:"inputs"`   // Inputs passed to the workflow
		Repository struct {
			Name     string `json:"name"`
			FullName string `json:"full_name"`
			NodeID   string `json:"node_id"`
		} `json:"repository"`
		Sender struct {
			Login  string `json:"login"`
			NodeID string `json:"node_id"`
		} `json:"sender"`
	}
	var payload WorkflowDispatchPayload
	err := json.Unmarshal(body, &payload)

	if err != nil {
		log.Printf("Failed to parse workflow_dispatch payload: %v", err)
		return
	}

	// Perform any additional actions like updating DynamoDB, if necessary
	// Example: updateDynamoDB(payload.Repository.NodeID, payload.Inputs)

	// Save the parsed payload to a file for debugging
	parsedPayloadFile := filepath.Join(folder, fmt.Sprintf("raw_%s_%s.json", timestamp, requestID))
	saveToFile(parsedPayloadFile, body)
	payloadBytes, err := json.MarshalIndent(payload, "", "  ") // Pretty print with indent
	if err != nil {
		log.Printf("Failed to marshal payload: %v", err)
		return
	}

	parsedPayloadFile = filepath.Join(folder, fmt.Sprintf("webhook_parsed_dispatch_%s_%s.json", timestamp, requestID))
	saveToFile(parsedPayloadFile, payloadBytes)
}

func saveToFile(filename string, data []byte) {
	err := ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		log.Printf("Failed to write parsed payload to file: %v", err)
	} else {
		log.Printf("Parsed payload written to: %s", filename)
	}
}

func handleWorkflowRun(body []byte, folder, timestamp, requestID string) {
	rawBodyFile := filepath.Join(folder, fmt.Sprintf("workflow_run_raw_%s_%s.json", timestamp, requestID))
	err := ioutil.WriteFile(rawBodyFile, body, 0644)
	if err != nil {
		log.Printf("Failed to write workflow_run raw body to file: %v", err)
	} else {
		log.Printf("workflow_run raw body written to: %s", rawBodyFile)
	}

	// Parse the workflow_run payload
	var payload WorkflowRunPayload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		log.Printf("Failed to parse workflow_run payload: %v", err)
		return
	}

	// Log the parsed payload to a file
	// parsedPayloadFile := filepath.Join(folder, fmt.Sprintf("workflow_run_parsed_%s_%s.json", timestamp, requestID))
	// payloadJSON, err := json.MarshalIndent(payload, "", "  ")
	// if err != nil {
	// 	log.Printf("Failed to marshal workflow_run payload: %v", err)
	// }
	// err = ioutil.WriteFile(parsedPayloadFile, payloadJSON, 0644)
	// if err != nil {
	// 	log.Printf("Failed to write workflow_run parsed payload to file: %v", err)
	// } else {
	// 	log.Printf("workflow_run parsed payload written to: %s", parsedPayloadFile)
	// }

	if payload.WorkflowRun.Status == "completed" {
		log.Printf("Fetching logs for workflow run: %d", payload.WorkflowRun.ID)
		err := fetchJobLogs(payload.WorkflowRun.JobsURL, folder, timestamp, requestID)
		if err != nil {
			log.Printf("Failed to fetch logs for workflow run: %v", err)
		}

		// Further processing: Update DynamoDB, send messages to Discord, etc.
		// ...
	}
}
// func fetchJobLogs(logsURL, folder, timestamp, requestID string) error {
// 	// Make API request to download logs
// 	client := &http.Client{}
// 	req, err := http.NewRequest("GET", logsURL, nil)
// 	if err != nil {
// 		return fmt.Errorf("failed to create request: %v", err)
// 	}

// 	// Add authorization header
// 	req.Header.Add("Authorization", "Bearer "+os.Getenv("GITHUB_TOKEN"))

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return fmt.Errorf("failed to fetch logs: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return fmt.Errorf("failed to fetch logs: received status %d", resp.StatusCode)
// 	}

// 	// Read logs response
// 	logsData, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return fmt.Errorf("failed to read logs: %v", err)
// 	}

// 	// Write logs to a file
// 	logsFile := filepath.Join(folder, fmt.Sprintf("logs_%s_%s.txt", timestamp, requestID))
// 	err = ioutil.WriteFile(logsFile, logsData, 0644)
// 	if err != nil {
// 		return fmt.Errorf("failed to write logs to file: %v", err)
// 	}

// 	log.Printf("Logs saved to: %s", logsFile)
// 	return nil
// }

type WorkflowJobPayload struct {
	Action      string `json:"action"`
	WorkflowJob struct {
		ID              int       `json:"id"`
		NodeID          string    `json:"node_id"`
		RunID           int       `json:"run_id"`
		RunURL          string    `json:"run_url"`
		Status          string    `json:"status"`
		Conclusion      string    `json:"conclusion"`
		StartedAt       time.Time `json:"started_at"`
		CompletedAt     time.Time `json:"completed_at"`
		Name            string    `json:"name"`
		RunnerID        int       `json:"runner_id"`
		RunnerName      string    `json:"runner_name"`
		RunnerGroupID   int       `json:"runner_group_id"`
		RunnerGroupName string    `json:"runner_group_name"`
	} `json:"workflow_job"`
	Repository struct {
		ID       int    `json:"id"`
		NodeID   string `json:"node_id"`
		Name     string `json:"name"`
		FullName string `json:"full_name"`
	} `json:"repository"`
}

func HandleWorkflowJob(body []byte, folder, timestamp, requestID string) {
	var payload WorkflowJobPayload
	err := json.Unmarshal(body, &payload)
	if err != nil {
		log.Printf("Failed to parse workflow_job payload: %v", err)
		return
	}

	// Writing the parsed payload to a file for debugging
	parsedPayloadFile := filepath.Join(folder, fmt.Sprintf("webhook_parsed_job_%s_%s.json", timestamp, requestID))
	// err = os.WriteFile(parsedPayloadFile, body)
	// if err != nil {
	// 	log.Printf("Failed to write parsed job payload to file: %v", err)
	// 	return
	// }
	// saveToFile(parsedPayloadFile, )

	log.Printf("Parsed payload written to: %s", parsedPayloadFile)

	// Correlate this job with the workflow run using the RunID (workflow_run_id)
	log.Printf("Workflow Job ID: %d, RunID: %d, Status: %s, Conclusion: %s",
		payload.WorkflowJob.ID,
		payload.WorkflowJob.RunID,
		payload.WorkflowJob.Status,
		payload.WorkflowJob.Conclusion,
	)

	// Further processing or updating data in DynamoDB can go here
	// For example, updating the job status in the database for a particular RunID
}







// Fetch the job logs and handle decompression if needed (e.g., zip archive)
func fetchJobLogs(logsURL, folder, timestamp, requestID string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", logsURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Add authorization header
	req.Header.Add("Authorization", "Bearer "+os.Getenv("GITHUB_TOKEN"))

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch logs: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch logs: received status %d", resp.StatusCode)
	}

	// Create a temporary file to save the zip archive
	zipFilePath := filepath.Join(folder, fmt.Sprintf("logs_%s_%s.zip", timestamp, requestID))
	out, err := os.Create(zipFilePath)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save zip archive: %v", err)
	}

	// Extract the contents of the zip archive
	err = extractZip(zipFilePath, folder)
	if err != nil {
		return fmt.Errorf("failed to extract zip file: %v", err)
	}

	fmt.Printf("Logs extracted to folder: %s\n", folder)
	return nil
}

// Extracts a zip archive to the specified folder
// Extracts a zip archive to the specified folder, and handles nested zip files if present.
func extractZip(zipFilePath, folder string) error {
	r, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %v", err)
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(folder, f.Name)

		// Create directory if it doesn't exist
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Check if the file is a zip file (based on its extension)
		if filepath.Ext(fpath) == ".zip" {
			// Save the nested zip file
			outFile, err := os.Create(fpath)
			if err != nil {
				return fmt.Errorf("failed to create nested zip file: %v", err)
			}

			rc, err := f.Open()
			if err != nil {
				return err
			}

			_, err = io.Copy(outFile, rc)
			outFile.Close()
			rc.Close()

			if err != nil {
				return fmt.Errorf("failed to save nested zip file: %v", err)
			}

			// Recursively extract the nested zip file
			err = extractZip(fpath, folder)
			if err != nil {
				return fmt.Errorf("failed to extract nested zip file: %v", err)
			}
		} else {
			// Regular file extraction
			if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}

			rc, err := f.Open()
			if err != nil {
				return err
			}

			_, err = io.Copy(outFile, rc)
			outFile.Close()
			rc.Close()

			if err != nil {
				return err
			}
		}
	}

	fmt.Printf("Extracted log files from %s\n", zipFilePath)
	return nil
}

