package interactions

import "github.com/bwmarrin/discordgo"

func InteractionHandlers(s *discordgo.Session, i *discordgo.InteractionCreate) {
        if i.Type == discordgo.InteractionApplicationCommand {
                switch i.ApplicationCommandData().Name {
                case "create":
                        HandleCreateInteraction(s, i)
                }
        }
}
