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
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

type WorkflowRunPayload struct {
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
		LogsURL    string `json:"logs_url"`
	} `json:"workflow_run"`
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

	if githubEvent == "workflow_run" {
		handleWorkflowRun(body, folder, timestamp, requestID)
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

func fetchJobLogs(logsURL, folder, timestamp, requestID string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", logsURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("Authorization", "Bearer "+os.Getenv("GITHUB_TOKEN"))

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch logs: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch logs: received status %d", resp.StatusCode)
	}

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

	err = extractZip(zipFilePath, folder)
	if err != nil {
		return fmt.Errorf("failed to extract zip file: %v", err)
	}

	fmt.Printf("Logs extracted to folder: %s\n", folder)

	// Parse the extracted files for outputs
	values, err := parseExtractedFiles(folder)
	if err != nil {
		return fmt.Errorf("failed to parse extracted files: %v", err)
	}

	log.Printf("Extracted values: %v", values)

	if len(values) > 0 {
		for _, value := range values {
			sendToDiscord(value["user_id"], value["game"], "success", value["server_ip"], value["run_id"])
		}
	}

	// cleanupExtractedFiles(folder)

	return nil
}

func extractZip(zipFilePath, folder string) error {
	r, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %v", err)
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(folder, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if filepath.Ext(fpath) == ".zip" {
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

			err = extractZip(fpath, folder)
			if err != nil {
				return fmt.Errorf("failed to extract nested zip file: %v", err)
			}
		} else {
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



func parseExtractedFiles(folder string) ([]map[string]string, error) {
	var extractedValues []map[string]string

	// List all files in the folder
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %v", folder, err)
	}

	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(folder, file.Name())
			
			// Read the file content
			data, err := ioutil.ReadFile(filePath)
			if err != nil {
				return nil, fmt.Errorf("failed to read file %s: %v", filePath, err)
			}

			content := string(data)
			
			// Check for relevant patterns in the file content
			if strings.Contains(content, "game") || strings.Contains(content, "run_id") || strings.Contains(content, "user_id") {
				extracted := extractValues(content)

				// Add to results only if there are relevant extracted values
				if extracted["game"] != "" || extracted["run_id"] != "" || extracted["user_id"] != "" || extracted["server_ip"] != "" {
					extractedValues = append(extractedValues, extracted)
				}
			}
		}
	}

	return extractedValues, nil
}

func extractValues(content string) map[string]string {
	values := map[string]string{
		"game":      "",
		"run_id":    "",
		"user_id":   "",
		"server_ip": "",
	}

	gamePattern := regexp.MustCompile(`"game":\s*"(.+?)"`)
	runIDPattern := regexp.MustCompile(`"run_id":\s*"(.+?)"`)
	userIDPattern := regexp.MustCompile(`"user_id":\s*"(.+?)"`)
	serverIPPattern := regexp.MustCompile(`(?i)(server_ip)\s*[=:]\s*"([^"]+)"`)

	if match := gamePattern.FindStringSubmatch(content); len(match) > 2 {
		values["game"] = match[2]
	}
	if match := runIDPattern.FindStringSubmatch(content); len(match) > 2 {
		values["run_id"] = match[2]
	}
	if match := userIDPattern.FindStringSubmatch(content); len(match) > 2 {
		values["user_id"] = match[2]
	}
	if match := serverIPPattern.FindStringSubmatch(content); len(match) > 2 {
		values["server_ip"] = match[2]
	}

	return values
}


func cleanupExtractedFiles(folder string) error {
	err := os.RemoveAll(folder)
	if err != nil {
		return fmt.Errorf("failed to cleanup folder: %v", err)
	}
	fmt.Printf("Folder %s cleaned up successfully\n", folder)
	return nil
}

func handleWorkflowRun(body []byte, folder, timestamp, requestID string) {
	rawBodyFile := filepath.Join(folder, fmt.Sprintf("workflow_run_raw_%s_%s.json", timestamp, requestID))
	err := ioutil.WriteFile(rawBodyFile, body, 0644)
	if err != nil {
		log.Printf("Failed to write workflow_run raw body to file: %v", err)
	} else {
		log.Printf("workflow_run raw body written to: %s", rawBodyFile)
	}

	var payload WorkflowRunPayload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		log.Printf("Failed to parse workflow_run payload: %v", err)
		return
	}

	if payload.WorkflowRun.Status == "completed" {
		log.Printf("Fetching logs for workflow run: %d", payload.WorkflowRun.ID)
		err := fetchJobLogs(payload.WorkflowRun.LogsURL, folder, timestamp, requestID)
		if err != nil {
			log.Printf("Failed to fetch logs for workflow run: %v", err)
		}
	}
}
