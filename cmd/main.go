package main

import (
	"fmt"
	"patchpal/config"
	"patchpal/service"
)

func main() {
	// Load environment variables into config
	config.LoadEnv()

	if err := service.ProcessAndFixManifests(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
}
