package entity

import "time"

type Resume struct {
	ID              int                    `json:"id"`
	UserID          int                    `json:"user_id"`
	FullName        string                 `json:"full_name"`
	DesiredPosition string                 `json:"desired_position"`
	Skills          []string               `json:"skills"`
	City            string                 `json:"city"`
	About           string                 `json:"about"`
	ParsedData      map[string]interface{} `json:"parsed_data"`
	CreatedAt       time.Time              `json:"created_at"`
}
