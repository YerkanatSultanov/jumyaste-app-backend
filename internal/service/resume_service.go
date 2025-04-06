package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"jumyste-app-backend/internal/ai"
	"jumyste-app-backend/internal/dto"
	"jumyste-app-backend/internal/entity"
	"jumyste-app-backend/internal/repository"
	"jumyste-app-backend/pkg/helper"
	"jumyste-app-backend/pkg/logger"
	"regexp"
	"strings"

	"github.com/ledongthuc/pdf"
)

type ResumeService struct {
	AIClient         *ai.OpenAIClient
	ResumeRepository *repository.ResumeRepository
}

func NewResumeService(aiClient *ai.OpenAIClient, resumeRepo *repository.ResumeRepository) *ResumeService {
	return &ResumeService{AIClient: aiClient, ResumeRepository: resumeRepo}
}

func (s *ResumeService) ProcessResume(file io.Reader) (entity.Resume, error) {
	logger.Log.Info("Reading PDF file...")

	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, file)
	if err != nil {
		logger.Log.Error("Failed to read PDF file", "error", err)
		return entity.Resume{}, err
	}

	reader := bytes.NewReader(buf.Bytes())
	pdfReader, err := pdf.NewReader(reader, int64(buf.Len()))
	if err != nil {
		logger.Log.Error("Failed to parse PDF", "error", err)
		return entity.Resume{}, err
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
		return entity.Resume{}, fmt.Errorf("empty text extracted")
	}

	cleanText := preprocessText(text)

	logger.Log.Info("Extracted and cleaned text from resume", "text", cleanText)

	parsedResume, err := s.AIClient.AnalyzeResume(cleanText)
	if err != nil {
		logger.Log.Error("Failed to analyze resume with AI", "error", err)
		return entity.Resume{}, err
	}

	resume := entity.Resume{
		FullName:        parsedResume.FullName,
		DesiredPosition: parsedResume.DesiredPosition,
		Skills:          parsedResume.Skills,
		City:            parsedResume.City,
		About:           parsedResume.AboutMe,
		ParsedData:      helper.StructToMap(parsedResume),
	}

	return resume, nil
}

func preprocessText(text string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9А-Яа-я\s.,:;!?()\[\]@-]`)
	text = re.ReplaceAllString(text, "")

	text = strings.Join(strings.Fields(text), " ")

	text = strings.ReplaceAll(text, ",", ", ")

	text = strings.TrimSpace(text)

	return text
}

func (s *ResumeService) SaveResume(ctx context.Context, resume entity.Resume) error {
	logger.Log.Info("Saving resume data...")

	if err := s.ResumeRepository.CreateResume(ctx, &resume); err != nil {
		logger.Log.Error("Failed to save resume", "error", err)
		return err
	}

	return nil
}

func (s *ResumeService) CreateResumeFromRequest(ctx context.Context, userID int, req dto.ResumeRequest) error {

	exists, err := s.CheckIfResumeExists(ctx, userID)
	if err != nil {
		logger.Log.Error("Failed to check if resume exists", "error", err)
		return err
	}

	if exists {
		return fmt.Errorf("resume already exists for user with ID %d", userID)
	}

	resume := entity.Resume{
		UserID:          userID,
		FullName:        req.FullName,
		DesiredPosition: req.DesiredPosition,
		Skills:          req.Skills,
		City:            req.City,
		About:           req.About,
	}

	if err := s.ResumeRepository.CreateResume(ctx, &resume); err != nil {
		logger.Log.Error("Failed to save resume", "error", err)
		return err
	}

	return nil
}

func (s *ResumeService) CheckIfResumeExists(ctx context.Context, userID int) (bool, error) {
	resume, _, err := s.ResumeRepository.GetResumeByUserID(ctx, userID)
	if err != nil {
		logger.Log.Error("Failed to check if resume exists", "error", err)
		return false, err
	}

	if resume != nil {
		return true, nil
	}

	return false, nil
}

func (s *ResumeService) GetResumeAndUserByUserID(ctx context.Context, userID int) (*entity.Resume, *entity.User, error) {
	resume, user, err := s.ResumeRepository.GetResumeByUserID(ctx, userID)
	if err != nil {
		logger.Log.Error("Failed to get resume and user", "error", err)
		return nil, nil, err
	}

	if resume == nil || user == nil {
		return nil, nil, fmt.Errorf("no resume or user found for user_id %d", userID)
	}

	return resume, user, nil
}

func (s *ResumeService) DeleteResumeByUserID(ctx context.Context, userID int) error {
	exists, err := s.CheckIfResumeExists(ctx, userID)
	if err != nil {
		logger.Log.Error("Failed to check if resume exists", "error", err)
		return err
	}

	if !exists {
		return fmt.Errorf("no resume found for user with ID %d", userID)
	}

	err = s.ResumeRepository.DeleteResumeByUserID(ctx, userID)
	if err != nil {
		logger.Log.Error("Failed to delete resume", "error", err)
		return err
	}

	return nil
}
