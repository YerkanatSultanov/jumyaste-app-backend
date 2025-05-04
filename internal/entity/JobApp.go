package entity

import "time"

type JobApplication struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	VacancyID       int       `json:"vacancy_id"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	Email           string    `json:"email"`
	Status          string    `json:"status"`
	AppliedAt       time.Time `json:"applied_at"`
	ResumeID        int       `json:"resume_id"`
	AIMatchingScore int       `json:"ai_matching_score"`
	AIStrengths     string    `json:"ai_strengths"`
	AIWeaknesses    string    `json:"ai_weaknesses"`
}

type JobApplicationWithResume struct {
	JobApplication
	Resume Resume `json:"resume"`
	User   User   `json:"user"`
}

type ApplicationStatusStat struct {
	Status     string  `json:"status"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}
