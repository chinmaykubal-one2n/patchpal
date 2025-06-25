package service

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// TrivyResult is a generic map for holding Trivy scan results
type TrivyResult map[string]interface{}

// ScanWithTrivy runs a Trivy scan on the given filePath and outputs the results to a JSON file
// Returns the path to the output file or an error if the scan fails
func ScanWithTrivy(filePath string) (string, error) {
	// Create a temporary directory for storing the Trivy JSON report
	tmpDir := filepath.Join(os.TempDir(), "patchpal-trivy-json-report")
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Step 1: Ensure Trivy is installed
	log.Printf("[Scanner] Checking for Trivy installation...")
	if _, err := exec.LookPath("trivy"); err != nil {
		return "", fmt.Errorf("trivy not found: %v", err)
	}

	// Step 2: Prepare output file path
	outputFile := filepath.Join(tmpDir, fmt.Sprintf("trivy-out-%d.json", time.Now().UnixNano()))
	log.Printf("[Scanner] Running Trivy scan on: %s", filePath)

	// Step 3: Run the Trivy command to scan the file and output JSON
	cmd := exec.Command(
		"trivy", "config", filePath,
		"--format", "json",
		"--output", outputFile,
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("trivy scan failed: %v\nstderr: %s", err, stderr.String())
	}
	log.Printf("[Scanner] Trivy scan completed. Output: %s", outputFile)
	return outputFile, nil
}
