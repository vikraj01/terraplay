package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

var ListGamesCmd = &cobra.Command{
	Use:   "list-games",
	Short: "List all available games",
	Long:  `Retrieve a list of all available games from the Zephyr server.`,
	Run:   listGames,
}

func listGames(cmd *cobra.Command, args []string) {
	req, err := http.NewRequest("GET", "http://localhost:8080/game/list", nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request to list games: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		log.Fatalf("Failed to list games. Status: %d, Response: %s", resp.StatusCode, string(bodyBytes))
	}

	var responseBody struct {
		Message string   `json:"message"`
		Games   []string `json:"games"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		log.Fatalf("Failed to parse server response: %v", err)
	}

	fmt.Println("🎮 Available Games 🎮")
	for _, game := range responseBody.Games {
		fmt.Printf("- %s\n", game)
	}
}
