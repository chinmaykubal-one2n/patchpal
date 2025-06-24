package main

import (
	"log"
	"patchpal/config"
	"patchpal/service"
)

func main() {
	// Load environment variables and configuration
	config.LoadEnv()

	log.Println("[PatchPal] Starting manifest processing and fixing...")
	if err := service.ProcessAndFixManifests(); err != nil {
		log.Fatalf("[PatchPal] Error: %v\n", err)
	}
	log.Println("[PatchPal] All manifests processed successfully.")
}
