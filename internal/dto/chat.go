package dto

type ChatResponse struct {
	ID        int    `json:"id"`
	User1ID   int    `json:"user1_id"`
	User2ID   int    `json:"user2_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
