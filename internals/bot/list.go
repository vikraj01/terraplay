package bot

import (
	"github.com/bwmarrin/discordgo"
)

// List game sessions command handler
func handleListSessionCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Listing all game sessions!")
}
