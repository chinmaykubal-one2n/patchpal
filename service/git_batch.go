package service

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"patchpal/config"

	"github.com/google/go-github/v57/github"
	"golang.org/x/oauth2"
)

func CreateBatchPR(fixedFiles []string) (string, error) {
	repoPath := config.Config.RepoPath
	githubOwner := config.Config.GithubOwner
	githubRepo := config.Config.GithubRepo
	githubToken := config.Config.GithubToken

	ctx := context.Background()
	branchName := fmt.Sprintf("patchpal-fix-%d", time.Now().Unix())

	// Step 1: Create new branch
	if err := runGit(repoPath, "checkout", "-b", branchName); err != nil {
		return "", fmt.Errorf("failed to checkout new branch: %w", err)
	}

	// Step 2: Add all fixed files
	for _, file := range fixedFiles {
		relPath, err := filepath.Rel(repoPath, file)
		if err != nil {
			return "", fmt.Errorf("failed to get relative path: %w", err)
		}
		if err := runGit(repoPath, "add", relPath); err != nil {
			return "", fmt.Errorf("failed to add file: %s, error: %w", relPath, err)
		}
	}

	// Step 3: Commit once
	if err := runGit(repoPath, "commit", "-m", "fix: patch Kubernetes misconfigurations"); err != nil {
		return "", fmt.Errorf("git commit failed: %w", err)
	}

	// Step 4: Set remote with token
	remoteURL := fmt.Sprintf("https://x-access-token:%s@github.com/%s/%s.git", githubToken, githubOwner, githubRepo)
	if err := runGit(repoPath, "remote", "set-url", "origin", remoteURL); err != nil {
		return "", fmt.Errorf("failed to set remote URL: %w", err)
	}

	// Step 5: Push branch
	if err := runGit(repoPath, "push", "-u", "origin", branchName); err != nil {
		return "", fmt.Errorf("git push failed: %w", err)
	}

	// Step 6: Create PR via GitHub API
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
	client := github.NewClient(oauth2.NewClient(ctx, ts))

	title := "PatchPal: Fix Kubernetes Misconfigurations"
	body := "This PR contains automated fixes for Kubernetes security misconfigurations found in multiple files."

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
