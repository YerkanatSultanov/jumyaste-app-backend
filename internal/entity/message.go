package entity

import "time"

type MessageType string

const (
	TextMessage  MessageType = "text"
	ImageMessage MessageType = "image"
	VideoMessage MessageType = "video"
	AudioMessage MessageType = "audio"
	FileMessage  MessageType = "file"
)

type Message struct {
	ID        int         `gorm:"primaryKey" json:"id"`
	ChatID    int         `gorm:"index" json:"chat_id"`
	SenderID  int         `gorm:"index" json:"sender_id"`
	Type      MessageType `gorm:"type:varchar(255)" json:"type"`
	Content   *string     `json:"content,omitempty"`
	FileURL   *string     `json:"file_url,omitempty"`
	SeenBy    []User      `gorm:"many2many:message_seen_by;" json:"seen_by"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}
