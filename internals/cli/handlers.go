package cli

import (
	"github.com/spf13/cobra"
	"github.com/vikraj01/terraplay/internals/cli/commands"
)

func RegisterCommands(rootCmd *cobra.Command) {
	rootCmd.AddCommand(commands.LoginCmd)
	rootCmd.AddCommand(commands.CreateGameCmd)
}
