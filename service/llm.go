package service

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

type LLMService struct {
	Client *openai.Client
	Model  string
}

func NewLLMService(apiKey string, model string) *LLMService {
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://openrouter.ai/api/v1"
	return &LLMService{
		Client: openai.NewClientWithConfig(config),
		Model:  model,
	}
}

func (l *LLMService) FixK8sManifest(ctx context.Context, vulnReport string, originalYAML string) (string, error) {
	fmt.Println("Fixing Kubernetes manifest using LLM...", vulnReport)
	prompt := fmt.Sprintf(`
You are a Kubernetes security expert.

Your task is to fix ONLY the misconfigurations listed below in the given Kubernetes YAML manifest. 
DO NOT modify any other field.  

Only fix misconfigurations with severity HIGH, CRITICAL, or MEDIUM — ignore all others.
---
Misconfiguration Report:

%s
---

Original YAML:

%s

Rules:
1. Fix ONLY what's required by the issues above.
2. Do NOT remove or add unrelated fields.
3. Output only the fixed YAML — no markdown, no explanations, no extra text.
`, vulnReport, originalYAML)

	req := openai.ChatCompletionRequest{
		Model: l.Model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	}

	resp, err := l.Client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("LLM call failed: %w", err)
	}

	return resp.Choices[0].Message.Content, nil
}
