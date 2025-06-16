package api

import (
	"context"
	"net/http"
	"os"
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
	llm := service.NewLLMService(os.Getenv("OPENROUTER_API_KEY"), os.Getenv("MODEL"))

	prompt := parser.FormatMisconfigsAsPrompt(misconfigs)
	ctx := context.Background()

	// Step 6: Call FixK8sManifest to get the fixed YAML
	fixedYaml, err := llm.FixK8sManifest(ctx, savedPath, prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fix manifest", "details": err.Error()})
		return
	}

	// Step 7: Return fixed YAML
	c.Data(http.StatusOK, "text/yaml", []byte(fixedYaml))
}
