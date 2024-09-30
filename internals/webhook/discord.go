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
		message = fmt.Sprintf(
			"🚨 <@%s>, there was an issue with the deployment of your game '%s'.\n"+
				"Run ID: `%s`\n"+
				"Status: `%s`\n\n"+
				"Please review the logs or contact support.",
			userID, game, runID, status)
	} else {
		message = fmt.Sprintf(
			"🎉 <@%s>, your game '%s' has been successfully deployed!\n"+
				"Run ID: `%s`\n"+
				"Status: `%s`",
			userID, game, runID, status)
		if serverIP != "" {
			message += fmt.Sprintf("\n🖥️ Server IP: `%s`\n", serverIP)
		}
		message += "\nEnjoy your game!"
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

func sendToDiscordForStopAction(runID, status string) error {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		return fmt.Errorf("DISCORD_BOT_TOKEN is not set")
	}

	channelID := os.Getenv("DISCORD_CHANNEL_ID")
	if channelID == "" {
		return fmt.Errorf("DISCORD_CHANNEL_ID is not set")
	}

	var message string
	if status == "terminated" {
		message = fmt.Sprintf(
			"🛑 **Session Terminated**\n"+
				"Run ID: `%s`\n"+
				"Status: `%s`\n\n"+
				"The session has been successfully destroyed.",
			runID, status)
	} else {
		message = fmt.Sprintf(
			"⚠️ **Session Termination Error**\n"+
				"Run ID: `%s`\n"+
				"Status: `%s`\n\n"+
				"There was an issue terminating the session.",
			runID, status)
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Printf("Failed to create Discord session: %v", err)
		return err
	}
	defer dg.Close()

	_, err = dg.ChannelMessageSend(channelID, message)
	if err != nil {
		log.Printf("Failed to send message to Discord: %v", err)
		return err
	}

	return nil
}
