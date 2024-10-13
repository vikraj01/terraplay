package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

var DestroySessionCmd = &cobra.Command{
	Use:   "destroy-session",
	Short: "Destroy a game session",
	Long:  `Initiate the destruction of a game session through the Zephyr server.`,
	Run:   destroySession,
}

func destroySession(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		log.Fatalf("Please provide a session ID as an argument. Example: zephyr destroy-session <session_id>")
	}

	sessionID := args[0]

	log.Print(sessionID)
	config, err := loadTokenConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if config.AccessToken == "" {
		log.Fatalf("Access token not found. Please login first using 'zephyr login'.")
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://localhost:8080/game/destroy/%s", sessionID), nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+config.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request to destroy session: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		log.Fatalf("Failed to destroy session. Status: %d, Response: %s", resp.StatusCode, string(bodyBytes))
	}

	var responseBody struct {
		Message string `json:"message"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		log.Fatalf("Failed to parse server response: %v", err)
	}

	fmt.Println(responseBody.Message)
}
