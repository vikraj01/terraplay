package bot

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/vikraj01/terraplay/internals/dynamodb"
)

func handleStopCommand(s *discordgo.Session, m *discordgo.Message) () {
	dynamodbService, err := dynamodb.InitializeDynamoDB()
	if err != nil {
		log.Printf("Error initializing DynamoDB: %v", err)
		s.ChannelMessageSend(m.ChannelID, "⚠️ Error: Could not initialize database. Please try again later.")
		return
	}

	args := strings.Fields(m.Content)
	if len(args) < 3 {
		s.ChannelMessageSend(m.ChannelID, "⚠️ Usage: `!stop <session_id>`")
		return
	}

	sessionId := args[2]
	details, err := dynamodbService.GetDetailsBySessionID(sessionId)
	if err != nil {
		log.Printf("Error fetching sessionId details: %v", err)
		s.ChannelMessageSend(m.ChannelID, "⚠️ Error: Could not find workspace for the given session ID.")
		return
	}
	fmt.Println(details)
}
