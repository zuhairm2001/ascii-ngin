package ai

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

func InitGeminiClient(ctx context.Context, apiKey string) (*genai.Client, error) {
	config := genai.ClientConfig{
		APIKey: apiKey,
	}

	client, err := genai.NewClient(ctx, &config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func generateImage(ctx context.Context, client *genai.Client, prompt string) (string, error) {
	resp, err := client.Models.GenerateImages(
		ctx,
		"gemini-3.1-flash-image-preview",
		prompt+PROMPT_SUFFIX,
		&genai.GenerateImagesConfig{},
	)
	if err != nil {
		return "", fmt.Errorf("Error in image generation : %s", err)
	}

	return resp.GeneratedImages[0].Image.GCSURI, nil
}
