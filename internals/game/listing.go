// internals/game/list.go
package game

import (
	"fmt"

	"github.com/vikraj01/terraplay/config"
	"github.com/vikraj01/terraplay/internals/dynamodb"
	"github.com/vikraj01/terraplay/pkg/models"
)


func GetGames() []string {
	return config.GetAllGames()
}

func ListSessions(userID, statusInput string, dynamoService *dynamodb.DynamoDBService) ([]models.Session, error) {
	var statusEnum config.Status
	var validStatus bool
	var statusString string

	if statusInput == "all" {
		statusString = "all"
	} else {
		statusEnum, validStatus = config.StatusValues[statusInput]
		if !validStatus {
			return nil, fmt.Errorf("Invalid status! Use one of the following: running, terminated, pending, all.")
		}
		statusString = config.StatusNames[statusEnum]
	}

	var sessions []models.Session
	var err error

	if statusString == "all" {
		sessions, err = dynamoService.GetSessionsBasedOnStatus(userID, "all")
	} else {
		sessions, err = dynamoService.GetSessionsBasedOnStatus(userID, statusString)
	}

	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve sessions: %v", err)
	}

	if len(sessions) == 0 {
		return nil, fmt.Errorf("No sessions found with status: %s", statusString)
	}

	return sessions, nil
}
