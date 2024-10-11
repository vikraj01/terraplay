package interactions

import (
    "fmt"
    "github.com/bwmarrin/discordgo"
)

func CreateInteraction() *discordgo.ApplicationCommand {
    return &discordgo.ApplicationCommand{
        Name:        "create",
        Description: "Create a new game server",
    }
}

func HandleCreateInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
    s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
        Type: discordgo.InteractionResponseModal,
        Data: &discordgo.InteractionResponseData{
            CustomID: "create_game_modal",
            Title:    "Create Game Server",
            Components: []discordgo.MessageComponent{
                discordgo.ActionsRow{
                    Components: []discordgo.MessageComponent{
                        discordgo.TextInput{
                            CustomID:    "server_name",
                            Label:       "Server Name",
                            Style:       discordgo.TextInputShort,
                            Placeholder: "Enter the name of the server",
                            Required:    true,
                        },
                    },
                },
                discordgo.ActionsRow{
                    Components: []discordgo.MessageComponent{
                        discordgo.TextInput{
                            CustomID:    "max_players",
                            Label:       "Maximum Players",
                            Style:       discordgo.TextInputShort,
                            Placeholder: "Enter max players (default is 10)",
                            Required:    false,
                        },
                    },
                },
            },
        },
    })
}

func HandleModalSubmit(s *discordgo.Session, i *discordgo.InteractionCreate) {
    if i.Type == discordgo.InteractionModalSubmit && i.ModalSubmitData().CustomID == "create_game_modal" {
        var serverName string
        var maxPlayers int64 = 10

        for _, component := range i.ModalSubmitData().Components {
            row := component.(*discordgo.ActionsRow)
            for _, input := range row.Components {
                textInput := input.(*discordgo.TextInput)
                switch textInput.CustomID {
                case "server_name":
                    serverName = textInput.Value
                case "max_players":
                    if textInput.Value != "" {
                        fmt.Sscanf(textInput.Value, "%d", &maxPlayers)
                    }
                }
            }
        }

        response := fmt.Sprintf("Creating server '%s' with max players: %d", serverName, maxPlayers)
        s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
            Type: discordgo.InteractionResponseChannelMessageWithSource,
            Data: &discordgo.InteractionResponseData{
                Content: response,
            },
        })
    }
}
