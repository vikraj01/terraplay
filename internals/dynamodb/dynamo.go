package dynamodb

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/vikraj01/terraplay/pkg/models"
)

type DynamoDBService struct {
	Client *dynamodb.DynamoDB
}

func InitializeDynamoDB() (*DynamoDBService, error) {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		log.Println("AWS REGION is not set")
		return nil, fmt.Errorf("AWS REGION is not set")
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		log.Printf("Failed to create AWS session: %v", err)
		return nil, err
	}

	dynamoClient := dynamodb.New(sess)

	return &DynamoDBService{Client: dynamoClient}, nil
}

func (svc *DynamoDBService) SaveSession(sessModel models.Session) error {
	table := os.Getenv("DYNAMO_TABLE")
	input := &dynamodb.PutItemInput{
		TableName: aws.String(table),
		Item: map[string]*dynamodb.AttributeValue{
			"session_id": {
				S: aws.String(sessModel.SessionId),
			},
			"user_id": {
				S: aws.String(sessModel.UserId),
			},
			"game_name": {
				S: aws.String(sessModel.GameName),
			},
			"status": {
				S: aws.String(sessModel.Status),
			},
			"start_time": {
				S: aws.String(sessModel.StartTime.Format(time.RFC3339)),
			},
			"created_at": {
				S: aws.String(sessModel.CreatedAt.Format(time.RFC3339)),
			},
			"updated_at": {
				S: aws.String(sessModel.UpdatedAt.Format(time.RFC3339)),
			},
			"workspace": {
				S: aws.String(sessModel.WorkSpace),
			},
			"server_ip": {
				S: aws.String(sessModel.ServerIP),
			},
		},
	}

	_, err := svc.Client.PutItem(input)
	if err != nil {
		log.Printf("Failed to save session to DynamoDB: %v", err)
		return err
	}

	log.Printf("Session saved successfully!")
	return nil
}

func (svc *DynamoDBService) GetActiveSessionsForUser(userID string, status string) ([]models.Session, error) {
	table := os.Getenv("DYNAMO_TABLE")
	
	input := &dynamodb.QueryInput{
		TableName:              aws.String(table),
		IndexName:              aws.String("user_id-index"),
		KeyConditionExpression: aws.String("user_id = :user_id"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":user_id": {
				S: aws.String(userID),
			},
		},
	}

	if status != "all" {
		input.FilterExpression = aws.String("#status = :status")
		input.ExpressionAttributeNames = map[string]*string{
			"#status": aws.String("status"),
		}
		input.ExpressionAttributeValues[":status"] = &dynamodb.AttributeValue{
			S: aws.String(status),
		}
	}

	result, err := svc.Client.Query(input)
	if err != nil {
		log.Printf("Failed to query active sessions: %v", err)
		return nil, err
	}

	var sessions []models.Session
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &sessions)
	if err != nil {
		log.Printf("Failed to unmarshal query result: %v", err)
		return nil, err
	}

	return sessions, nil
}

func (svc *DynamoDBService) UpdateSessionStatusAndIP(sessionID, status, serverIP string) error {
	table := os.Getenv("DYNAMO_TABLE")

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(table),
		Key: map[string]*dynamodb.AttributeValue{
			"session_id": {
				S: aws.String(sessionID),
			},
		},
		UpdateExpression: aws.String("SET #status = :status, #server_ip = :server_ip, #updated_at = :updated_at"),
		ExpressionAttributeNames: map[string]*string{
			"#status":     aws.String("status"),
			"#server_ip":  aws.String("server_ip"),
			"#updated_at": aws.String("updated_at"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":status": {
				S: aws.String(status),
			},
			":server_ip": {
				S: aws.String(serverIP),
			},
			":updated_at": {
				S: aws.String(time.Now().Format(time.RFC3339)),
			},
		},
		ReturnValues: aws.String("UPDATED_NEW"),
	}

	result, err := svc.Client.UpdateItem(input)
	if err != nil {
		log.Printf("Failed to update session: %v", err)
		return err
	}

	log.Printf("Session updated successfully: %v", result.Attributes)
	return nil
}


func (svc *DynamoDBService) GetWorkspaceBySessionID(sessionID string) (string, error) {
	table := os.Getenv("DYNAMO_TABLE")

	// Define the query input to fetch the session based on the session_id
	input := &dynamodb.GetItemInput{
		TableName: aws.String(table),
		Key: map[string]*dynamodb.AttributeValue{
			"session_id": {
				S: aws.String(sessionID),
			},
		},
		ProjectionExpression: aws.String("workspace"), // Only return the workspace field
	}

	// Perform the GetItem query
	result, err := svc.Client.GetItem(input)
	if err != nil {
		log.Printf("Failed to get workspace for session ID %s: %v", sessionID, err)
		return "", err
	}

	// Check if the item was found
	if result.Item == nil {
		return "", fmt.Errorf("no session found with session_id: %s", sessionID)
	}

	// Unmarshal the result to extract the workspace value
	var workspace struct {
		Workspace string `json:"workspace"`
	}
	err = dynamodbattribute.UnmarshalMap(result.Item, &workspace)
	if err != nil {
		log.Printf("Failed to unmarshal workspace result: %v", err)
		return "", err
	}

	return workspace.Workspace, nil
}
