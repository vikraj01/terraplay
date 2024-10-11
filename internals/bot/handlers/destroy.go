package handlers

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/vikraj01/terraplay/internals/dynamodb"
	"github.com/vikraj01/terraplay/internals/game"
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

	workspaceID, err := game.DestroyGameSession(sessionId, dynamodbService)
	if err != nil {
		log.Printf("Error during game session destruction: %v", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("⚠️ Error: %s", err.Error()))
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
