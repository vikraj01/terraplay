package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/vikraj01/terraplay/internals/github"
	"github.com/vikraj01/terraplay/internals/utils"
)

var commandMap = map[string]func(*discordgo.Session, *discordgo.MessageCreate){
	"!ping":         handlePingCommand,
	"!create":       handleCreateCommand,
	"!destroy":      handleDestroyCommand,
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
	args := strings.Fields(m.Content)
	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Usage: !create <game>")
		return
	}
	gameName := args[2]

	inputs := map[string]string{
		"game":    gameName,
		"user_id": m.Author.ID,
		"run_id":  utils.GenerateUUID(),
	}
	fmt.Print(inputs)

	err := github.TriggerGithubAction("vikraj01", "terraplay", "start.game.yml", "main", inputs)
	fmt.Print(err)
	if err != nil {
		fmt.Println(err)
		s.ChannelMessageSend(m.ChannelID, "Failed to trigger GitHub Action to create game session!")
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Game session created! GitHub Action triggered for game: "+gameName)
}

func handleDestroyCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Game session destroyed!")
}

func handleListSessionCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Listing all game sessions!")
}
