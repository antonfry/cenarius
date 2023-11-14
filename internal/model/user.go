package model

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                int    `json:"id"`
	Login             string `json:"login"`
	Password          string `json:"password,omitempty"`
	EncryptedPassword string `json:"encrypted_password,omitempty"`
}

func (u *User) String() string {
	return fmt.Sprintf("ID: %d, Login: %s, Password: %s, EncryptedPassword: %s", u.ID, u.Login, u.Password, u.EncryptedPassword)
}

func (u *User) Sanitaze() {
	u.Password = ""
}

func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)) == nil
}

func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Login, validation.Required, is.Alphanumeric),
		validation.Field(
			&u.Password,
			validation.By(requiredIf(u.EncryptedPassword == "")),
			validation.Length(8, 32),
		),
	)
}

func (u *User) BeforeCreate() error {
	if len(u.Password) > 0 {
		enc, err := HashFromString(u.Password)
		if err != nil {
			return err
		}
		u.EncryptedPassword = enc
	}
	u.Sanitaze()
	return nil
}

func HashFromString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
