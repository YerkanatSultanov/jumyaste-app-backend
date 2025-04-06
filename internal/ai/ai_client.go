package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"jumyste-app-backend/pkg/helper"
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

type Resume struct {
	FullName        string   `json:"full_name"`
	DesiredPosition string   `json:"desired_position"`
	Skills          []string `json:"skills"`
	City            string   `json:"city"`
	AboutMe         string   `json:"about_me"`
}

func (c *OpenAIClient) AnalyzeResume(text string) (*Resume, error) {
	prompt := fmt.Sprintf(`
Parse the following resume text and return a JSON object with the following structure:

{
  "full_name": "Full name",
  "desired_position": "Desired job position",
  "skills": ["Skill1", "Skill2", "Skill3"],
  "city": "City of residence",
  "about_me": "Everything else important from the resume (experience, achievements, etc.)"
}

If some information is missing — leave an empty string or an empty array.

Resume text:
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
	log.Printf("Raw AI Response: %s", string(body))

	var openAIResp OpenAIResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		log.Printf("Invalid OpenAI response: %s", string(body))
		return nil, fmt.Errorf("failed to parse OpenAI response: %w", err)
	}

	if len(openAIResp.Choices) == 0 {
		return nil, fmt.Errorf("no response from AI")
	}

	responseText := strings.TrimSpace(openAIResp.Choices[0].Message.Content)
	log.Printf("AI Response: %s", responseText)

	// Важная часть — достаём JSON из текста
	jsonStr := helper.ExtractJSON(responseText)
	if jsonStr == "" {
		log.Printf("Failed to extract JSON: %s", responseText)
		return nil, fmt.Errorf("failed to extract JSON from response")
	}

	var resume Resume
	if err := json.Unmarshal([]byte(jsonStr), &resume); err != nil {
		log.Printf("Failed to parse JSON from AI response: %s", jsonStr)
		return nil, fmt.Errorf("failed to parse JSON from AI response: %w", err)
	}

	return &resume, nil
}
