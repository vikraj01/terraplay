package handlers

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/vikraj01/terraplay/config"
	"github.com/vikraj01/terraplay/internals/dynamodb"
	"github.com/vikraj01/terraplay/internals/game"
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

	newServerIP, err := game.RestartEC2(details.InstanceId)
	if err != nil {
		log.Printf("Error retrieving new server IP: %v", err)
		s.ChannelMessageSend(m.ChannelID, "⚠️ Error retrieving new server IP.")
		return
	}

	user := "ec2-user"
	sshConfig, err := utils.GetSSHConfig(newServerIP, user)
	if err != nil {
		log.Printf("Error establishing ssh connection: %v", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("⚠️ Error: %v", err))
	}

	gameConfig := config.FindGameConfig(details.GameName)
	backupFile := "/tmp/backup.tar.gz"
	backupPath := gameConfig.VolumePath
	s3Bucket := os.Getenv("AWS_GAME_BACKUP_BUCKET")

	err = game.RestoreEC2(sshConfig, backupPath, s3Bucket, backupFile, newServerIP, details.InstanceId)
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
