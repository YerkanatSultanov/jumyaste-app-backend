package service

import (
	"bytes"
	"fmt"
	"io"
	"jumyste-app-backend/internal/ai"
	"jumyste-app-backend/pkg/logger"
	"regexp"
	"strings"

	"github.com/ledongthuc/pdf"
)

type ResumeService struct {
	AIClient *ai.OpenAIClient
}

func NewResumeService(aiClient *ai.OpenAIClient) *ResumeService {
	return &ResumeService{AIClient: aiClient}
}

func (s *ResumeService) ProcessResume(file io.Reader) (map[string]interface{}, error) {
	logger.Log.Info("Reading PDF file...")

	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, file)
	if err != nil {
		logger.Log.Error("Failed to read PDF file", "error", err)
		return nil, err
	}

	reader := bytes.NewReader(buf.Bytes())
	pdfReader, err := pdf.NewReader(reader, int64(buf.Len()))
	if err != nil {
		logger.Log.Error("Failed to parse PDF", "error", err)
		return nil, err
	}

	var text string
	numPages := pdfReader.NumPage()
	for i := 1; i <= numPages; i++ {
		page := pdfReader.Page(i)
		pageText, err := page.GetPlainText(nil)
		if err != nil {
			logger.Log.Warn("Failed to extract text from page", "page", i, "error", err)
			continue
		}
		text += pageText + "\n"
	}

	if text == "" {
		logger.Log.Error("No text extracted from resume")
		return nil, fmt.Errorf("empty text extracted")
	}

	cleanText := preprocessText(text)

	logger.Log.Info("Extracted and cleaned text from resume", "text", cleanText)

	parsedResume, err := s.AIClient.AnalyzeResume(cleanText)
	if err != nil {
		logger.Log.Error("Failed to analyze resume with AI", "error", err)
		return nil, err
	}

	return parsedResume, nil
}

func preprocessText(text string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9А-Яа-я\s.,:;!?()\[\]@-]`)
	text = re.ReplaceAllString(text, "")

	text = strings.Join(strings.Fields(text), " ")

	text = strings.ReplaceAll(text, ",", ", ")

	text = strings.TrimSpace(text)

	return text
}
