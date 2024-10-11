package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/vikraj01/terraplay/internals/dynamodb"
	"github.com/vikraj01/terraplay/internals/game"
	"github.com/vikraj01/terraplay/pkg/models"
)

func HandleListSessionCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	dynamoService, err := dynamodb.InitializeDynamoDB()
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "⚠️ Error: Could not connect to the database.")
		return
	}

	args := strings.Fields(m.Content)
	if len(args) <= 2 {
		s.ChannelMessageSend(m.ChannelID, "Usage: `!list-session <status>` (options: running, terminated, pending, all)")
		return
	}

	statusInput := strings.ToLower(args[2])
	userID := m.Author.ID

	sessions, err := game.ListSessions(userID, statusInput, dynamoService)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("⚠️ %s", err.Error()))
		return
	}

	message := formatSessionDetails(sessions)
	s.ChannelMessageSend(m.ChannelID, message)
}

func HandleListGamesCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	games := game.GetGames()
	message := formatGameList(games)
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

func formatGameList(games []string) string {
	var gameListBuilder strings.Builder
	gameListBuilder.WriteString("**🎮 Available Games 🎮**\n\n")
	for _, game := range games {
		gameListBuilder.WriteString(fmt.Sprintf("• %s\n", game))
	}
	gameListBuilder.WriteString("\nType `!create <game_name>` to create a game server!")
	return gameListBuilder.String()
}
