package utils

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// FindAllYAMLFiles recursively scans a directory and returns all .yml or .yaml files
// It skips the .github folder and any files that cause errors
func FindAllYAMLFiles(root string) []string {
	var files []string
	log.Printf("[Utils] Scanning for YAML files in: %s\n", root)
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("[Utils] Error accessing path %s: %v", path, err)
			return nil // ignore files that cause errors
		}
		// Skip the .github folder and its contents
		if info.IsDir() && strings.Contains(path, ".github") {
			return filepath.SkipDir
		}
		if !info.IsDir() && (strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml")) {
			log.Printf("[Utils] Found YAML file: %s", path)
			files = append(files, path)
		}
		return nil
	})
	log.Printf("[Utils] Total YAML files found: %d", len(files))
	return files
}

// OverwriteFile writes the given content to the specified file path, overwriting it if it exists
func OverwriteFile(path string, content string) error {
	log.Printf("[Utils] Overwriting file: %s\n", path)
	return os.WriteFile(path, []byte(content), 0644)
}
