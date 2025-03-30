package entity

import "time"

type Chat struct {
	ID            int            `json:"id"`
	Users         []UserResponse `json:"users"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	LastMessage   string         `json:"last_message"`
	LastMessageAt time.Time      `json:"last_message_at"`
	IsRead        bool           `json:"is_read"`
}
