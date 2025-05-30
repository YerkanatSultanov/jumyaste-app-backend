package dto

type CreateVacancyRequest struct {
	Title          string   `json:"title" example:"Frontend developer"`
	EmploymentType string   `json:"employment_type" example:"Полная занятость"`
	WorkFormat     string   `json:"work_format" example:"Гибрид"`
	Experience     string   `json:"experience" example:"Без опыта"`
	SalaryMin      int      `json:"salary_min" example:"253000"`
	SalaryMax      int      `json:"salary_max" example:"909000"`
	Location       string   `json:"location" example:"Almaty"`
	Category       string   `json:"category" example:"IT"`
	Skills         []string `json:"skills" example:"[\"Python\",\"Node.js\"]"`
	Description    string   `json:"description" example:"<h3><strong><em><s>Hello</s></em></strong></h3>"`
}

type UpdateVacancyRequest struct {
	Title          string   `json:"title" binding:"required"`
	EmploymentType string   `json:"employment_type" binding:"required"`
	WorkFormat     string   `json:"work_format" binding:"required"`
	Experience     string   `json:"experience" binding:"required"`
	SalaryMin      int      `json:"salary_min" binding:"required"`
	SalaryMax      int      `json:"salary_max" binding:"required"`
	Location       string   `json:"location" binding:"required"`
	Category       string   `json:"category" binding:"required"`
	Skills         []string `json:"skills" binding:"required"`
	Description    string   `json:"description" binding:"required"`
}

type UpdateVacancyStatusRequest struct {
	Status string `json:"status"`
}

type FeedDataResponse struct {
	NewVacanciesCount int `json:"new_vacancies_count"`
}

type VacancyInput struct {
	Title          string   `json:"title"`
	EmploymentType string   `json:"employment_type"`
	WorkFormat     string   `json:"work_format"`
	Skills         []string `json:"skills"`
	Location       string   `json:"location"`
	Experience     string   `json:"experience"`
	SalaryMin      int      `json:"salary_min"`
	SalaryMax      int      `json:"salary_max"`
}

type DescriptionResponse struct {
	Description string `json:"description"`
}
