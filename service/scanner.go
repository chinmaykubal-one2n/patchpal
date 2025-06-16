package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// TrivyResult is a generic structure to unmarshal trivy output.
type TrivyResult map[string]interface{}

// ScanWithTrivy runs a Trivy config scan on the given file and returns the parsed JSON output
func ScanWithTrivy(filePath string) (TrivyResult, error) {
	// Ensure Trivy is installed
	if _, err := exec.LookPath("trivy"); err != nil {
		return nil, fmt.Errorf("trivy not found: %v", err)
	}

	// Prepare output file path
	outputFile := filepath.Join(os.TempDir(), fmt.Sprintf("trivy-out-%d.json", time.Now().UnixNano()))

	// Run the Trivy command
	cmd := exec.Command(
		"trivy", "config", filePath,
		"--format", "json",
		"--output", outputFile,
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("trivy scan failed: %v\nstderr: %s", err, stderr.String())
	}

	// Read and parse JSON output
	data, err := os.ReadFile(outputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read trivy output: %w", err)
	}

	var result TrivyResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse trivy json: %w", err)
	}

	return result, nil
}
