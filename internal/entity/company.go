package entity

import "time"

type Company struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	OwnerId     int       `json:"ownerId"`
	CreatedAt   time.Time `json:"created_at"`
	PhotoUrl    string    `json:"photoUrl"`
	Description string    `json:"description"`
}
