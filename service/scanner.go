package service

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"
)

type TrivyResult map[string]interface{}

func ScanWithTrivy(filePath string) (string, error) {
	const tmpDir = "/home/one2n/Desktop/patchpal_out" // trivy is failing in dir /tmp
	// Ensure Trivy is installed
	if _, err := exec.LookPath("trivy"); err != nil {
		return "", fmt.Errorf("trivy not found: %v", err)
	}

	// Prepare output file path
	outputFile := filepath.Join(tmpDir, fmt.Sprintf("trivy-out-%d.json", time.Now().UnixNano()))

	// Run the Trivy command
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

	return outputFile, nil
}
