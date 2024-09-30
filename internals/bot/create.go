package bot

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/vikraj01/terraplay/internals/dynamodb"
	"github.com/vikraj01/terraplay/internals/github"
	"github.com/vikraj01/terraplay/internals/utils"
	"github.com/vikraj01/terraplay/pkg/models"
)
var dynamoService *dynamodb.DynamoDBService

func handleCreateCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Fields(m.Content)
	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Usage: !create <game>")
		return
	}
	gameName := args[1] // Fixed index for game name

	userID := m.Author.ID
	sessions, err := dynamoService.GetActiveSessionsForUser(userID)
	if err != nil {
		log.Printf("Error fetching active sessions for user %s: %v", userID, err)
		s.ChannelMessageSend(m.ChannelID, "Failed to retrieve active sessions.")
		return
	}
	if len(sessions) > 5 {
		s.ChannelMessageSend(m.ChannelID, "One User Can Only Create 5 Game Sessions At Most")
		return
	}

	runId := utils.GenerateUUID()
	uniqueId := fmt.Sprintf("%s_%s_%d", userID, runId, len(sessions)+1)
	inputs := map[string]string{
		"game":    gameName,
		"user_id": uniqueId,
		"run_id":  runId,
	}

	err = github.TriggerGithubAction("vikraj01", "terraplay", "start.game.yml", "main", inputs)
	if err != nil {
		fmt.Println(err)
		s.ChannelMessageSend(m.ChannelID, "Failed to trigger GitHub Action to create game session!")
		return
	}

	workspace := fmt.Sprintf("%s_%d@%s", userID, len(sessions)+1, gameName)
	sessionModel := models.Session{
		SessionId: runId,
		UserId:    userID,
		GameName:  gameName,
		Status:    "pending",
		StartTime: time.Now(),
		ServerIP:  "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		WorkSpace: workspace,
	}

	err = dynamoService.SaveSession(sessionModel)
	if err != nil {
		log.Printf("Failed to save session for user %s: %v", userID, err)
		s.ChannelMessageSend(m.ChannelID, "Failed to save game session data.")
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Game session created! GitHub Action triggered for game: "+gameName)
}
