package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/vikraj01/terraplay/internals/interactions"
	// "github.com/joho/godotenv"
)

func StartBot() {
	// if os.Getenv("APP_ENV") != "production" {
	// 	err := godotenv.Load()
	// 	if err != nil {
	// 		log.Println("Error loading .env file, using environment variables")
	// 	}
	// }

	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		log.Fatalf("Discord bot token is missing")
	}

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	session.AddHandler(messageHandler)
	session.AddHandler(interactions.InteractionHandlers)

	session.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages



	err = session.Open()
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
		return
	}

	interactions.RegisterInteractionCommands(session)
	fmt.Println("Bot is now running. Press CTRL+C to exit.")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	session.Close()
}

// # For First Phase I want to achieve these
// !create game
// !destroy <sessionid>
// !list-session

// # For Next Phase
// !stop
// !restart

// # Next Phase
// !logs
// !cost
// !config - config permission for other team member to handle server permission


// V2

// 1. Setup Interaction
// 2. Add cli
// 3. Config Power To User For many different things
// 3. Cosmetics Heavy UI / UX
// 3. Expenses & Logs with cloudwatch [ many different sources ]
// 4. Update config Option - For Everything Possible!
// 4. One Health Manager, who do workspace cleanup and pending cleanups + Check global infra health! + Global infra spending
// 5. Tightening Security - Phase 1