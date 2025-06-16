package main

import (
	"patchpal/api"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Register API route
	r.POST("/fix", api.HandleFix)

	// Run the server
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
