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
	prompt := fmt.Sprintf(`You are an expert in Kubernetes security. 
Given the following K8s YAML manifest and the Trivy scan results, suggest a fixed version of the manifest that resolves the issues.

Trivy Report (JSON):
%s

Original YAML:
%s

Respond with only the corrected YAML.`, vulnReport, originalYAML)

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
