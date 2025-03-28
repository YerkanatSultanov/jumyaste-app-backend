package dto

// SuccessResponse стандартный ответ на успешный запрос
type SuccessResponse struct {
	Message string `json:"message" example:"User registered successfully"`
}

// ErrorResponse стандартный ответ на ошибку
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid input"`
}
