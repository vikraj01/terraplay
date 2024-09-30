package bot

import (
	"github.com/bwmarrin/discordgo"
)

func handleListSessionCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Listing all game sessions!")
}

func handleListGamesCommand(s *discordgo.Session, m *discordgo.MessageCreate) {

	gameList := "**🎮 Available Games 🎮**\n\n"
	for _, game := range validGames {
		gameList += game + "\n"
	}
	gameList += "\nType `!create <game_name>` to create a game server!"

	s.ChannelMessageSend(m.ChannelID, gameList)
}
