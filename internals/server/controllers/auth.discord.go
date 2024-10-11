package controllers

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "net/url"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/go-resty/resty/v2"
)

var (
    redirectURI = "http://localhost:8080/auth/discord/callback"
    accessToken = ""
    userID      = ""
    username    = ""
)

func InitiateDiscordOAuth(c *gin.Context) {
    clientID := os.Getenv("DISCORD_CLIENT_ID")
    clientSecret := os.Getenv("DISCORD_CLIENT_SECRET")
    log.Printf("InitiateDiscordOAuth called with clientID: %s, clientSecret: %s", clientID, clientSecret)

    if clientID == "" || clientSecret == "" {
        log.Println("Missing Discord OAuth configuration in InitiateDiscordOAuth")
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Missing Discord OAuth configuration"})
        return
    }

    authURL := fmt.Sprintf(
        "https://discord.com/api/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=identify",
        clientID, url.QueryEscape(redirectURI),
    )

    c.JSON(http.StatusOK, gin.H{"auth_url": authURL})
}

func DiscordOAuthCallback(c *gin.Context) {
    code := c.Query("code")
    if code == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code is missing"})
        return
    }

    clientID := os.Getenv("DISCORD_CLIENT_ID")
    clientSecret := os.Getenv("DISCORD_CLIENT_SECRET")

    tokenURL := "https://discord.com/api/oauth2/token"
    client := resty.New()

    response, err := client.R().
        SetFormData(map[string]string{
            "client_id":     clientID,
            "client_secret": clientSecret,
            "grant_type":    "authorization_code",
            "code":          code,
            "redirect_uri":  redirectURI,
        }).
        SetHeader("Content-Type", "application/x-www-form-urlencoded").
        Post(tokenURL)

    if err != nil || response.StatusCode() != http.StatusOK {
        log.Printf("Error exchanging code for token: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange code for token"})
        return
    }

    var tokenResponse struct {
        AccessToken string `json:"access_token"`
    }
    if err := json.Unmarshal(response.Body(), &tokenResponse); err != nil {
        log.Printf("Failed to parse token response: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse token response"})
        return
    }

    accessToken = tokenResponse.AccessToken

    userInfoURL := "https://discord.com/api/v10/users/@me"
    userResponse, err := client.R().
        SetAuthToken(accessToken).
        Get(userInfoURL)

    if err != nil || userResponse.StatusCode() != http.StatusOK {
        log.Printf("Error fetching user details: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user details"})
        return
    }

    var userInfo struct {
        ID       string `json:"id"`
        Username string `json:"username"`
    }

    if err := json.Unmarshal(userResponse.Body(), &userInfo); err != nil {
        log.Printf("Failed to parse user info response: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user info"})
        return
    }

    // Store user details
    userID = userInfo.ID
    username = userInfo.Username

    c.JSON(http.StatusOK, gin.H{
        "message":      "Login successful",
        "access_token": accessToken,
        "user_id":      userID,
        "username":     username,
    })
}

func CheckAuthStatus(c *gin.Context) {
    if accessToken == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"message": "User not authenticated"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "access_token": accessToken,
        "user_id":      userID,
        "username":     username,
    })
}
