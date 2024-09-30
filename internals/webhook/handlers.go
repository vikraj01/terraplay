package webhook

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "time"
    "path/filepath"

    "github.com/google/uuid"
    "github.com/bwmarrin/discordgo"
)

type WorkflowRunPayload struct {
    Action      string `json:"action"`
    WorkflowRun struct {
        Status     string `json:"status"`
        Conclusion string `json:"conclusion"`
        Outputs    struct {
            ServerIP string `json:"server_ip"`
        } `json:"outputs"`
        Inputs struct {
            RunID  string `json:"run_id"`
            UserID string `json:"user_id"`
            Game   string `json:"game"`
        } `json:"inputs"`
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

    rawBodyFile := filepath.Join(folder, fmt.Sprintf("webhook_raw_body_%s_%s.json", timestamp, requestID))

    err = ioutil.WriteFile(rawBodyFile, body, 0644)
    if err != nil {
        log.Printf("Failed to write raw body to file: %v", err)
    } else {
        log.Printf("Raw body written to: %s", rawBodyFile)
    }

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

    var payload WorkflowRunPayload
    err = json.Unmarshal(body, &payload)
    if err != nil {
        log.Printf("Failed to parse webhook payload: %v", err)
        http.Error(w, "Invalid payload", http.StatusBadRequest)
        return
    }

    parsedPayloadFile := filepath.Join(folder, fmt.Sprintf("webhook_parsed_payload_%s_%s.json", timestamp, requestID))

    payloadJSON, err := json.MarshalIndent(payload, "", "  ")
    if err != nil {
        log.Printf("Failed to marshal JSON payload: %v", err)
    }	

    err = ioutil.WriteFile(parsedPayloadFile, payloadJSON, 0644)
    if err != nil {
        log.Printf("Failed to write JSON payload to file: %v", err)
    } else {
        log.Printf("Parsed payload written to: %s", parsedPayloadFile)
    }

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
