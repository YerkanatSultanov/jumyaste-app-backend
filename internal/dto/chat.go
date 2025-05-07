package dto

type ChatResponse struct {
	ID        int    `json:"id"`
	User1ID   int    `json:"user1_id"`
	User2ID   int    `json:"user2_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CreateChatRequest struct {
	SecondUserId int `json:"second_user_id" binding:"required"`
}
