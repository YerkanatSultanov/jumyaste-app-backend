package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"jumyste-app-backend/internal/dto"
	"jumyste-app-backend/pkg/helper"
	"log"
	"net/http"
	"os"
	"strconv"
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

func (c *OpenAIClient) GenerateVacancyDescription(input dto.VacancyInput) (string, error) {
	prompt := fmt.Sprintf(`
Ты HR-специалист. На основе следующих данных о вакансии сгенерируй качественное, подробное и привлекательное описание вакансии на русском языке в формате HTML.

Данные:
Название: %s  
Тип занятости: %s  
Формат работы: %s  
Навыки: %s  
Локация: %s  
Опыт: %s  
Зарплата: от %d до %d  

Описание должно быть:
- структурированным и легко читаемым,
- оформленным в HTML с использованием тегов <b>, <i>, <ul>, <li>, <p>, <h2> и других при необходимости,
- с выделением ключевых разделов: «Обязанности», «Требования», «Мы предлагаем».

Результат должен быть валидным HTML-блоком:
- Каждую информацию представь в соответствующем блоке: заголовок в <h2>, текстовые блоки в <p>,
- Используй <ul> и <li> для списка обязанностей, требований и предложений,
- Обязательно разделяй разделы с помощью <h2> для заголовков,
- Не используйте переносы строк с \n или <br>, оформляй структуру с помощью HTML.

Важно:
- не выдумывай информацию — опирайся строго на предоставленные данные;
- избегай излишней общности — пиши конкретно и по делу;
- сделай описание привлекательным для кандидатов.

Результат должен быть валидным HTML-блоком, без объяснений и комментариев.
`, input.Title, input.EmploymentType, input.WorkFormat, strings.Join(input.Skills, ", "), input.Location, input.Experience, input.SalaryMin, input.SalaryMax)

	requestBody, err := json.Marshal(OpenAIRequest{
		Model: "gpt-4o-mini",
		Messages: []Message{
			{Role: "system", Content: "Ты помощник HR, создающий описания вакансий."},
			{Role: "user", Content: prompt},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to encode request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var openAIResp OpenAIResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return "", fmt.Errorf("failed to parse OpenAI response: %w", err)
	}

	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	description := strings.TrimSpace(openAIResp.Choices[0].Message.Content)
	return description, nil
}

type MatchingResult struct {
	Score      int
	Strengths  string
	Weaknesses string
}

func (c *OpenAIClient) GetMatchingScore(resumeText string, vacancyDescription string) (*MatchingResult, error) {
	prompt := fmt.Sprintf(`
Оцени соответствие резюме описанию вакансии по шкале от 0 до 100.
Затем обязательно укажи 2–3 сильные стороны и 2–3 слабые стороны кандидата.

Ответ **только строго в формате**:
Оценка: <число>
Преимущества: <список через запятую>
Недостатки: <список через запятую>

Резюме:
%s

Вакансия:
%s
`, resumeText, vacancyDescription)

	requestBody, err := json.Marshal(OpenAIRequest{
		Model: "gpt-4o-mini",
		Messages: []Message{
			{Role: "system", Content: "Ты AI-рекрутер. Оцени соответствие кандидата вакансии, дай сильные и слабые стороны. Ответ строго в нужном формате."},
			{Role: "user", Content: prompt},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to encode analysis request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create analysis request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send analysis request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read analysis response: %w", err)
	}

	var openAIResp OpenAIResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return nil, fmt.Errorf("failed to parse analysis response: %w", err)
	}

	if len(openAIResp.Choices) == 0 {
		return nil, fmt.Errorf("empty analysis response")
	}

	text := strings.TrimSpace(openAIResp.Choices[0].Message.Content)
	lines := strings.Split(text, "\n")

	result := &MatchingResult{}

	for _, line := range lines {
		line = strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(line, "Оценка:"):
			scoreStr := strings.TrimSpace(strings.TrimPrefix(line, "Оценка:"))
			score, err := strconv.Atoi(scoreStr)
			if err != nil {
				return nil, fmt.Errorf("invalid score format: %w", err)
			}
			result.Score = score

		case strings.HasPrefix(line, "Преимущества:"),
			strings.HasPrefix(line, "Плюсы:"):
			value := strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
			result.Strengths = value

		case strings.HasPrefix(line, "Недостатки:"),
			strings.HasPrefix(line, "Минусы:"):
			value := strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
			result.Weaknesses = value
		}
	}

	return result, nil
}
