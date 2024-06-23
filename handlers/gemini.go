package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofor-little/env"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type Snippet struct {
    Title    string    `json:"title"`
    Category string    `json:"category"`
    Content  string    `json:"content"`
    DateTime time.Time `json:"date"`
    Notes    string    `json:"customNotes"`
    Type     string    `json:"type"`
    Source   string    `json:"source"`
    ID       string    `json:"id"`
}

func GeminiHandler(w http.ResponseWriter, r *http.Request) {

    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var bodyBytes bytes.Buffer
    if _, err := io.Copy(&bodyBytes, r.Body); err != nil {
        http.Error(w, "Error reading request body", http.StatusBadRequest)
        return
    }
    r.Body.Close()

    requestBody := bodyBytes.String()

    var snippets []Snippet
	err := json.Unmarshal([]byte(requestBody), &snippets)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

    var mergedContent string
	for _, snippet := range snippets {
		mergedContent += snippet.Content
	}

    // Assuming getResponseFromGemini returns the generated text as a string
    generatedText, err := getResponseFromGemini(w, mergedContent)

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

func getResponseFromGemini(w http.ResponseWriter, mergedContent string) (string, error) {
    ctx := context.Background()

    // Replace with your actual API key
    client, err := genai.NewClient(ctx, option.WithAPIKey(goDotEnvVariable("GEMINI_API_KEY")))
    if err != nil {
        return "", err
    }
    defer client.Close()

    model := client.GenerativeModel("gemini-1.5-flash")
    resp, err := model.GenerateContent(ctx, genai.Text("create a summary from the following content: " + mergedContent))
    if err != nil {
        return "", err
    }
	encodeJSON(w, resp)
	
    return "", nil
}

func goDotEnvVariable(key string) string {

	// load .env file
	err := env.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

