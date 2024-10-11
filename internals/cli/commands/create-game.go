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
)

var CreateGameCmd = &cobra.Command{
	Use:   "create-game",
	Short: "Create a game server",
	Long:  `Create a new game server session through the Zephyr server.`,
	Run:   createGame,
}

func createGame(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		log.Fatalf("Please provide a game name as an argument. Example: zephyr create-game minecraft")
	}

	gameName := args[0]

	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if config.AccessToken == "" {
		log.Fatalf("Access token not found. Please login first using 'zephyr login'.")
	}

	payload := map[string]string{
		"game":     gameName,
		"user_id":  config.UserID,
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

func loadConfig() (*Config, error) {
	configPath := os.ExpandEnv("$HOME/.zephyr/config.json")

	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config: %v", err)
	}

	return &config, nil
}

