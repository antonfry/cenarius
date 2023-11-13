package model

type Encrypter interface {
	Encrypt(string, string) error
	Decrypt(string, string) error
}
type SecretData struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
	Meta   string `json:"meta"`
}
