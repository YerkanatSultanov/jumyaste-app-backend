package entity

import "time"

type User struct {
	ID             int    `json:"id"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	ProfilePicture string `json:"profile_picture"`
	RoleId         int    `json:"role_id"`
	CompanyID      int    `json:"company_id"`
	DepID          int    `json:"department_id"`
}

type UserResponse struct {
	ID             int       `json:"id"`
	Email          string    `json:"email"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	ProfilePicture string    `json:"profile_picture,omitempty"`
	Company        *Company  `json:"company"`
	CreatedAt      time.Time `json:"created_at"`
	IsOwner        bool      `json:"is_owner"`
}
