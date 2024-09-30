package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)



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





