package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	OpenRouterAPIKey string
	Model            string
	GithubToken      string
	GithubOwner      string
	GithubRepo       string
	RepoPath         string // Local path to the repository
}

var Config AppConfig

// LoadEnv loads environment variables into the Config struct
func LoadEnv() {
	log.Println("[Config] Loading environment variables from .env...")
	_ = godotenv.Load()

	Config = AppConfig{
		OpenRouterAPIKey: mustGet("OPENROUTER_API_KEY"),
		Model:            mustGet("MODEL"),
		GithubToken:      mustGet("GITHUB_TOKEN"),
		GithubOwner:      mustGet("GITHUB_OWNER"),
		GithubRepo:       mustGet("GITHUB_REPO"),
		RepoPath:         mustGet("REPO_PATH"),
	}
	log.Println("[Config] Environment variables loaded successfully.")
}

// mustGet retrieves an environment variable
func mustGet(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("[Config] Required environment variable %s not set", key)
	}
	log.Printf("[Config] Loaded %s\n", key)
	return val
}
