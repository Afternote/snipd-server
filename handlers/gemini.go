package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofor-little/env"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func GeminiHandler(w http.ResponseWriter, r *http.Request) {

    // Assuming getResponseFromGemini returns the generated text as a string
    generatedText, err := getResponseFromGemini(w)

    if err != nil {
        // Handle error appropriately (e.g., log the error and return a generic error message)
        log.Println("Error generating content:", err)
        fmt.Fprintf(w, "Error generating content")
        return
    }

    response := struct {
        Message string `json:"message"`
        Content string `json:"content"`
    }{
        Message: "Response from Gemini",
        Content: generatedText,
    }

    err = json.NewEncoder(w).Encode(response)
    if err != nil {
        // Handle encoding error (e.g., log the error and return a bad request status)
        log.Println("Error encoding response:", err)
        http.Error(w, "Error encoding response", http.StatusBadRequest)
        return
    }
}

func getResponseFromGemini(w http.ResponseWriter) (string, error) {
    ctx := context.Background()

    // Replace with your actual API key
    client, err := genai.NewClient(ctx, option.WithAPIKey(goDotEnvVariable("GEMINI_API_KEY")))
    if err != nil {
        return "", err
    }
    defer client.Close()

    model := client.GenerativeModel("gemini-1.5-flash")
    resp, err := model.GenerateContent(ctx, genai.Text("Write a story about a magic backpack."))
    if err != nil {
        return "", err
    }
	encodeJSON(w, resp)
	
    return "hellooo from gljgldg", nil
}

func goDotEnvVariable(key string) string {

	// load .env file
	err := env.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

