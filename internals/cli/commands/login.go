package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/vikraj01/terraplay/config"
	"github.com/vikraj01/terraplay/internals/utils"
)

type Config struct {
	AccessToken string `json:"access_token"`
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
}

const (
	maxRetries     = 5
	authCheckURL   = "http://localhost:8080/auth/discord/status"
	authTimeout    = 3 * time.Minute
)

var LoginCmd = &cobra.Command{
	Use:   string(config.Login),
	Short: "Login to Zephyr using Discord",
	Long:  `Login to Zephyr using Discord's OAuth to authenticate and receive an access token.`,
	Run:   login,
}

func login(cmd *cobra.Command, args []string) {
	fmt.Println("🚀 Initiating login flow. Please wait...")

	response, err := http.Get("http://localhost:8080/auth/discord/initiate")
	if err != nil {
		log.Fatalf("❌ Failed to initiate OAuth: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Fatalf("❌ Failed to initiate OAuth, status code: %d", response.StatusCode)
	}

	var responseBody struct {
		AuthURL string `json:"auth_url"`
	}
	if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
		log.Fatalf("❌ Failed to parse OAuth initiation response: %v", err)
	}

	err = utils.OpenBrowser(responseBody.AuthURL)
	if err != nil {
		log.Fatalf("❌ Failed to open browser: %v", err)
	}

	fmt.Println("🌐 Please complete the login in your browser...")
	waitForAuthCompletion()
}

func waitForAuthCompletion() {
	progressLogger := utils.NewProgressLogger(maxRetries, authTimeout)
	fmt.Println("⌛ Waiting for authentication to complete...")

	timeout := time.After(authTimeout)
	retries := 0

	for {
		select {
		case <-timeout:
			progressLogger.LogCompletion(false, "Authentication timed out")
			log.Fatalf("❌ Authentication timed out after %s. Please try again.", authTimeout)
		default:
			if retries >= maxRetries {
				progressLogger.LogCompletion(false, "Max retries reached")
				log.Fatalf("❌ Max retries reached. Could not verify authentication status.")
			}

			response, err := http.Get(authCheckURL)
			if err != nil {
				progressLogger.LogProgress(retries, "Error checking authentication status")
				retries++
				time.Sleep(utils.CalculateBackoff(retries))
				continue
			}
			defer response.Body.Close()

			if response.StatusCode == http.StatusOK {
				var responseBody struct {
					AccessToken string `json:"access_token"`
					UserID      string `json:"user_id"`
					Username    string `json:"username"`
				}
				if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
					progressLogger.LogProgress(retries, "Error parsing authentication response")
					retries++
					time.Sleep(utils.CalculateBackoff(retries))
					continue
				}

				if err := saveTokenLocally(responseBody.AccessToken, responseBody.UserID, responseBody.Username); err != nil {
					progressLogger.LogCompletion(false, "Failed to save access token")
					log.Fatalf("❌ Failed to save access token: %v", err)
				}

				progressLogger.LogCompletion(true, "Login successful")
				fmt.Println("✅ Login successful! Access token and user details saved.")
				return
			} else if response.StatusCode == http.StatusUnauthorized {
				progressLogger.LogProgress(retries, "Authentication not yet completed")
				retries++
				time.Sleep(utils.CalculateBackoff(retries))
			} else {
				progressLogger.LogProgress(retries, fmt.Sprintf("Unexpected status code %d", response.StatusCode))
				retries++
				time.Sleep(utils.CalculateBackoff(retries))
			}
		}
	}
}

func saveTokenLocally(token, userID, username string) error {
	config := Config{
		AccessToken: token,
		UserID:      userID,
		Username:    username,
	}

	baseDir := os.Getenv("HOME") + "/.zephyr"
	if err := utils.SaveConfig(config, baseDir, "config.json"); err != nil {
		return fmt.Errorf("failed to save configuration: %v", err)
	}

	log.Printf("💾 Configuration saved successfully")
	return nil
}
