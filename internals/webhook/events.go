package webhook

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/vikraj01/terraplay/internals/dynamodb"
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
		err := fetchJobLogs(payload.WorkflowRun.LogsURL, folder, timestamp, requestID, payload)
		if err != nil {
			log.Printf("Failed to fetch logs for workflow run: %v", err)
		}
	}
}
func fetchJobLogs(logsURL, folder, timestamp, requestID string, payload WorkflowRunPayload) error {
	client := &http.Client{}
	dynamoService, err := dynamodb.InitializeDynamoDB()
	if err != nil {
		log.Fatalf("Failed to initialize DynamoDB: %v", err)
	}

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

	values, err := parseExtractedFiles(folder)
	if err != nil || values == nil {
		errorMessage := fmt.Sprintf("Failed to extract necessary values from logs: %v", err)
		sendToDiscord("", "", "error", "", errorMessage)
		return err
	}

	userID := values["user_id"]
	game := values["game"]
	serverIP := values["server_ip"]
	runID := values["run_id"]

	if userID == "" || game == "" || serverIP == "" || runID == "" {
		errorMessage := fmt.Sprintf("Missing values in logs: game=%s, user_id=%s, server_ip=%s, run_id=%s", game, userID, serverIP, runID)
		sendToDiscord("", "", "error", "", errorMessage)
		return fmt.Errorf(errorMessage)
	}

	status := payload.WorkflowRun.Conclusion
	sendToDiscord(userID, game, status, serverIP, runID)
	if dynamoService == nil {
		log.Println("DynamoDBService is not initialized")
		return fmt.Errorf("DynamoDBService is not initialized")
	}

	dynamoService.UpdateSessionStatusAndIP(runID, "running", serverIP)

	cleanupExtractedFiles(folder)
	return nil
}

func cleanupExtractedFiles(folder string) error {
	err := os.RemoveAll(folder)
	if err != nil {
		return fmt.Errorf("failed to cleanup folder: %v", err)
	}
	fmt.Printf("Folder %s cleaned up successfully\n", folder)
	return nil
}
