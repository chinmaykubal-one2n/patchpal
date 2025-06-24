package main

import (
	"fmt"
	"patchpal/config"
	"patchpal/service"
)

func main() {
	// Load environment variables into config
	config.LoadEnv()

	// Process and fix k8s manifests
	if err := service.ProcessAndFixManifests(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
}
