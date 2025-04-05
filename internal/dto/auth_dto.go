package dto

type RegisterUserRequest struct {
	Email     string `json:"email" binding:"required" example:"user@example.com"`
	Password  string `json:"password" binding:"required" example:"securepassword"`
	FirstName string `json:"first_name" binding:"required" example:"John"`
	LastName  string `json:"last_name" binding:"required" example:"Doe"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"securepassword"`
}

type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type RequestPasswordResetRequest struct {
	Email string `json:"email" binding:"required" example:"user@example.com"`
}

type RequestPasswordResetResponse struct {
	Message string `json:"message" example:"Reset code sent to your email"`
}

type ResetPasswordRequest struct {
	Email           string `json:"email" binding:"required" example:"user@example.com"`
	ResetCode       string `json:"reset_code" binding:"required" example:"123456"`
	NewPassword     string `json:"new_password" binding:"required" example:"newSecurePass"`
	ConfirmPassword string `json:"confirm_password" binding:"required" example:"newSecurePass"`
}

type ResetPasswordResponse struct {
	Message string `json:"message" example:"Password reset successful"`
}
type RegisterHRRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	DepID     int    `json:"dep_id" binding:"required"`
	CompanyID int    `json:"company_id" binding:"required"`
}
