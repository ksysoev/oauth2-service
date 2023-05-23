package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig *oauth2.Config
	oauthStateString  = "random"
)

func main() {
	// Initialize the Google OAuth2 configuration
	googleOauthConfig = &oauth2.Config{
		ClientID:     "YOUR_CLIENT_ID",
		ClientSecret: "YOUR_CLIENT_SECRET",
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	// Initialize the Gin router
	router := gin.Default()

	// Load the HTML template
	router.LoadHTMLGlob("templates/*")

	// Define the routes
	router.GET("/", handleHome)
	router.GET("/auth/google/login", handleGoogleLogin)
	router.GET("/auth/google/callback", handleGoogleCallback)

	// Start the server
	router.Run(":8080")
}

// handleHome is the handler for the home page.
func handleHome(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", gin.H{})
}

// handleGoogleLogin is the handler for the Google OAuth2 login page.
func handleGoogleLogin(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// handleGoogleCallback is the handler for the Google OAuth2 callback page.
func handleGoogleCallback(c *gin.Context) {
	state := c.Query("state")
	if state != oauthStateString {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Invalid oauth state"))
		return
	}

	code := c.Query("code")
	token, err := googleOauthConfig.Exchange(c, code)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	client := googleOauthConfig.Client(c, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer resp.Body.Close()

	// TODO: Handle the user information returned by the API

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Logged in with Google",
	})
}
