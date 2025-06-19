package api

import (
	"context"
	"net/http"
	"patchpal/config"
	"patchpal/parser"
	"patchpal/service"
	"patchpal/utils"

	"github.com/gin-gonic/gin"
)

func HandleFix(c *gin.Context) {
	// Step 1: Accept uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File not provided"})
		return
	}

	originalFilename := file.Filename

	// Step 2: Save file
	savedPath, err := utils.SaveUploadedFile(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save uploaded file"})
		return
	}

	// Step 3: Run Trivy Scan
	trivyOutputPath, err := service.ScanWithTrivy(savedPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Trivy scan failed", "details": err.Error()})
		return
	}

	// Step 4: Parse relevant misconfigurations
	misconfigs, err := parser.ExtractRelevantMisconfigs(trivyOutputPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse misconfigs", "details": err.Error()})
		return
	}

	// Step 5: Initialize LLM service
	llm := service.NewLLMService(config.Config.OpenRouterAPIKey, config.Config.Model)

	prompt := parser.FormatMisconfigsAsPrompt(misconfigs)
	ctx := context.Background()

	// Step 6: Call FixK8sManifest to get the fixed YAML
	fixedYaml, err := llm.FixK8sManifest(ctx, prompt, savedPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fix manifest", "details": err.Error()})
		return
	}

	// Step 7: Save the fixed YAML to a file
	fixedPath, err := utils.SaveFixedManifest(fixedYaml, originalFilename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save fixed manifest", "details": err.Error()})
		return
	}

	// Step 8: Commit and create PR
	prURL, err := service.CreatePRWithFixedYAML(fixedPath, originalFilename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create PR", "details": err.Error()})
		return
	}

	// Step 9: Return PR URL and fixed YAML
	c.JSON(http.StatusOK, gin.H{
		"pr_url": prURL,
		// "fixed_yaml": fixedYaml,
	})
}
