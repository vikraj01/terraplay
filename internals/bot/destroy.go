package bot

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/vikraj01/terraplay/internals/dynamodb"
	"github.com/vikraj01/terraplay/internals/github"
)

func handleDestroyCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	dynamodbService, err := dynamodb.InitializeDynamoDB()
	if err != nil {
		return
	}
	args := strings.Fields(m.Content)
	sessionId := args[2]

	workspaceID, err := dynamodbService.GetWorkspaceBySessionID(sessionId)
	if err != nil {
		log.Printf("Error fetching workspace: %v", err)
	} else {
		log.Printf("Workspace ID: %s", workspaceID)
	}

	inputs := map[string]string{
		"run_id":    sessionId,
		"workspace": workspaceID,
	}

	err = github.TriggerGithubAction("vikraj01", "terraplay", "stop.game.yml", "main", inputs)
	if err != nil {
		fmt.Println(err)
		s.ChannelMessageSend(m.ChannelID, "Failed to trigger GitHub Action to create game session!")
		return
	}
	s.ChannelMessageSend(m.ChannelID, "Game session destroyed!")
}
