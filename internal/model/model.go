package model

type SecretData struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
	Meta   string `json:"meta"`
}
