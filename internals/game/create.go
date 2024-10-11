package game

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/vikraj01/terraplay/config"
	"github.com/vikraj01/terraplay/internals/dynamodb"
	"github.com/vikraj01/terraplay/internals/github"
	"github.com/vikraj01/terraplay/internals/utils"
	"github.com/vikraj01/terraplay/pkg/models"
)

var validGames = []string{"minetest", "minecraft", "fortnite", "apex", "csgo"}

func IsValidGame(gameName string) bool {
	for _, game := range validGames {
		if strings.EqualFold(game, gameName) {
			return true
		}
	}
	return false
}

func CreateGameSession(userID, globalName, gameName string) (models.Session, string, error) {
	dynamoService, err := dynamodb.InitializeDynamoDB()
	if err != nil {
		log.Fatalf("Failed to initialize DynamoDB: %v", err)
	}
	status := "running"
	statusEnum, validStatus := config.StatusValues[status]
	if !validStatus {
		return models.Session{}, "", fmt.Errorf("invalid status")
	}
	statusString := config.StatusNames[statusEnum]
	sessions, err := dynamoService.GetActiveSessionsForUser(userID, statusString)
	if err != nil {
		log.Printf("Error fetching active sessions for user %s: %v", userID, err)
		return models.Session{}, "", fmt.Errorf("failed to retrieve active sessions")
	}

	if len(sessions) >= 5 {
		return models.Session{}, "", fmt.Errorf("one user can only create 5 game sessions at most")
	}

	runId := utils.GenerateCryptoID(12)
	uniqueId := fmt.Sprintf("%s_%s_%s_%d", globalName, userID, runId, len(sessions)+1)
	inputs := map[string]string{
		"game":    gameName,
		"user_id": uniqueId,
		"run_id":  runId,
	}

	err = github.TriggerGithubAction("vikraj01", "terraplay", "start.game.yml", "main", inputs)
	if err != nil {
		log.Printf("Failed to trigger GitHub Action: %v", err)
		return models.Session{}, "", fmt.Errorf("failed to trigger GitHub Action")
	}

	workspace := fmt.Sprintf("%s@%s", uniqueId, gameName)
	sessionModel := models.Session{
		SessionId:  runId,
		UserId:     userID,
		GameName:   gameName,
		Status:     "pending",
		StartTime:  time.Now(),
		ServerIP:   "",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		WorkSpace:  workspace,
		InstanceId: "",
	}

	err = dynamoService.SaveSession(sessionModel)
	if err != nil {
		log.Printf("Failed to save session for user %s: %v", userID, err)
		return models.Session{}, "", fmt.Errorf("failed to save game session data")
	}

	return sessionModel, workspace, nil
}
