package ai

import (
	"context"

	"google.golang.org/genai"
)

func InitializeGeminiClient(ctx context.Context, apiKey string) (*genai.Client, error) {
	config := genai.ClientConfig{
		APIKey: apiKey,
	}

	client, err := genai.NewClient(ctx, &config)

	if err != nil {
		return nil, err
	}
	return client, nil
}
