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
	CompanyId      int       `json:"company_id"`
}

type VacancyFilter struct {
	Title          string
	Skills         []string
	Experience     string
	EmploymentType []string
	WorkFormat     []string
	Location       string
	CompanyId      int
	Query          string
}
