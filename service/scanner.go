package service

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// TrivyResult is a generic map for holding Trivy scan results
type TrivyResult map[string]interface{}

// ScanWithTrivy runs a Trivy scan on the given filePath and outputs the results to a JSON file
// Returns the path to the output file or an error if the scan fails
func ScanWithTrivy(filePath string) (string, error) {
	// Step 1: Ensure Trivy is installed
	log.Printf("[Scanner] Checking for Trivy installation...")
	if _, err := exec.LookPath("trivy"); err != nil {
		return "", fmt.Errorf("trivy not found: %v", err)
	}

	// Step 2: Prepare local report folder next to the YAML
	yamlDir := filepath.Dir(filePath)
	reportDir := filepath.Join(yamlDir, "patchpal-json-report")
	if err := os.MkdirAll(reportDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create report directory: %w", err)
	}

	// Step 3: Prepare output JSON filename
	baseName := filepath.Base(filePath)
	safeName := strings.TrimSuffix(baseName, filepath.Ext(baseName))
	outputFile := filepath.Join(reportDir, fmt.Sprintf("trivy-out-%d-%s.json", time.Now().UnixNano(), safeName))

	// Step 4: Run Trivy
	log.Printf("[Scanner] Running Trivy scan on: %s", filePath)
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
