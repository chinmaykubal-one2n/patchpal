package service

import (
	"context"
	"fmt"
	"patchpal/config"
	"patchpal/parser"
	"patchpal/utils"
)

func ProcessAndFixManifests() error {
	allFiles := utils.FindAllYAMLFiles(config.Config.RepoPath)
	if len(allFiles) == 0 {
		return fmt.Errorf("no Kubernetes YAML files found in repo")
	}

	llm := NewLLMService(config.Config.OpenRouterAPIKey, config.Config.Model)

	var filesToCommit []string

	for _, file := range allFiles {
		fmt.Println("Processing:", file)

		// Step 1: Trivy Scan
		trivyJSONPath, err := ScanWithTrivy(file)
		if err != nil {
			fmt.Printf("Trivy scan failed for %s: %v\n", file, err)
			continue
		}

		// Step 2: Parse Trivy JSON
		misconfigs, err := parser.ExtractRelevantMisconfigs(trivyJSONPath)
		if err != nil || len(misconfigs) == 0 {
			fmt.Printf("No high-severity misconfigs in %s\n", file)
			continue
		}

		// Step 3: Format prompt and call LLM
		prompt := parser.FormatMisconfigsAsPrompt(misconfigs)
		ctx := context.Background()

		fixedYAML, err := llm.FixK8sManifest(ctx, prompt, file)
		if err != nil {
			fmt.Printf("LLM fix failed for %s: %v\n", file, err)
			continue
		}

		// Step 4: Overwrite original file
		if err := utils.OverwriteFile(file, fixedYAML); err != nil {
			fmt.Printf("Failed to overwrite file %s: %v\n", file, err)
			continue
		}

		filesToCommit = append(filesToCommit, file)
	}

	// Step 5: Commit all changes and raise PR
	if len(filesToCommit) > 0 {
		prURL, err := CreateBatchPR(filesToCommit)
		if err != nil {
			return fmt.Errorf("failed to create PR: %w", err)
		}
		fmt.Println("PR Created:", prURL)
	}

	return nil
}
