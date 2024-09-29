package bot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

var commandMap = map[string]func(*discordgo.Session, *discordgo.MessageCreate){
	"!ping":        handlePingCommand,
	"!create":      handleCreateCommand,
	"!destroy":     handleDestroyCommand,
	"!list-session": handleListSessionCommand,
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

func handlePingCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Pong!")
}

func handleCreateCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Game session created!")
}

func handleDestroyCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Game session destroyed!")
}

func handleListSessionCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Listing all game sessions!")
}
