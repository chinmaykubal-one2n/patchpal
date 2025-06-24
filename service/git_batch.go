package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"patchpal/config"

	"github.com/google/go-github/v57/github"
	"golang.org/x/oauth2"
)

// CreateBatchPR creates a new branch, commits fixed files, pushes the branch, and opens a pull request on GitHub
// It returns the URL of the created pull request
func CreateBatchPR(fixedFiles []string) (string, error) {
	repoPath := config.Config.RepoPath
	githubOwner := config.Config.GithubOwner
	githubRepo := config.Config.GithubRepo
	githubToken := config.Config.GithubToken

	ctx := context.Background()
	branchName := fmt.Sprintf("patchpal-fix-%d", time.Now().Unix())

	// Step 1: Create new branch
	log.Printf("[GitBatch] Creating new branch: %s\n", branchName)
	if err := runGit(repoPath, "checkout", "-b", branchName); err != nil {
		return "", fmt.Errorf("failed to checkout new branch: %w", err)
	}

	// Step 2: Add all fixed files to the new branch
	for _, file := range fixedFiles {
		relPath, err := filepath.Rel(repoPath, file)
		if err != nil {
			return "", fmt.Errorf("failed to get relative path: %w", err)
		}
		log.Printf("[GitBatch] Adding file to commit: %s\n", relPath)
		if err := runGit(repoPath, "add", relPath); err != nil {
			return "", fmt.Errorf("failed to add file: %s, error: %w", relPath, err)
		}
	}

	// Step 3: Commit all changes in one commit
	log.Println("[GitBatch] Committing changes...")
	if err := runGit(repoPath, "commit", "-m", "fix: patch Kubernetes misconfigurations"); err != nil {
		return "", fmt.Errorf("git commit failed: %w", err)
	}

	// Step 4: Push the new branch to origin
	log.Println("[GitBatch] Pushing branch to origin...")
	if err := runGit(repoPath, "push", "-u", "origin", branchName); err != nil {
		return "", fmt.Errorf("git push failed: %w", err)
	}

	// Step 5: Create a pull request using the GitHub API
	log.Println("[GitBatch] Creating pull request via GitHub API...")
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
	client := github.NewClient(oauth2.NewClient(ctx, ts))

	title := "PatchPal: Fix Kubernetes Misconfigurations"
	body := "This PR contains automated fixes for Kubernetes security misconfigurations."

	pr := &github.NewPullRequest{
		Title: &title,
		Head:  &branchName,
		Base:  github.String("dev"),
		Body:  &body,
	}

	createdPR, _, err := client.PullRequests.Create(ctx, githubOwner, githubRepo, pr)
	if err != nil {
		return "", fmt.Errorf("failed to create GitHub PR: %w", err)
	}

	return createdPR.GetHTMLURL(), nil
}

func runGit(repoPath string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = repoPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
