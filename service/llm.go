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
Understand Description, Message and based on Resolution make exact changes to given YAML file and nothing else, do not do anything extra to file, do not remove anything that does not concern with the trivy report.

---
This is a trivy report:
%s
---
For the this manifest:
%s
`, vulnReport, originalYAML)

	req := openai.ChatCompletionRequest{
		Model:       l.Model,
		Temperature: 0,
		TopP:        1,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: `Fix the manifest with respect to the Resolution suggested in the report and give me updated manifest. DO NOT explain. DO NOT use markdown. DO NOT add comments. DO NOT return anything other than plain corrected YAML.`,
			},
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
