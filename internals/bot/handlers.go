package bot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/vikraj01/terraplay/internals/bot/handlers"
)

var commandMap = map[string]func(*discordgo.Session, *discordgo.MessageCreate){
	"!ping":          handlers.HandlePingCommand,
	"!create":        handlers.HandleCreateCommand,
	"!destroy":       handlers.HandleDestroyCommand,
	"!list-sessions": handlers.HandleListSessionCommand,
	"!list-games":    handlers.HandleListGamesCommand,
	"!stop":          handlers.HandleStopCommand,
	"!restart":       handlers.HandleRestartCommand,
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	content := m.Content
	botMention := "<@" + s.State.User.ID + ">"
	content = strings.TrimPrefix(content, botMention)
	content = strings.TrimSpace(content)

	parts := strings.Fields(content)
	if len(parts) == 0 {
		return
	}

	if handler, exists := commandMap[parts[0]]; exists {
		handler(s, m)
	}
}
