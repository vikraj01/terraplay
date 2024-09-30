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

var commandMap = map[string]func(*discordgo.Session, *discordgo.MessageCreate){
	"!ping":         handlePingCommand,
	"!create":       handleCreateCommand,
	"!destroy":      handleDestroyCommand,
	"!list-session": handleListSessionCommand,
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	content := m.Content
	botMention := "<@" + s.State.User.ID + ">"
	content = strings.TrimPrefix(content, botMention)
	content = strings.TrimSpace(content)

	parts := strings.Fields(content)
	if len(parts) == 0 {
		return
	}

	if handler, exists := commandMap[parts[0]]; exists {
		handler(s, m)
	}
}

func handlePingCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Pong!")
}

func handleCreateCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Fields(m.Content)
	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Usage: !create <game>")
		return
	}
	gameName := args[2]

	userID := m.Author.ID
	sessions, err := dynamoService.GetActiveSessionsForUser(userID)
	if err != nil {
		log.Printf("Error fetching active sessions for user %s: %v", userID, err)
		s.ChannelMessageSend(m.ChannelID, "Failed to retrieve active sessions.")
		return
	}
	if len(sessions) > 5 {
		s.ChannelMessageSend(m.ChannelID, "One User Can Only Create 5 Game Session At Most")
		return
	}

	runId := utils.GenerateUUID()
	// uniqueId := fmt.Sprintf("%s_%s_%d", userID, runId, len(sessions)+1)
	inputs := map[string]string{
		"game":    gameName,
		"user_id": userID,
		"run_id":  runId,
	}
	fmt.Print(inputs)

	err = github.TriggerGithubAction("vikraj01", "terraplay", "start.game.yml", "main", inputs)
	fmt.Print(err)
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
		Status:    "waiting",
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

func handleDestroyCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Game session destroyed!")
}

func handleListSessionCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Listing all game sessions!")
}
