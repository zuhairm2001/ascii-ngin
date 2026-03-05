package ai

import (
	"context"
	"os"
	"testing"

	"google.golang.org/genai"
)

func TestGenerateImage(t *testing.T) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		t.Skip("GEMINI_API_KEY not set, skipping integration test")
	}

	ctx := context.Background()

	client, err := InitGeminiClient(ctx, apiKey)
	if err != nil {
		t.Fatalf("failed to init Gemini client: %v", err)
	}

	resp, err := client.Models.List(ctx, &genai.ListModelsConfig{})
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Log(resp.SDKHTTPResponse.Body)

	prompt := "A white cat sitting on a table"
	uri, err := generateImage(ctx, client, prompt)
	if err != nil {
		t.Fatalf("generateImage failed: %v", err)
	}

	if uri == "" {
		t.Fatal("expected non-empty URI, got empty string")
	}

	t.Logf("generated image URI: %s", uri)
}
