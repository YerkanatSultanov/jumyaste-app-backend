package dto

type SuccessResponse struct {
	Message string `json:"message" example:"User registered successfully"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"Invalid input"`
}
