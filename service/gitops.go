package service

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"patchpal/config"
	"path/filepath"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/go-github/v57/github"
	"golang.org/x/oauth2"
)

var (
	repoPath       = config.Config.RepoPath
	githubOwner    = config.Config.GithubOwner
	githubRepo     = config.Config.GithubRepo
	githubToken    = config.Config.GithubToken
	commitAuthor   = &object.Signature{Name: "PatchPal Bot", Email: "bot@patchpal.io", When: time.Now()}
	prBranchPrefix = "patchpal-fix"
)

func CreatePRWithFixedYAML(fixedFilePath, originalFilename string) (string, error) {
	ctx := context.Background()
	branchName := fmt.Sprintf("patchpal-fix-%d", time.Now().Unix())
	repoPath := config.Config.RepoPath

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
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: config.Config.GithubToken})
	client := github.NewClient(oauth2.NewClient(ctx, ts))

	title := "PatchPal: Fix Kubernetes Misconfiguration"
	body := "This PR contains fixes generated automatically by PatchPal for a vulnerable manifest."
	pr := &github.NewPullRequest{
		Title: &title,
		Head:  &branchName,
		Base:  github.String("dev"), // adjust if your base is "main" or "master"
		Body:  &body,
	}

	createdPR, _, err := client.PullRequests.Create(ctx, config.Config.GithubOwner, config.Config.GithubRepo, pr)
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

func CreatePRWithFixedYAML1(fixedFilePath, originalFilename string) (string, error) {
	ctx := context.Background()

	// Step 1: Open repo
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", fmt.Errorf("failed to open repo: %w", err)
	}

	// Step 2: Create new branch
	worktree, err := r.Worktree()
	if err != nil {
		return "", fmt.Errorf("failed to get worktree: %w", err)
	}

	branchName := fmt.Sprintf("%s-%d", prBranchPrefix, time.Now().Unix())
	refName := plumbing.NewBranchReferenceName(branchName)

	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: refName,
		Create: true,
		Force:  false,
	})
	if err != nil {
		return "", fmt.Errorf("failed to checkout new branch: %w", err)
	}

	// Step 3: Replace the file in repo
	destPath := filepath.Join(repoPath, originalFilename)
	srcData, err := os.ReadFile(fixedFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read fixed file: %w", err)
	}

	err = os.WriteFile(destPath, srcData, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to overwrite repo file: %w", err)
	}

	// Step 4: Add and commit
	_, err = worktree.Add(originalFilename)
	if err != nil {
		return "", fmt.Errorf("failed to add file: %w", err)
	}

	_, err = worktree.Commit("fix: patch Kubernetes misconfig", &git.CommitOptions{
		Author: commitAuthor,
	})
	if err != nil {
		return "", fmt.Errorf("failed to commit: %w", err)
	}

	// Step 5: Push
	err = r.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: githubOwner, // Can be anything except blank
			Password: githubToken,
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to push: %w", err)
	}

	// Step 6: Create PR via GitHub API
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
	client := github.NewClient(oauth2.NewClient(ctx, ts))

	title := "PatchPal: Fix Kubernetes Misconfiguration"
	body := "This PR contains fixes generated automatically by PatchPal for a vulnerable manifest."
	pr := &github.NewPullRequest{
		Title: &title,
		Head:  &branchName,
		Base:  github.String("dev"), // Or "master" depending on your repo
		Body:  &body,
	}

	createdPR, _, err := client.PullRequests.Create(ctx, githubOwner, githubRepo, pr)
	if err != nil {
		return "", fmt.Errorf("failed to create GitHub PR: %w", err)
	}

	return createdPR.GetHTMLURL(), nil
}
