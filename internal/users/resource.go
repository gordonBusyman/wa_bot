package users

// Resource represents a user.
type Resource struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	ChatID int64  `json:"chat_id"`
}
