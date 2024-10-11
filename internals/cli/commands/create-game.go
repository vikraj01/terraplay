package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/vikraj01/terraplay/config"
	"github.com/vikraj01/terraplay/internals/utils"
)

var CreateGameCmd = &cobra.Command{
	Use:   string(config.CreateGame),
	Short: "Create a game server",
	Long:  `Create a new game server session through the Zephyr server.`,
	Run:   createGame,
}

func createGame(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		log.Fatalf("Please provide a game name as an argument. Example: zephyr create-game minecraft")
	}

	gameName := args[0]

	config, err := loadTokenConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if config.AccessToken == "" {
		log.Fatalf("Access token not found. Please login first using 'zephyr login'.")
	}

	payload := map[string]string{
		"game":      gameName,
		"user_id":   config.UserID,
		"user_name": config.Username,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Failed to marshal request payload: %v", err)
	}

	req, err := http.NewRequest("POST", "http://localhost:8080/game/create", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+config.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request to create game: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		log.Fatalf("Failed to create game. Status: %d, Response: %s", resp.StatusCode, string(bodyBytes))
	}

	var responseBody struct {
		Message  string `json:"message"`
		ServerIP string `json:"server_ip"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		log.Fatalf("Failed to parse server response: %v", err)
	}

	fmt.Printf("Game server created successfully! Server IP: %s\n", responseBody.ServerIP)
}

func loadTokenConfig() (*Config, error) {
	var config Config
	baseDir := os.Getenv("HOME") + "/.zephyr"
	if err := utils.LoadConfig(baseDir, "config.json", &config); err != nil {
		return nil, fmt.Errorf("failed to load configuration: %v", err)
	}

	log.Printf("Loaded configuration: %v", config)
	return &config, nil
}
