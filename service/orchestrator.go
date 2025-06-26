package service

import (
	"context"
	"fmt"
	"log"
	"patchpal/config"
	"patchpal/parser"
	"patchpal/utils"
)

// ProcessAndFixManifests scans all YAML files in the repo, fixes misconfigurations using LLM, and creates a PR
// It orchestrates the full PatchPal workflow
func ProcessAndFixManifests(ctx context.Context) error {
	allFiles := utils.FindAllYAMLFiles(config.Config.RepoPath)
	if len(allFiles) == 0 {
		return fmt.Errorf("no Kubernetes YAML files found in repo")
	}

	llm := NewLLMService(config.Config.OpenRouterAPIKey, config.Config.Model)

	var filesToCommit []string

	for _, file := range allFiles {
		log.Println("Processing:", file)

		// Step 1: Trivy Scan for misconfigurations
		trivyJSONPath, err := ScanWithTrivy(file)
		if err != nil {
			log.Printf("Trivy scan failed for %s: %v", file, err)
			continue
		}

		// Step 2: Parse Trivy JSON for relevant misconfigurations
		misconfigs, err := parser.ExtractRelevantMisconfigs(trivyJSONPath)
		if err != nil || len(misconfigs) == 0 {
			log.Printf("No high-severity misconfigs in %s", file)
			continue
		}

		// Step 3: Format prompt and call LLM to fix manifest
		prompt := parser.FormatMisconfigsAsPrompt(misconfigs)

		fixedYAML, err := llm.FixK8sManifest(ctx, prompt, file)
		if err != nil {
			log.Printf("LLM fix failed for %s: %v", file, err)
			continue
		}

		// Step 4: Overwrite original file with fixed YAML
		if err := utils.OverwriteFile(file, fixedYAML); err != nil {
			log.Printf("Failed to overwrite file %s: %v", file, err)
			continue
		}

		filesToCommit = append(filesToCommit, file)
	}

	// Step 5: Commit all changes and raise PR if any files were fixed
	if len(filesToCommit) > 0 {
		prURL, err := CreateBatchPR(ctx, filesToCommit)
		if err != nil {
			return fmt.Errorf("failed to create PR: %w", err)
		}
		log.Println("PR Created:", prURL)
	}

	return fmt.Errorf("failed to create PR")
}
