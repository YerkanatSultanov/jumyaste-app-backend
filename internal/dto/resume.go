package dto

type ResumeRequest struct {
	FullName        string   `json:"full_name" binding:"required"`
	DesiredPosition string   `json:"desired_position" binding:"required"`
	Skills          []string `json:"skills" binding:"required"`
	City            string   `json:"city"`
	About           string   `json:"about"`
}

type ResumeResponse struct {
	FullName        string       `json:"full_name"`
	DesiredPosition string       `json:"desired_position"`
	Skills          []string     `json:"skills"`
	City            string       `json:"city"`
	About           string       `json:"about"`
	ParsedData      interface{}  `json:"parsed_data"`
	User            UserResponse `json:"user"`
}

type UserResponse struct {
	ID             int    `json:"id"`
	Email          string `json:"email"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	ProfilePicture string `json:"profile_picture"`
	RoleID         int    `json:"role_id"`
}
