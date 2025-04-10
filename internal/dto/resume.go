package dto

type ResumeRequest struct {
	FullName        string   `json:"full_name" binding:"required"`
	DesiredPosition string   `json:"desired_position" binding:"required"`
	Skills          []string `json:"skills" binding:"required"`
	City            string   `json:"city"`
	About           string   `json:"about"`
}

type ResumeResponse struct {
	FullName        string       `json:"full_name"`        // Полное имя
	DesiredPosition string       `json:"desired_position"` // Желаемая должность
	Skills          []string     `json:"skills"`           // Навыки
	City            string       `json:"city"`             // Город
	About           string       `json:"about"`            // Информация о себе
	ParsedData      interface{}  `json:"parsed_data"`      // Дополнительные данные, извлеченные из резюме
	User            UserResponse `json:"user"`             // Информация о пользователе
}

type UserResponse struct {
	ID             int    `json:"id"`              // ID пользователя
	Email          string `json:"email"`           // Электронная почта
	FirstName      string `json:"first_name"`      // Имя пользователя
	LastName       string `json:"last_name"`       // Фамилия пользователя
	ProfilePicture string `json:"profile_picture"` // Профильное изображение
	RoleID         int    `json:"role_id"`         // ID роли пользователя
}
