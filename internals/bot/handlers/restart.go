package handlers

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/vikraj01/terraplay/internals/dynamodb"
	"github.com/vikraj01/terraplay/internals/utils"
)

func HandleRestartCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	dynamodbService, err := dynamodb.InitializeDynamoDB()
	if err != nil {
		log.Printf("Error initializing DynamoDB: %v", err)
		s.ChannelMessageSend(m.ChannelID, "⚠️ Error: Could not initialize database. Please try again later.")
		return
	}

	args := strings.Fields(m.Content)
	if len(args) < 3 {
		s.ChannelMessageSend(m.ChannelID, "⚠️ Usage: `!restart <session_id>`")
		return
	}

	sessionId := args[2]
	details, err := dynamodbService.GetDetailsBySessionID(sessionId)
	if err != nil {
		log.Printf("Error fetching sessionId details: %v", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("⚠️ Error: Could not find workspace for the given session ID: %v", err))
		return
	}

	if details.InstanceId == "" {
		log.Println("⚠️ Error: Instance ID is missing in session details.")
		s.ChannelMessageSend(m.ChannelID, "⚠️ Error: Instance ID is missing in the session data.")
		return
	}

	sshKeyBase64 := os.Getenv("EC2_SSH_KEY_BASE64")
	if sshKeyBase64 == "" {
		log.Println("⚠️ Error: EC2_SSH_KEY_BASE64 is not set")
		s.ChannelMessageSend(m.ChannelID, "⚠️ Error: SSH private key is missing.")
		return
	}

	privateKey, err := base64.StdEncoding.DecodeString(sshKeyBase64)
	if err != nil {
		log.Printf("Error decoding base64 private key: %v", err)
		s.ChannelMessageSend(m.ChannelID, "⚠️ Error decoding SSH private key.")
		return
	}

	err = utils.StartEC2Instance(details.InstanceId, os.Getenv("AWS_REGION"))
	if err != nil {
		log.Printf("Error starting EC2 instance: %v", err)
		s.ChannelMessageSend(m.ChannelID, "⚠️ Error starting EC2 instance.")
		return
	}

	awsRegion := os.Getenv("AWS_REGION")
	newServerIP, err := utils.GetPublicIPByInstanceID(details.InstanceId, awsRegion)
	if err != nil {
		log.Printf("Error retrieving new server IP: %v", err)
		s.ChannelMessageSend(m.ChannelID, "⚠️ Error retrieving new server IP.")
		return
	}

	sshConfig := utils.SSHConfig{
		Host:       newServerIP,
		Port:       "22",
		User:       "ec2-user",
		PrivateKey: privateKey,
	}

	backupFile := "/tmp/backup.tar.gz"
	backupPath := "/opt/minetest/data"
	s3Bucket := "global-bucket-893606"

	err = RestoreAndRestartEC2(sshConfig, backupPath, s3Bucket, backupFile, newServerIP, details.InstanceId)
	if err != nil {
		log.Printf("Error restarting and restoring EC2: %v", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("⚠️ Error: %v", err))
		return
	}

	err = dynamodbService.UpdateSessionStatusAndIP(sessionId, "running", newServerIP)
	if err != nil {
		log.Printf("Error updating session with new IP: %v", err)
		s.ChannelMessageSend(m.ChannelID, "⚠️ Error updating session with new IP.")
		return
	}

	message := fmt.Sprintf(
		"EC2 instance with IP `%s` has been restarted and data restored. Workspace: `%s`", newServerIP, details.Workspace)
	s.ChannelMessageSend(m.ChannelID, message)
}

func RestoreAndRestartEC2(sshConfig utils.SSHConfig, backupPath, s3Bucket, backupFile, publicIP string, instanceID string) error {
	awsRegion := os.Getenv("AWS_REGION")

	err := utils.WaitForInstanceRunning(instanceID, awsRegion)
	if err != nil {
		return fmt.Errorf("error waiting for instance to reach running state: %v", err)
	}

	client, err := utils.ConnectToEC2ViaSSHWithRetry(sshConfig)
	if err != nil {
		return fmt.Errorf("error connecting to EC2 via SSH: %v", err)
	}
	defer client.Close()

	scriptContent := `
		#!/bin/bash
		BACKUP_FILE=$1
		BACKUP_PATH=$2
		S3_BUCKET=$3
		AWS_SECRET_ACCESS_KEY=$4
		AWS_ACCESS_KEY_ID=$5
		AWS_REGION=$6

		if ! command -v aws &> /dev/null; then
			echo "AWS CLI not found, installing..."
			sudo yum install -y awscli
		fi

		export AWS_ACCESS_KEY_ID=$5
		export AWS_SECRET_ACCESS_KEY=$4
		export AWS_DEFAULT_REGION=$6

		echo "Downloading backup from S3 bucket $S3_BUCKET..."
		aws s3 cp s3://$S3_BUCKET/backup.tar.gz $BACKUP_FILE --region $AWS_REGION

		echo "Extracting backup..."
		tar -xzf $BACKUP_FILE -C $BACKUP_PATH

		echo "Cleaning up the backup file..."
		rm -f $BACKUP_FILE
	`

	createScriptCommand := fmt.Sprintf("echo '%s' > /home/ec2-user/restore.sh && chmod +x /home/ec2-user/restore.sh", scriptContent)
	err = utils.RunCommandOnEC2(client, createScriptCommand)
	if err != nil {
		return fmt.Errorf("error creating restore.sh script on EC2: %v", err)
	}
	log.Println("Restore script created successfully on EC2 instance.")

	restoreCommand := fmt.Sprintf("/bin/bash /home/ec2-user/restore.sh %s %s %s %s %s %s", backupFile, backupPath, s3Bucket, os.Getenv("AWS_SECRET_ACCESS_KEY"), os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_REGION"))
	err = utils.RunCommandOnEC2(client, restoreCommand)
	if err != nil {
		return fmt.Errorf("error running restore script on EC2: %v", err)
	}

	log.Println("Data restore completed successfully on EC2 instance.")
	return nil
}
