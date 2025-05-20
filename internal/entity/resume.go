package entity

import "time"

type Resume struct {
	ID              int                    `json:"id"`
	UserID          int                    `json:"user_id"`
	FullName        string                 `json:"full_name"`
	DesiredPosition string                 `json:"desired_position"`
	Skills          []string               `json:"skills"`
	City            string                 `json:"city"`
	About           string                 `json:"about"`
	ParsedData      map[string]interface{} `json:"parsed_data"`
	Experiences     []WorkExperience       `json:"experiences"`
	CreatedAt       time.Time              `json:"created_at"`
}

type WorkExperience struct {
	CompanyName    string `json:"company_name"`
	Position       string `json:"position"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
	Location       string `json:"location"`
	EmploymentType string `json:"employment_type"`
	Description    string `json:"description"`
}
