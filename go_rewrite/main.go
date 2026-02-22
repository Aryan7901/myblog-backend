package main

import (
	"backend/config"
	"backend/routes"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	tokenSecret := os.Getenv("TOKEN_SECRET")
	if tokenSecret == "" {
		log.Fatal("TOKEN_SECRET environment variable not set")
	}

	dbURI := os.Getenv("db")
	if dbURI == "" {
		log.Fatal("db environment variable not set")
	}

	err := config.ConnectDB(dbURI)
	if err != nil {
		log.Fatal("Could not connect to database")
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://aryan7901.github.io", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.UserRoutes(router)
	routes.BlogRoutes(router)

	router.Use(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Could not find this route."})
	})

	router.Use(func(c *gin.Context) {
		c.Next()
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	log.Printf("Server running on port %s", port)
	router.Run(":" + port)
}
