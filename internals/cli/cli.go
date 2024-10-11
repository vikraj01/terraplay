package cli

import (
	"log"
	"os"
	"strings"

	"github.com/mgutz/ansi"
	"github.com/spf13/cobra"
)

func StartCli() {
	var rootCmd = &cobra.Command{
		Use:   "zephyr",
		Short: "Zephyr - Manage game servers via CLI",
		Long:  getWelcomeMessage(),
	}

	rootCmd.SetHelpTemplate(`{{.Long}}`)

	RegisterCommands(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error starting Zephyr CLI: %v", err)
		os.Exit(1)
	}
}

func getWelcomeMessage() string {
	colors := []string{"cyan+h", "green+h", "yellow+h", "magenta+h", "blue+h"}

	asciiArtLines := []string{
		"       ___           ___           ___         ___                       ___     ",
		"     /  /\\         /  /\\         /  /\\       /__/\\          ___        /  /\\    ",
		"    /  /::|       /  /:/_       /  /::\\      \\  \\:\\        /__/|      /  /::\\   ",
		"   /  /:/:|      /  /:/ /\\     /  /:/\\:\\      \\__\\:\\      |  |:|     /  /:/\\:\\  ",
		"  /  /:/|:|__   /  /:/ /:/_   /  /:/~/:/  ___ /  /::\\     |  |:|    /  /:/~/:/  ",
		" /__/:/ |:| /\\ /__/:/ /:/ /\\ /__/:/ /:/  /__/\\  /:/\\:\\  __|__|:|   /__/:/ /:/___",
		" \\__\\/  |:|/:/ \\  \\:\\/:/ /:/ \\  \\:\\/:/   \\  \\:\\/:/__\\/ /__/::::\\   \\  \\:\\/:::::/",
		"     |  |:/:/   \\  \\::/ /:/   \\  \\::/     \\  \\::/         ~\\~~\\:\\   \\  \\::/~~~~ ",
		"     |  |::/     \\  \\:\\/:/     \\  \\:\\      \\  \\:\\           \\  \\:\\   \\  \\:\\     ",
		"     |  |:/       \\  \\::/       \\  \\:\\      \\  \\:\\           \\__\\/    \\  \\:\\    ",
		"     |__|/         \\__\\/         \\__\\/       \\__\\/                     \\__\\/    ",
	}

	var coloredArtLines []string
	for i, line := range asciiArtLines {
		colorFunc := ansi.ColorFunc(colors[i%len(colors)])
		coloredArtLines = append(coloredArtLines, colorFunc(line))
	}
	asciiArt := strings.Join(coloredArtLines, "\n")

	description := `
` + ansi.Color("Welcome to Zephyr!", "magenta+b") + `
` + ansi.Color("Zephyr", "white+b") + ` is a CLI tool designed for managing game servers through Terraplay.

` + ansi.Color("Commands:", "yellow+b") + `
  ` + ansi.Color("create-game", "green") + `   Create a game server
  ` + ansi.Color("login", "blue") + `         Authenticate with Discord
`
	return asciiArt + description
}
