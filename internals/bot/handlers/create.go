package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/vikraj01/terraplay/config"
	"github.com/vikraj01/terraplay/internals/game"
)

func HandleCreateCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Fields(m.Content)
	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "⚠️ **Usage:** `!create <game>`")
		return
	}
	gameName := args[2]

	if !config.IsValidGame(gameName) {
		s.ChannelMessageSend(m.ChannelID, "❌ **Invalid game!** Please choose a valid game: `minetest`, `minecraft`, `fortnite`, `apex`, `csgo`")
		return
	}

	userID := m.Author.ID
	globalName := m.Author.GlobalName

	sessionModel, workspace, err := game.CreateGameSession(userID, globalName, gameName)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "🚫 **Error:** "+err.Error())
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "🎮 Game Session Created! 🎮",
		Description: fmt.Sprintf("🚀 **GitHub Action triggered for game:** `%s`", gameName),
		Color:       0x00FF00,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "👤 UserID",
				Value:  fmt.Sprintf("`%s`", userID),
				Inline: true,
			},
			{
				Name:   "🌐 GlobalName",
				Value:  fmt.Sprintf("`%s`", globalName),
				Inline: true,
			},
			{
				Name:   "🕹️ Game",
				Value:  fmt.Sprintf("`%s`", gameName),
				Inline: true,
			},
			{
				Name:   "💠 Session ID",
				Value:  fmt.Sprintf("`%s`", sessionModel.SessionId),
				Inline: true,
			},
			{
				Name:   "💎 Run ID",
				Value:  fmt.Sprintf("`%s`", sessionModel.SessionId),
				Inline: true,
			},
			{
				Name:   "📁 Workspace",
				Value:  fmt.Sprintf("`%s`", workspace),
				Inline: true,
			},
			{
				Name:   "📅 Created At",
				Value:  fmt.Sprintf("`%s`", sessionModel.CreatedAt.Format("2006-01-02 15:04:05")),
				Inline: true,
			},
			{
				Name:   "🔄 Status",
				Value:  fmt.Sprintf("`%s`", sessionModel.Status),
				Inline: true,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Powered by TerraPlay",
			IconURL: "https://example.com/your-logo.png",
		},
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}
