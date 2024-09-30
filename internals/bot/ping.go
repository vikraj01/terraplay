package bot

import (
	"github.com/bwmarrin/discordgo"
)

func handlePingCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Pong!")
}
