package game

import (
	"fmt"
	"log"

	"github.com/vikraj01/terraplay/internals/dynamodb"
	"github.com/vikraj01/terraplay/internals/github"
)

func DestroyGameSession(sessionId string, dynamodbService *dynamodb.DynamoDBService) (string, error) {
	details, err := dynamodbService.GetDetailsBySessionID(sessionId)
	if err != nil {
		log.Printf("Error fetching sessionId details: %v", err)
		return "", fmt.Errorf("could not find workspace for the given session ID")
	}

	workspaceID := details.Workspace
	log.Printf("Workspace ID: %s", workspaceID)

	inputs := map[string]string{
		"run_id":    sessionId,
		"workspace": workspaceID,
	}

	err = github.TriggerGithubAction("vikraj01", "terraplay", "stop.game.yml", "main", inputs)
	if err != nil {
		log.Printf("Failed to trigger GitHub Action: %v", err)
		return "", fmt.Errorf("failed to trigger the destruction of the game session")
	}

	return workspaceID, nil
}
