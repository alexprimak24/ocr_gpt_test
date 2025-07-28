package main

import (
	"context"
	"fmt"
	"log"
	"os"

	vision "cloud.google.com/go/vision/apiv1"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"

	openai "github.com/sashabaranov/go-openai"
)

func main() {
	// Load .env into environment
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	credPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credPath == "" {
		log.Fatal("GOOGLE_APPLICATION_CREDENTIALS not set")
	}

	gptKey := os.Getenv("OPENAI_API_KEY")
	if credPath == "" {
		log.Fatal("OPENAI_API_KEY not set")
	}

	gptClient := openai.NewClient(gptKey)

	var frontText string
	var backText string

	ctx := context.Background()

	// Load credentials from JSON file path stored in env var
	// Example: export GOOGLE_APPLICATION_CREDENTIALS="service-account.json"
	client, err := vision.NewImageAnnotatorClient(ctx, option.WithCredentialsFile(credPath))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// // Replace with path to your test image
	// fileName := "back_photo.jpg"

	// image := vision.NewImageFromFilename(fileName)
	file, err := os.Open("1.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	image, err := vision.NewImageFromReader(file)
	if err != nil {
		log.Fatal(err)
	}

	// Run OCR
	annotations, err := client.DetectTexts(ctx, image, nil, 1)
	if err != nil {
		log.Fatalf("Failed to detect text: %v", err)
	}

	if len(annotations) > 0 {
		fmt.Println("Detected text:", annotations[0].Description)
		frontText = annotations[0].Description
	} else {
		fmt.Println("No text found")
	}

	file1, err := os.Open("2.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	image1, err := vision.NewImageFromReader(file1)
	if err != nil {
		log.Fatal(err)
	}

	// Run OCR
	annotations1, err := client.DetectTexts(ctx, image1, nil, 1)
	if err != nil {
		log.Fatalf("Failed to detect text: %v", err)
	}

	if len(annotations1) > 0 {
		fmt.Println("Detected text:", annotations1[0].Description)
		backText = annotations1[0].Description
	} else {
		fmt.Println("No text found")
	}

	fmt.Println("%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%")

	prompt := fmt.Sprintf(`
Here is the front label text: "%s"
Here is the back label text (ingredients): "%s"

Extract:
1. The product name
2. The list of ingredients as a clean JSON array
3. Return JSON like { "productName": ..., "ingredients": [...] }
`, frontText, backText)

	resp, err := gptClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4o,
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: "You are an assistant that extracts product names and ingredients from OCR text."},
			{Role: "user", Content: prompt},
		},
	})
	if err != nil {
		log.Fatalf("ChatCompletion error: %v", err)
	}

	fmt.Println(resp.Choices[0].Message.Content)

}


