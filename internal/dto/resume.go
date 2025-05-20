package dto

type UserResponse struct {
	ID             int    `json:"id"`
	Email          string `json:"email"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	ProfilePicture string `json:"profile_picture"`
	RoleID         int    `json:"role_id"`
}

type WorkExperienceRequest struct {
	CompanyName    string `json:"company_name" binding:"required"`
	Position       string `json:"position" binding:"required"`
	StartDate      string `json:"start_date" binding:"required"`
	EndDate        string `json:"end_date"`
	Location       string `json:"location"`
	EmploymentType string `json:"employment_type"`
	Description    string `json:"description"`
}

type ResumeRequest struct {
	FullName        string                  `json:"full_name" binding:"required"`
	DesiredPosition string                  `json:"desired_position" binding:"required"`
	Skills          []string                `json:"skills" binding:"required"`
	City            string                  `json:"city"`
	About           string                  `json:"about"`
	WorkExperiences []WorkExperienceRequest `json:"work_experiences"`
}

type WorkExperienceResponse struct {
	CompanyName    string `json:"company_name"`
	Position       string `json:"position"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
	City           string `json:"city"`
	EmploymentType string `json:"employment_type"`
	Description    string `json:"description"`
}

type ResumeResponse struct {
	FullName        string                   `json:"full_name"`
	DesiredPosition string                   `json:"desired_position"`
	Skills          []string                 `json:"skills"`
	City            string                   `json:"city"`
	About           string                   `json:"about"`
	ParsedData      interface{}              `json:"parsed_data"`
	User            UserResponse             `json:"user"`
	WorkExperiences []WorkExperienceResponse `json:"work_experiences"`
}
