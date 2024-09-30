package webhook

import (
    "fmt"
    "log"
    "os"

    "github.com/bwmarrin/discordgo"
)

func sendToDiscord(userID, game, status, serverIP, runID string) {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		log.Fatalf("DISCORD_BOT_TOKEN is not set")
	}

	channelID := os.Getenv("DISCORD_CHANNEL_ID")
	if channelID == "" {
		log.Fatalf("DISCORD_CHANNEL_ID is not set")
	}

	if userID == "" {
		userID = "unknown"
	}

	var message string
	if status != "success" {
		message = fmt.Sprintf("Error encountered during deployment: %s", runID)
	} else {
		message = fmt.Sprintf("<@%s>, your game deployment for '%s' (Run ID: %s) completed with status: %s", userID, game, runID, status)
		if serverIP != "" {
			message += fmt.Sprintf("\nServer IP: %s", serverIP)
		}
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Printf("Failed to create Discord session: %v", err)
		return
	}
	defer dg.Close()

	_, err = dg.ChannelMessageSend(channelID, message)
	if err != nil {
		log.Printf("Failed to send message to Discord: %v", err)
	}
}
