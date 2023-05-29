package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig *oauth2.Config
	oauthStateString  = "random"
)

var mongoClient *mongo.Client

func main() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://root:example@localhost:27017")

	var err error
	// Connect to MongoDB
	mongoClient, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		fmt.Println("Failed to connect to MongoDB:", err)
		return
	}

	// Check the connection
	err = mongoClient.Ping(context.Background(), nil)
	if err != nil {
		fmt.Println("Failed to ping MongoDB:", err)
		return
	}

	fmt.Println("Connected to MongoDB!")
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
	router.POST("/registration", handleProcessRegistration)

	// Start the server
	router.Run(":8080")
}

// handleHome is the handler for the home page.
func handleHome(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", gin.H{})
}

func handleRegistration(c *gin.Context) {
	c.HTML(http.StatusOK, "registration.html", gin.H{})
}

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func handleProcessRegistration(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := mongoClient.Ping(context.Background(), nil)
	if err != nil {
		fmt.Println("Failed to ping MongoDB:", err)
		return
	}

	res, err := mongoClient.Database("oauth2").Collection("users").InsertOne(context.Background(), user)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(res)

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})

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
