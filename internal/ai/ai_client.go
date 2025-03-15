package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type OpenAIClient struct {
	APIKey string
}

func NewOpenAIClient() *OpenAIClient {
	return &OpenAIClient{
		APIKey: os.Getenv("OPENAI_API_KEY"),
	}
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// AnalyzeResume отправляет текст резюме в AI и получает JSON-ответ
func (c *OpenAIClient) AnalyzeResume(text string) (map[string]interface{}, error) {
	prompt := fmt.Sprintf(`
					You are an AI specializing in analyzing resumes and extracting key fields.
					Your task is to restore spaces, correct errors, and return the result in **pure JSON format**.

					Output **only JSON**, without any explanations or additional text.

					Example response:

					{
  						"name": "Yerkanat Sultanov",
  						"contacts": { "phone": "+77471089155", "email": "example@mail.com" },
  						"skills": ["Go", "Python", "PostgreSQL"]
					}

					⚠️ Important:
					- Output **only JSON**, with **no explanations or extra text**.
					- **Do not include empty fields**.
					- Fix missing spaces and formatting issues in the text.

					Here is the resume text:
					%s
					`, text)

	requestBody, err := json.Marshal(OpenAIRequest{
		Model: "gpt-4o-mini",
		Messages: []Message{
			{Role: "system", Content: "You are an AI that extracts structured resume information in JSON format."},
			{Role: "user", Content: prompt},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var openAIResp OpenAIResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		log.Printf("Invalid OpenAI response: %s", string(body))
		return nil, fmt.Errorf("failed to parse OpenAI response: %w", err)
	}

	if len(openAIResp.Choices) == 0 {
		return nil, fmt.Errorf("no response from AI")
	}

	responseText := strings.TrimSpace(openAIResp.Choices[0].Message.Content)

	if !strings.HasPrefix(responseText, "{") {
		log.Printf("Unexpected AI response format: %s", responseText)
		return nil, fmt.Errorf("unexpected response format")
	}

	var parsedResume map[string]interface{}
	if err := json.Unmarshal([]byte(responseText), &parsedResume); err != nil {
		log.Printf("Failed to parse JSON from AI response: %s", responseText)
		return nil, fmt.Errorf("failed to parse JSON from AI response: %w", err)
	}

	return parsedResume, nil
}
