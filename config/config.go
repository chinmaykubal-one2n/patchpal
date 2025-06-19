package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	OpenRouterAPIKey string
	Model            string

	GithubToken string
	GithubOwner string
	GithubRepo  string

	RepoPath string
}

var Config AppConfig

func LoadEnv() {
	_ = godotenv.Load() // Load from .env if present

	Config = AppConfig{
		OpenRouterAPIKey: mustGet("OPENROUTER_API_KEY"),
		Model:            mustGet("MODEL"),
		GithubToken:      mustGet("GITHUB_TOKEN"),
		GithubOwner:      mustGet("GITHUB_OWNER"),
		GithubRepo:       mustGet("GITHUB_REPO"),
		RepoPath:         mustGet("REPO_PATH"),
	}
}

func mustGet(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Required environment variable %s not set", key)
	}
	return val
}
