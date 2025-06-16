package main

import (
	"patchpal/api"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load env variables if .env file exists
	_ = godotenv.Load()

	r := gin.Default()

	// Register API route
	r.POST("/fix", api.HandleFix)

	// Run the server
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
