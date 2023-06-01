package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/ksysoev/oauth2-service/pkg/services"
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

	// Initialize the Redis store for sessions
	store, err := redis.NewStore(5, "tcp", "localhost:6379", "", []byte("secret"))
	if err != nil {
		log.Fatal(err)
	}
	router.Use(sessions.Sessions("mysession", store))

	router.Use(csrf.Middleware(csrf.Options{
		Secret: "secret123",
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	}))

	// Load the HTML template
	router.LoadHTMLGlob("templates/*")

	// Define the routes
	router.GET("/", handleHome)
	router.POST("/signup", SignUp)
	router.GET("/signup", SignUp)

	router.POST("/signin", SignIn)
	router.GET("/signin", SignIn)

	// Start the server
	router.Run(":8080")
}

// handleHome is the handler for the home page.
func handleHome(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/signin")
}

func SignUp(c *gin.Context) {
	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "signup.tmpl", gin.H{})
		return
	}

	var user services.SignUpRequest
	if err := c.ShouldBind(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.SignUp(c, &user, mongoClient)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "home.tmpl", gin.H{})
	return
}

func SignIn(c *gin.Context) {
	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "signin.tmpl", gin.H{})
		return
	}

	var signInRequest services.SignInRequest
	if err := c.ShouldBind(&signInRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.SignIn(c, &signInRequest, mongoClient)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "home.tmpl", gin.H{})
	return
}
