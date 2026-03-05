package ai

import (
	"context"
	"fmt"

	"google.golang.org/genai"

	"github.com/zuhairm2001/ascii-ngin/pkg/config"
)

func InitGeminiClient(ctx context.Context) (*genai.Client, error) {
	if err := config.Load(); err != nil {
		return nil, err
	}
	vars := config.Get()

	cfg := genai.ClientConfig{
		APIKey: vars.GeminiAPIKey,
	}

	client, err := genai.NewClient(ctx, &cfg)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func GenerateImage(ctx context.Context, client *genai.Client, prompt string) ([]byte, error) {
	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-3.1-flash-image-preview",
		genai.Text(prompt+PROMPT_SUFFIX),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("image generation: %w", err)
	}

	if result == nil || len(result.Candidates) == 0 || result.Candidates[0].Content == nil {
		return nil, fmt.Errorf("image generation: empty response")
	}

	for _, part := range result.Candidates[0].Content.Parts {
		if part.InlineData != nil {
			return part.InlineData.Data, nil
		}
	}

	return nil, fmt.Errorf("image generation: no image data in response")
}
