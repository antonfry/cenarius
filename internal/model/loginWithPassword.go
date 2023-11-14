package model

import (
	"cenarius/internal/encrypt"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type LoginWithPassword struct {
	SecretData
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (s *LoginWithPassword) String() string {
	return fmt.Sprintf("ID: %d, Name: %s, Login: %s, Password: %s, Meta: %s", s.ID, s.Name, s.Login, s.Password, s.Meta)
}

func (s *LoginWithPassword) Validate() error {
	return validation.ValidateStruct(
		s,
		validation.Field(&s.Login, validation.Required, is.Alphanumeric),
		validation.Field(&s.Password, validation.Required, is.Alphanumeric),
	)
}

func (s *LoginWithPassword) Encrypt(key, iv string) error {
	encLogin, err := encrypt.AESEncrypted(s.Login, key, iv)
	if err != nil {
		return err
	}
	encPassword, err := encrypt.AESEncrypted(s.Password, key, iv)
	if err != nil {
		return err
	}
	s.Login = encLogin
	s.Password = encPassword
	return nil
}

func (s *LoginWithPassword) Decrypt(key, iv string) error {
	decLogin, err := encrypt.AESDecrypted(s.Login, key, iv)
	if err != nil {
		return err
	}
	decPassword, err := encrypt.AESDecrypted(s.Password, key, iv)
	if err != nil {
		return err
	}
	s.Login = decLogin
	s.Password = decPassword
	return nil
}
