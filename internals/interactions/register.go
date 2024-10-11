package interactions

import (
    "log"
    "github.com/bwmarrin/discordgo"
)

func RegisterInteractionCommands(s *discordgo.Session) {
    _, err := s.ApplicationCommandCreate(s.State.User.ID, "", CreateInteraction())
    if err != nil {
        log.Fatalf("Cannot create 'create' command: %v", err)
    }

    log.Println("Slash commands registered successfully.")
}
