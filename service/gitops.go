package service

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"patchpal/config"
	"path/filepath"
	"time"

	"github.com/google/go-github/v57/github"
	"golang.org/x/oauth2"
)

func CreatePRWithFixedYAML(fixedFilePath, originalFilename string) (string, error) {
	var (
		repoPath    = config.Config.RepoPath
		githubOwner = config.Config.GithubOwner
		githubRepo  = config.Config.GithubRepo
		githubToken = config.Config.GithubToken
	)

	ctx := context.Background()
	branchName := fmt.Sprintf("patchpal-fix-%d", time.Now().Unix())

	// Step 1: Create new branch
	if err := runGit(repoPath, "checkout", "-b", branchName); err != nil {
		return "", fmt.Errorf("failed to checkout new branch: %w", err)
	}

	// Step 2: Replace the file in repo
	destPath := filepath.Join(repoPath, originalFilename)
	srcData, err := os.ReadFile(fixedFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read fixed file: %w", err)
	}
	if err := os.WriteFile(destPath, srcData, 0644); err != nil {
		return "", fmt.Errorf("failed to overwrite repo file: %w", err)
	}

	// Step 3: Add & Commit
	if err := runGit(repoPath, "add", originalFilename); err != nil {
		return "", fmt.Errorf("git add failed: %w", err)
	}
	if err := runGit(repoPath, "commit", "-m", "fix: patch Kubernetes misconfig"); err != nil {
		return "", fmt.Errorf("git commit failed: %w", err)
	}

	// Step 4: Push
	if err := runGit(repoPath, "push", "-u", "origin", branchName); err != nil {
		return "", fmt.Errorf("git push failed: %w", err)
	}

	// Step 5: Create PR using GitHub API
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
	client := github.NewClient(oauth2.NewClient(ctx, ts))

	title := "PatchPal: Fix Kubernetes Misconfiguration"
	body := "This PR contains fixes generated automatically by PatchPal for a vulnerable manifest."
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
