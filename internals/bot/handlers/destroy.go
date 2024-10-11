package handlers

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/vikraj01/terraplay/internals/dynamodb"
	"github.com/vikraj01/terraplay/internals/github"
)

func HandleDestroyCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	dynamodbService, err := dynamodb.InitializeDynamoDB()
	if err != nil {
		log.Printf("Error initializing DynamoDB: %v", err)
		s.ChannelMessageSend(m.ChannelID, "⚠️ Error: Could not initialize database. Please try again later.")
		return
	}

	args := strings.Fields(m.Content)
	if len(args) < 3 {
		s.ChannelMessageSend(m.ChannelID, "⚠️ Usage: `!destroy <session_id>`")
		return
	}
	sessionId := args[2]

	details, err := dynamodbService.GetDetailsBySessionID(sessionId)
	workspaceID := details.Workspace
	if err != nil {
		log.Printf("Error fetching sessionId details: %v", err)
		s.ChannelMessageSend(m.ChannelID, "⚠️ Error: Could not find workspace for the given session ID.")
		return
	}
	log.Printf("Workspace ID: %s", workspaceID)

	inputs := map[string]string{
		"run_id":    sessionId,
		"workspace": workspaceID,
	}

	err = github.TriggerGithubAction("vikraj01", "terraplay", "stop.game.yml", "main", inputs)
	if err != nil {
		log.Printf("Failed to trigger GitHub Action: %v", err)
		s.ChannelMessageSend(m.ChannelID, "⚠️ Error: Failed to trigger the destruction of the game session.")
		return
	}

	message := fmt.Sprintf(
		"🛠️ **Destruction of the game session has been initiated!**\n\n"+
			"**Session ID:** `%s`\n"+
			"**Workspace ID:** `%s`\n\n"+
			"Your game session is now in the process of being stopped. "+
			"Thank you for using the Terraplay service!",
		sessionId, workspaceID,
	)
	s.ChannelMessageSend(m.ChannelID, message)
}
