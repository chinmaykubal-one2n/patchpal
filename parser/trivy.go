package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type TrivyReport struct {
	Results []struct {
		Target            string `json:"Target"`
		Type              string `json:"Type"`
		Misconfigurations []struct {
			ID            string `json:"ID"`
			Title         string `json:"Title"`
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

type SimplifiedMisconfig struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	Message    string   `json:"message"`
	Resolution string   `json:"resolution"`
	Severity   string   `json:"severity"`
	CodeLines  []string `json:"code_lines"`
}

func FormatMisconfigsAsPrompt(misconfigs []SimplifiedMisconfig) string {
	var b strings.Builder
	for i, m := range misconfigs {
		fmt.Fprintf(&b, "Issue #%d:\n", i+1)
		fmt.Fprintf(&b, "ID: %s\nTitle: %s\nDescription: %s\nSeverity: %s\nResolution: %s\n\n",
			m.ID, m.Title, m.Message, m.Severity, m.Resolution)
	}
	return b.String()
}

func ExtractRelevantMisconfigs(jsonPath string) ([]SimplifiedMisconfig, error) {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON file: %w", err)
	}

	var report TrivyReport
	if err := json.Unmarshal(data, &report); err != nil {
		return nil, fmt.Errorf("failed to parse Trivy JSON: %w", err)
	}

	var misconfigs []SimplifiedMisconfig
	for _, result := range report.Results {
		if result.Type != "kubernetes" {
			continue
		}
		for _, m := range result.Misconfigurations {
			var lines []string
			for _, l := range m.CauseMetadata.Code.Lines {
				if l.IsCause && strings.TrimSpace(l.Content) != "" {
					lines = append(lines, l.Content)
				}
			}

			misconfigs = append(misconfigs, SimplifiedMisconfig{
				ID:         m.ID,
				Title:      m.Title,
				Message:    m.Message,
				Resolution: m.Resolution,
				Severity:   m.Severity,
				CodeLines:  lines,
			})
		}
	}

	return misconfigs, nil
}
