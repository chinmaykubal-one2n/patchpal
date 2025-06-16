package api

import (
	"net/http"
	"patchpal/service"
	"patchpal/utils"

	"github.com/gin-gonic/gin"
)

func HandleFix(c *gin.Context) {
	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File not provided"})
		return
	}

	// Save file to tmp directory
	savedPath, err := utils.SaveUploadedFile(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save uploaded file"})
		return
	}

	// Run Trivy scan on the uploaded file
	trivyResult, err := service.ScanWithTrivy(savedPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Trivy scan failed",
			"details": err.Error(),
		})
		return
	}

	// Respond with scan result (can be large; limit fields later)
	c.JSON(http.StatusOK, gin.H{
		"message":      "Trivy scan completed successfully",
		"file_path":    savedPath,
		"scan_summary": trivyResult,
	})
}
