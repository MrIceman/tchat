package types

type Message struct {
	UserID      string `json:"userID"`
	DisplayName string `json:"displayName"`
	Content     string `json:"content"`
	CreatedAt   string `json:"createdAt"`
}
