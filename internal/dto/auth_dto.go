package dto

type RegisterUserRequest struct {
	Email          string `json:"email" binding:"required,email" example:"user@example.com"`
	Password       string `json:"password" binding:"required" example:"securepassword"`
	FirstName      string `json:"first_name" binding:"required" example:"John"`
	LastName       string `json:"last_name" binding:"required" example:"Doe"`
	ProfilePicture string `json:"profile_picture" binding:"required" example:"/static/images/profile.jpg"`
}

type RegisterHRRequest struct {
	Email     string `json:"email" binding:"required,email" example:"user@example.com"`
	Password  string `json:"password" binding:"required,min=6" example:"securepassword"`
	FirstName string `json:"first_name" binding:"required" example:"John"`
	LastName  string `json:"last_name" binding:"required" example:"Doe"`
	DepID     int    `json:"dep_id" binding:"required" example:"1"`
	CompanyID int    `json:"company_id" binding:"required" example:"1"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"securepassword"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type RequestPasswordResetRequest struct {
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
}

type RequestPasswordResetResponse struct {
	Message string `json:"message" example:"Reset code sent to your email"`
}

type ResetPasswordRequest struct {
	Email           string `json:"email" binding:"required,email" example:"user@example.com"`
	ResetCode       string `json:"reset_code" binding:"required" example:"123456"`
	NewPassword     string `json:"new_password" binding:"required,min=6" example:"newSecurePass"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=6" example:"newSecurePass"`
}

type ResetPasswordResponse struct {
	Message string `json:"message" example:"Password reset successful"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}
