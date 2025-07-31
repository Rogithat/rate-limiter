package main

import (
	"log"
	"os"

	"rate-limiter/limiter"
	"rate-limiter/middleware"
	"rate-limiter/storage"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load("config.env"); err != nil {
		log.Fatal("Error loading config.env file")
	}

	// Initialize Redis storage
	redisStorage, err := storage.NewRedisStorage()
	if err != nil {
		log.Fatalf("Failed to initialize Redis storage: %v", err)
	}

	// Initialize rate limiter
	config := limiter.NewConfig()
	rateLimiter := limiter.NewRateLimiter(redisStorage, config)

	// Initialize Gin router
	router := gin.Default()

	// Apply rate limiter middleware
	router.Use(middleware.RateLimitMiddleware(rateLimiter))

	// Add a test endpoint
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Rate limit test successful",
		})
	})

	// Start the server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
