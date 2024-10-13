package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

var ListSessionsCmd = &cobra.Command{
	Use:   "list-sessions",
	Short: "List game sessions",
	Long:  `Retrieve a list of game sessions for the user with a specified status.`,
	Run:   listSessions,
}

func listSessions(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		log.Fatalf("Please provide a status as an argument. Example: zephyr list-sessions running")
	}

	status := args[0]

	config, err := loadTokenConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if config.AccessToken == "" {
		log.Fatalf("Access token not found. Please login first using 'zephyr login'.")
	}

	userID := config.UserID
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:8080/game/sessions?status=%s&user_id=%s", status, userID), nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+config.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request to list sessions: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		log.Fatalf("Failed to list sessions. Status: %d, Response: %s", resp.StatusCode, string(bodyBytes))
	}

	var responseBody struct {
		Message  string                   `json:"message"`
		Sessions []map[string]interface{} `json:"sessions"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		log.Fatalf("Failed to parse server response: %v", err)
	}

	fmt.Println("📋 Listing Sessions 📋")
	for _, session := range responseBody.Sessions {
		fmt.Printf(
			"Game: %s\nStatus: %s\nRun ID: %s\nStart Time: %s\nServer IP: %s\n",
			session["game"],
			session["status"],
			session["run_id"],
			session["start_time"],
			session["server_ip"],
		)
	}
}
