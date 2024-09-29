package webhook

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var payload WorkflowRunPayload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		log.Printf("Failed to parse webhook payload: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if payload.WorkflowRun.Status == "completed" {
		sendToDiscord(payload.WorkflowRun.Inputs.UserID, payload.WorkflowRun.Inputs.Game, payload.WorkflowRun.Conclusion, payload.WorkflowRun.Outputs.ServerIP, payload.WorkflowRun.Inputs.RunID)
	}
}

func sendToDiscord(userID, game, status, serverIP, runID string) {
	token := "YOUR_DISCORD_BOT_TOKEN"
	channelID := "YOUR_DISCORD_CHANNEL_ID"

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Failed to create Discord session: %v", err)
	}

	message := fmt.Sprintf("<@%s>, your game deployment for '%s' (Run ID: %s) completed with status: %s", userID, game, runID, status)

	if serverIP != "" {
		message += fmt.Sprintf("\nServer IP: %s", serverIP)
	}

	_, err = dg.ChannelMessageSend(channelID, message)
	if err != nil {
		log.Printf("Failed to send message to Discord: %v", err)
	}

	dg.Close()
}
