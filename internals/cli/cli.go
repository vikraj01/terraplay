package cli

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

func StartCli() {
	var rootCmd = &cobra.Command{
		Use:   "zephyr",
		Short: "Zephyr - Manage game servers via CLI",
		Long:  `Zephyr is a CLI tool to manage game servers through Terraplay.`,
	}
	RegisterCommands(rootCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error starting Zephyr CLI: %v", err)
		os.Exit(1)
	}
}
