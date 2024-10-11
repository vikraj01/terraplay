package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" || len(token) < 8 || token[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing access token"})
			c.Abort()
			return
		}

		accessToken := token[7:]

		userID, err := validateDiscordToken(accessToken)
		if err != nil {
			log.Printf("Failed to validate access token: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid access token"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}

func validateDiscordToken(accessToken string) (string, error) {
	client := resty.New()
	resp, err := client.R().
		SetAuthToken(accessToken).
		Get("https://discord.com/api/users/@me")

	if err != nil {
		return "", err
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("discord token validation failed, status code: %d", resp.StatusCode())
	}

	var result struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return "", fmt.Errorf("failed to parse Discord user response: %v", err)
	}

	return result.ID, nil
}
