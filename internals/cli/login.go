package cli

import (
    "fmt"
    "github.com/spf13/cobra"
)

var LoginCmd = &cobra.Command{
    Use:   "login",
    Short: "Login to Zephyr using Discord",
    Long:  `Login to Zephyr using Discord's OAuth to authenticate and receive an access token.`,
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Placeholder: Initiating login flow...")
        fmt.Println("Login flow complete. Token saved to your machine.")
    },
}
