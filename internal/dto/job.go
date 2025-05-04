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
	ID                int            `json:"id"`
	UserID            int            `json:"user_id"`
	VacancyID         int            `json:"vacancy_id"`
	FirstName         string         `json:"first_name"`
	LastName          string         `json:"last_name"`
	Email             string         `json:"email"`
	Status            string         `json:"status"`
	AppliedAt         string         `json:"applied_at"`
	Resume            ResumeResponse `json:"resume"`
	AIMatchingScore   int            `json:"ai_matching_score"`
	AIStrengths       string         `json:"ai_strengths,omitempty"`
	AIMatchWeaknesses string         `json:"ai_weaknesses,omitempty"`
}

type JobAppStatusAnalytics struct {
	Status     string `json:"status" example:"new"`
	Count      int    `json:"count" example:"1"`
	Percentage int    `json:"percentage" example:"50"`
}
