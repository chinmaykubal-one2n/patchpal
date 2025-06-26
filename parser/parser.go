package parser

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

// TrivyReport represents the structure of a Trivy JSON output for misconfigurations
type TrivyReport struct {
	Results []struct {
		Target            string `json:"Target"`
		Type              string `json:"Type"`
		Misconfigurations []struct {
			ID            string `json:"ID"`
			Title         string `json:"Title"`
			Description   string `json:"Description"`
			Message       string `json:"Message"`
			Resolution    string `json:"Resolution"`
			Severity      string `json:"Severity"`
			CauseMetadata struct {
				Code struct {
					Lines []struct {
						Number  int    `json:"Number"`
						Content string `json:"Content"`
						IsCause bool   `json:"IsCause"`
					} `json:"Lines"`
				} `json:"Code"`
			} `json:"CauseMetadata"`
		} `json:"Misconfigurations"`
	} `json:"Results"`
}

// SimplifiedMisconfig holds only the relevant fields for LLM prompt generation
type SimplifiedMisconfig struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Message     string   `json:"message"`
	Resolution  string   `json:"resolution"`
	Severity    string   `json:"severity"`
	CodeLines   []string `json:"code_lines"`
}

// FormatMisconfigsAsPrompt formats misconfigurations for LLM prompt input
func FormatMisconfigsAsPrompt(misconfigs []SimplifiedMisconfig) string {
	log.Printf("[Parser] Formatting %d misconfigurations for LLM prompt...\n", len(misconfigs))
	var b strings.Builder
	for _, m := range misconfigs {
		fmt.Fprintf(&b, "Description: %s\nMessage: %s\nSeverity: %s\nResolution: %s\n", m.Description, m.Message, m.Severity, m.Resolution)
	}
	// log.Println("[Parser] Formatted prompt for LLM with misconfigurations.", b.String())
	return b.String()
}

// ExtractRelevantMisconfigs parses a Trivy JSON file and extracts relevant misconfigurations for fixing
func ExtractRelevantMisconfigs(jsonPath string) ([]SimplifiedMisconfig, error) {
	log.Printf("[Parser] Reading Trivy JSON from: %s\n", jsonPath)
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON file: %w", err)
	}

	var report TrivyReport
	if err := json.Unmarshal(data, &report); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Trivy JSON: %w", err)
	}
	log.Printf("[Parser] Parsed %d result(s) from Trivy report.\n", len(report.Results))

	var misconfigs []SimplifiedMisconfig
	for _, result := range report.Results {
		if result.Type != "kubernetes" {
			continue
		}
		for _, m := range result.Misconfigurations {
			// Skip LOW and UNKNOWN severity
			if m.Severity == "LOW" || m.Severity == "UNKNOWN" {
				continue
			}

			var lines []string
			for _, l := range m.CauseMetadata.Code.Lines {
				if l.IsCause && strings.TrimSpace(l.Content) != "" {
					lines = append(lines, l.Content)
				}
			}

			misconfigs = append(misconfigs, SimplifiedMisconfig{
				ID:          m.ID,
				Title:       m.Title,
				Description: m.Description,
				Message:     m.Message,
				Resolution:  m.Resolution,
				Severity:    m.Severity,
				CodeLines:   lines,
			})
		}
	}

	return misconfigs, nil
}
