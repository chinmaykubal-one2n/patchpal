package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// FindAllYAMLFiles recursively scans a directory and returns .yml or .yaml files
func FindAllYAMLFiles(root string) []string {
	var files []string

	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // ignore files that cause errors
		}
		// Skip the .github folder and its contents
		if info.IsDir() && strings.Contains(path, ".github") {
			return filepath.SkipDir
		}
		if !info.IsDir() && (strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml")) {
			files = append(files, path)
		}
		return nil
	})

	return files
}

func OverwriteFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}
