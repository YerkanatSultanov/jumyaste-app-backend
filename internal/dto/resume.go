package dto

type ResumeRequest struct {
	FullName        string   `json:"full_name" binding:"required"`
	DesiredPosition string   `json:"desired_position" binding:"required"`
	Skills          []string `json:"skills" binding:"required"`
	City            string   `json:"city"`
	About           string   `json:"about"`
}
