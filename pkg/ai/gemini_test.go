package ai

import (
	"context"
	"os"
	"testing"
)

func TestGenerateImage(t *testing.T) {
	ctx := context.Background()

	client, err := InitGeminiClient(ctx)
	if err != nil {
		t.Fatalf("failed to init Gemini client: %v", err)
	}

	prompt := "A white cat sitting on a table"
	imageBytes, err := GenerateImage(ctx, client, prompt)
	if err != nil {
		t.Fatalf("GenerateImage failed: %v", err)
	}

	if len(imageBytes) == 0 {
		t.Fatal("expected non-empty image data, got empty")
	}

	outputFile := "test_generated_image.png"
	if err := os.WriteFile(outputFile, imageBytes, 0o644); err != nil {
		t.Fatalf("failed to write image to %s: %v", outputFile, err)
	}

	t.Logf("wrote generated image to %s (%d bytes)", outputFile, len(imageBytes))
}
