package bot

import (
	"github.com/bwmarrin/discordgo"
)

func handleDestroyCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Game session destroyed!")
}
