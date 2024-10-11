package bot

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/vikraj01/terraplay/config"
	"github.com/vikraj01/terraplay/internals/dynamodb"
	"github.com/vikraj01/terraplay/pkg/models"
)
var validGames = []string{"minetest", "minecraft", "fortnite", "apex", "csgo"}


func handleListSessionCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	dynamoService, err := dynamodb.InitializeDynamoDB()
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error Occured while connecting to the dynamodb-database")
		return 
	}

	args := strings.Fields(m.Content)
	if len(args) <= 2 {
		s.ChannelMessageSend(m.ChannelID, "Usage: !list-session <status> (options: running, terminated, pending, all)")
		return
	}

	statusInput := strings.ToLower(args[2])

	var statusEnum config.Status
	var validStatus bool
	var statusString string

	if statusInput == "all" {
		statusString = "all"
	} else {
		statusEnum, validStatus = config.StatusValues[statusInput]
		if !validStatus {
			s.ChannelMessageSend(m.ChannelID, "Invalid status! Use one of the following: running, terminated, pending, all.")
			return
		}
		statusString = config.StatusNames[statusEnum]
	}

	userID := m.Author.ID
	var sessions []models.Session

	if statusString == "all" {
		sessions, err = dynamoService.GetActiveSessionsForUser(userID, "all")
	} else {
		sessions, err = dynamoService.GetActiveSessionsForUser(userID, statusString)
	}

	if err != nil {
		log.Printf("Error fetching sessions for user %s: %v", userID, err)
		s.ChannelMessageSend(m.ChannelID, "Failed to retrieve sessions.")
		return
	}

	if len(sessions) == 0 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No sessions found with status: %s", statusString))
		return
	}

	message := formatSessionDetails(sessions)
	s.ChannelMessageSend(m.ChannelID, message)
}

func formatSessionDetails(sessions []models.Session) string {
	var messageBuilder strings.Builder
	messageBuilder.WriteString("**Listing Sessions:**\n\n")
	for _, session := range sessions {
		messageBuilder.WriteString(fmt.Sprintf(
			"**Game:** %s\n"+
				"**Status:** %s\n"+
				"**Run ID:** `%s`\n"+
				"**Start Time:** %s\n\n",
			session.GameName,
			session.Status,
			session.SessionId,
			session.StartTime.Format(time.RFC1123),
		))
	}
	return messageBuilder.String()
}

func handleListGamesCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	gameList := "**🎮 Available Games 🎮**\n\n"
	for _, game := range validGames {
		gameList += game + "\n"
	}
	gameList += "\nType `!create <game_name>` to create a game server!"

	s.ChannelMessageSend(m.ChannelID, gameList)
}
