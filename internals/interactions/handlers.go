package interactions

import "github.com/bwmarrin/discordgo"

func InteractionHandlers(s *discordgo.Session, i *discordgo.InteractionCreate) {
    if i.Type == discordgo.InteractionApplicationCommand {
        switch i.ApplicationCommandData().Name {
        case "create":
            HandleCreateInteraction(s, i)
        }
    } else if i.Type == discordgo.InteractionModalSubmit {
        switch i.ModalSubmitData().CustomID {
        case "create_game_modal":
            HandleModalSubmit(s, i)
        }
    }
}


