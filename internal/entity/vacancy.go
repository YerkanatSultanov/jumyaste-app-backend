package entity

import "time"

type Vacancy struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	EmploymentType string    `json:"employment_type"`
	WorkFormat     string    `json:"work_format"`
	Experience     string    `json:"experience"`
	SalaryMin      *int      `json:"salary_min,omitempty"`
	SalaryMax      *int      `json:"salary_max,omitempty"`
	Location       *string   `json:"location,omitempty"`
	Category       *string   `json:"category,omitempty"`
	Skills         []string  `json:"skills"`
	Description    string    `json:"description"`
	CreatedBy      int       `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
}

type VacancyFilter struct {
	Title          string   `form:"title"`
	Skills         []string `form:"skills[]"`
	Experience     string   `form:"experience"`
	EmploymentType []string `form:"employment_type[]"`
	WorkFormat     []string `form:"work_format[]"`
	Location       string   `form:"location"`
}
