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

func HandleStopCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	dynamodbService, err := dynamodb.InitializeDynamoDB()
	if err != nil {
		log.Printf("Error initializing DynamoDB: %v", err)
		s.ChannelMessageSend(m.ChannelID, "⚠️ Error: Could not initialize database. Please try again later.")
		return
	}

	args := strings.Fields(m.Content)
	if len(args) < 3 {
		s.ChannelMessageSend(m.ChannelID, "⚠️ Usage: `!stop <session_id>`")
		return
	}

	sessionId := args[2]
	details, err := dynamodbService.GetDetailsBySessionID(sessionId)
	if err != nil {
		log.Printf("Error fetching sessionId details: %v", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("⚠️ Error: Could not find workspace for the given session ID: %v", err))
		return
	}

	user := "ec2-user"
	sshConfig, err := utils.GetSSHConfig(details.ServerIP, user)
	if err != nil {
		log.Printf("Error establishing ssh connection: %v", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("⚠️ Error: %v", err))
	}

	gameConfig := config.FindGameConfig(details.GameName)
	backupFile := "/tmp/backup.tar.gz"
	backupPath := gameConfig.VolumePath
	s3Bucket := os.Getenv("AWS_GAME_BACKUP_BUCKET")

	err = game.BackupAndStopEC2(sshConfig, backupPath, s3Bucket, backupFile, details.ServerIP)
	if err != nil {
		log.Printf("Error executing backup and stop: %v", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("⚠️ Error: %v", err))
		return
	}
	dynamodbService.UpdateSessionStatusAndIP(sessionId, "halted", details.ServerIP)

	message := fmt.Sprintf(
		"EC2 instance with IP `%s` has been backed up and stopped. Workspace: `%s`", details.ServerIP, details.Workspace)
	s.ChannelMessageSend(m.ChannelID, message)
}
