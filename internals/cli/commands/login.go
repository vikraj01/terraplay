package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/spf13/cobra"
)

var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Zephyr using Discord",
	Long:  `Login to Zephyr using Discord's OAuth to authenticate and receive an access token.`,
	Run:   login,
}

func login(cmd *cobra.Command, args []string) {
	fmt.Println("Initiating login flow...")

	response, err := http.Get("http://localhost:8080/auth/discord/initiate")
	if err != nil {
		log.Fatalf("Failed to initiate OAuth: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Fatalf("Failed to initiate OAuth, status code: %d", response.StatusCode)
	}

	var responseBody struct {
		AuthURL string `json:"auth_url"`
	}
	if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
		log.Fatalf("Failed to parse OAuth initiation response: %v", err)
	}

	err = openBrowser(responseBody.AuthURL)
	if err != nil {
		log.Fatalf("Failed to open browser: %v", err)
	}

	fmt.Println("Please complete the login in your browser...")
	waitForAuthCompletion()
}

func openBrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	}
	return err
}

const (
	maxRetries     = 5
	baseRetryDelay = 2 * time.Second
	authCheckURL   = "http://localhost:8080/auth/discord/status"
	authTimeout    = 3 * time.Minute
)

func waitForAuthCompletion() {
	fmt.Println("Waiting for authentication to complete...")
	timeout := time.After(authTimeout)
	retries := 0

	for {
		select {
		case <-timeout:
			log.Fatalf("Authentication timed out after %s. Please try again.", authTimeout)
		default:
			if retries >= maxRetries {
				log.Fatalf("Max retries reached. Could not verify authentication status.")
			}

			response, err := http.Get(authCheckURL)
			if err != nil {
				log.Printf("Error checking authentication status: %v. Retrying...", err)
				retries++
				time.Sleep(calculateBackoff(retries))
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
					log.Printf("Error parsing authentication response: %v. Retrying...", err)
					retries++
					time.Sleep(calculateBackoff(retries))
					continue
				}

				if err := saveTokenLocally(responseBody.AccessToken, responseBody.UserID, responseBody.Username); err != nil {
					log.Fatalf("Failed to save access token: %v", err)
				}

				fmt.Println("Login successful! Access token and user details saved.")
				return
			} else if response.StatusCode == http.StatusUnauthorized {
				log.Println("Authentication not yet completed. Please approve the request in your browser.")
				retries++
				time.Sleep(calculateBackoff(retries))
			} else {
				log.Printf("Unexpected status code %d from server. Retrying...", response.StatusCode)
				retries++
				time.Sleep(calculateBackoff(retries))
			}
		}
	}
}

func calculateBackoff(retry int) time.Duration {
	return time.Duration(retry) * baseRetryDelay
}

type Config struct {
	AccessToken string `json:"access_token"`
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
}

func saveTokenLocally(token, userID, username string) error {
	configPath := os.ExpandEnv("$HOME/.zephyr/config.json")

	configData := Config{
		AccessToken: token,
		UserID:      userID,
		Username:    username,
	}

	file, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %v", err)
	}
	defer file.Close()

	log.Print(configData)
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(configData); err != nil {
		return fmt.Errorf("failed to write to config file: %v", err)
	}

	return nil
}
