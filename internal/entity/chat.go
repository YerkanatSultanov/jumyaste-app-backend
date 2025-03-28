package entity

import "time"

type Chat struct {
	ID        int            `gorm:"primaryKey" json:"id"`
	Users     []UserResponse `gorm:"many2many:chat_users;" json:"users"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}
