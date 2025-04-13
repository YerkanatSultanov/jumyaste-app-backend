package dto

import "time"

type JobApplicationResponse struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	VacancyID int       `json:"vacancy_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Status    string    `json:"status"`
	AppliedAt time.Time `json:"applied_at"`
}

type JobApplicationWithResumeResponse struct {
	ID              int            `json:"id"`
	UserID          int            `json:"user_id"`
	VacancyID       int            `json:"vacancy_id"`
	FirstName       string         `json:"first_name"`
	LastName        string         `json:"last_name"`
	Email           string         `json:"email"`
	Status          string         `json:"status"`
	AppliedAt       string         `json:"applied_at"`
	Resume          ResumeResponse `json:"resume"`
	AIMatchingScore int            `json:"ai_matching_score"`
}
