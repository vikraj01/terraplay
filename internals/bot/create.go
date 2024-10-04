package bot

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/vikraj01/terraplay/config"
	"github.com/vikraj01/terraplay/internals/dynamodb"
	"github.com/vikraj01/terraplay/internals/github"
	"github.com/vikraj01/terraplay/internals/utils"
	"github.com/vikraj01/terraplay/pkg/models"
)

var validGames = []string{"minetest", "minecraft", "fortnite", "apex", "csgo"}

func isValidGame(gameName string) bool {
	for _, game := range validGames {
		if strings.EqualFold(game, gameName) {
			return true
		}
	}
	return false
}

func handleCreateCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	dynamoService, err := dynamodb.InitializeDynamoDB()
	if err != nil {
		log.Fatalf("Failed to initialize DynamoDB: %v", err)
	}
	args := strings.Fields(m.Content)
	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Usage: !create <game>")
		return
	}
	gameName := args[2]

	if !isValidGame(gameName) {
		s.ChannelMessageSend(m.ChannelID, "Invalid game! Please choose a valid game: minetest, minecraft, fortnite, apex, csgo")
		return
	}

	userID := m.Author.ID
	GlobalName := m.Author.GlobalName

	status := "running"
	statusEnum, validStatus := config.StatusValues[status]
	if !validStatus {
		return
	}
	statusString := config.StatusNames[statusEnum]
	sessions, err := dynamoService.GetActiveSessionsForUser(userID, statusString)
	if err != nil {
		log.Printf("Error fetching active sessions for user %s: %v", userID, err)
		s.ChannelMessageSend(m.ChannelID, "Failed to retrieve active sessions.")
		return
	}

	if len(sessions) >= 5 {
		s.ChannelMessageSend(m.ChannelID, "One User Can Only Create 5 Game Sessions At Most")
		return
	}

	runId := utils.GenerateCryptoID(12)
	uniqueId := fmt.Sprintf("%s_%s_%s_%d", GlobalName, userID, runId, len(sessions)+1)
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

	workspace := fmt.Sprintf("%s@%s", uniqueId, gameName)
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
		InstanceId: "",
	}

	err = dynamoService.SaveSession(sessionModel)
	if err != nil {
		log.Printf("Failed to save session for user %s: %v", userID, err)
		s.ChannelMessageSend(m.ChannelID, "Failed to save game session data.")
		return
	}

	message := fmt.Sprintf(
		"```\n"+
			"Game Session Details:\n"+
			"----------------------------\n"+
			"UserID      : %s\n"+
			"GlobalName  : %s\n"+
			"Game        : %s\n"+
			"Session ID  : %s\n"+
			"Run ID      : %s\n"+
			"Workspace   : %s\n"+
			"Created At  : %s\n"+
			"Status      : %s\n"+
			"```",
		userID, GlobalName, gameName, sessionModel.SessionId, runId, workspace, sessionModel.CreatedAt.Format(time.RFC822), sessionModel.Status)

	s.ChannelMessageSend(m.ChannelID, message)
	s.ChannelMessageSend(m.ChannelID, "Game session created! GitHub Action triggered for game: "+gameName)
}
