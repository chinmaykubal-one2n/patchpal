package utils

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

// SaveUploadedFile saves the uploaded file to /tmp and returns the file path
func SaveUploadedFile(file *multipart.FileHeader) (string, error) {
	const tmpDir = "/home/one2n/Desktop/patchpal_files" // trivy is failing in dir /tmp

	dst := filepath.Join(tmpDir, fmt.Sprintf("patchpal-%d-%s", time.Now().UnixNano(), filepath.Base(file.Filename)))
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return "", err
	}
	if err := saveFile(file, dst); err != nil {
		return "", err
	}
	return dst, nil
}

func saveFile(file *multipart.FileHeader, dst string) error {
	return os.WriteFile(dst, readUploadedFile(file), 0644)
}

func readUploadedFile(file *multipart.FileHeader) []byte {
	f, _ := file.Open()
	defer f.Close()

	data := make([]byte, file.Size)
	f.Read(data)
	return data
}

// SaveFixedManifest saves the fixed YAML content to a specific directory and returns the path
func SaveFixedManifest(yamlContent string, originalFilename string) (string, error) {
	const outputDir = "/home/one2n/Desktop/patchpal_out"

	// Ensure directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create a unique filename based on timestamp
	basename := filepath.Base(originalFilename)
	fixedName := fmt.Sprintf("fixed-%d-%s", time.Now().UnixNano(), basename)
	outputPath := filepath.Join(outputDir, fixedName)

	// Write YAML content to file
	if err := os.WriteFile(outputPath, []byte(yamlContent), 0644); err != nil {
		return "", fmt.Errorf("failed to write fixed manifest: %w", err)
	}

	return outputPath, nil
}
