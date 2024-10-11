package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/vikraj01/terraplay/config"
)

var RestartGameCmd = &cobra.Command{
	Use:   string(config.RestartGame),
	Short: "Restart a game server",
	Long:  `Restart an existing game server session through the Zephyr server.`,
	Run:   restartGame,
}

func restartGame(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		log.Fatalf("Please provide a session ID as an argument. Example: zephyr restart-game <session_id>")
	}

	sessionID := args[0]

	config, err := loadTokenConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if config.AccessToken == "" {
		log.Fatalf("Access token not found. Please login first using 'zephyr login'.")
	}

	payload := map[string]string{
		"session_id": sessionID,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Failed to marshal request payload: %v", err)
	}

	req, err := http.NewRequest("POST", "http://localhost:8080/game/restart", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+config.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request to restart game: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		log.Fatalf("Failed to restart game. Status: %d, Response: %s", resp.StatusCode, string(bodyBytes))
	}

	var responseBody struct {
		Message  string `json:"message"`
		ServerIP string `json:"server_ip"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		log.Fatalf("Failed to parse server response: %v", err)
	}

	fmt.Printf("✅ Game server has been restarted successfully! Server IP: %s\n", responseBody.ServerIP)
}
