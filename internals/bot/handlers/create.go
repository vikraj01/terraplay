package handlers

import (
	"fmt"
	"strings"
	"github.com/bwmarrin/discordgo"
	"github.com/vikraj01/terraplay/internals/game"
)

func HandleCreateCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Fields(m.Content)
	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Usage: !create <game>")
		return
	}
	gameName := args[2]

	if !game.IsValidGame(gameName) {
		s.ChannelMessageSend(m.ChannelID, "Invalid game! Please choose a valid game: minetest, minecraft, fortnite, apex, csgo")
		return
	}

	userID := m.Author.ID
	globalName := m.Author.GlobalName

	sessionModel, workspace, err := game.CreateGameSession(userID, globalName, gameName)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	message := fmt.Sprintf(
		"```\n"+
			"Game Session Details:\n"+
			"----------------------------\n"+
			"UserID      : %s\n"+
			"GlobalName  : %s\n"+
			"Game        : %s\n"+
			"Session ID  : %s\n"+
			"Run ID      : %s\n"+
			"Workspace   : %s\n"+
			"Created At  : %s\n"+
			"Status      : %s\n"+
			"```",
		userID, globalName, gameName, sessionModel.SessionId, sessionModel.SessionId, workspace, sessionModel.CreatedAt.Format("2006-01-02 15:04:05"), sessionModel.Status)

	s.ChannelMessageSend(m.ChannelID, message)
	s.ChannelMessageSend(m.ChannelID, "Game session created! GitHub Action triggered for game: "+gameName)
}
