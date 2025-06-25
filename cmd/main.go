package main

import (
	"context"
	"log"
	"patchpal/config"
	"patchpal/service"
)

func main() {
	//  Initialize context
	ctx := context.Background()

	// Load environment variables and configuration
	config.LoadEnv()

	log.Println("[PatchPal] Starting manifest processing and fixing...")
	if err := service.ProcessAndFixManifests(ctx); err != nil {
		log.Fatalf("[PatchPal] Error: %v\n", err)
	}
	log.Println("[PatchPal] All manifests processed successfully.")
}
