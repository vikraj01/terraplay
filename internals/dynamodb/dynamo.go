package dynamodb

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/vikraj01/terraplay/pkg/models"
)

func SaveSession(sessModel models.Session) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})
	if err != nil {
		log.Printf("Failed to create AWS session: %v", err)
		return err
	}

	svc := dynamodb.New(sess)

	input := &dynamodb.PutItemInput{
		TableName: aws.String("game_sessions"),
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
			"instance_id": {
				S: aws.String(sessModel.InstanceID),
			},
			"state_file": {
				S: aws.String(sessModel.StateFile),
			},
			"created_at": {
				S: aws.String(sessModel.CreatedAt.Format(time.RFC3339)),
			},
			"updated_at": {
				S: aws.String(sessModel.UpdatedAt.Format(time.RFC3339)),
			},
		},
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Printf("Failed to save session to DynamoDB: %v", err)
		return err
	}

	log.Printf("Session saved successfully!")
	return nil
}
