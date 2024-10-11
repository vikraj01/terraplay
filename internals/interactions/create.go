package interactions

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func CreateInteraction() *discordgo.ApplicationCommand {
    return &discordgo.ApplicationCommand{
        Name:        "create",
        Description: "Create a new game server",
        Options: []*discordgo.ApplicationCommandOption{
            {
                Name:        "name",
                Description: "The name of the game server",
                Type:        discordgo.ApplicationCommandOptionString,
                Required:    true,
            },
            {
                Name:        "max_players",
                Description: "Maximum number of players",
                Type:        discordgo.ApplicationCommandOptionInteger,
                Required:    false,
            },
        },
    }
}

func HandleCreateInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
    options := i.ApplicationCommandData().Options
    serverName := options[0].StringValue()
    var maxPlayers int64 = 10

    if len(options) > 1 {
        maxPlayers = options[1].IntValue()
    }

    response := fmt.Sprintf("Creating server '%s' with max players: %d", serverName, maxPlayers)
    s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
        Type: discordgo.InteractionResponseChannelMessageWithSource,
        Data: &discordgo.InteractionResponseData{
            Content: response,
        },
    })
}

