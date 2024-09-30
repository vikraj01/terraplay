package bot

import (
	"strings"
	"github.com/bwmarrin/discordgo"
)

var commandMap = map[string]func(*discordgo.Session, *discordgo.MessageCreate){
	"!ping":         handlePingCommand,
	"!create":       handleCreateCommand,
	"!destroy":      handleDestroyCommand,
	"!list-sessions": handleListSessionCommand,
	"!list-games"  : handleListGamesCommand,
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

