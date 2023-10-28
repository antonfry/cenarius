package model

import (
	"cenarius/internal/encrypt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type SecretText struct {
	SecretData
	Text string `json:"text"`
}

func (s *SecretText) Validate() error {
	return validation.ValidateStruct(
		s,
		validation.Field(&s.Text, validation.Required, is.ASCII),
	)
}

func (s *SecretText) Encrypt(key, iv string) error {
	encText, err := encrypt.AESEncrypted(s.Text, key, iv)
	if err != nil {
		return err
	}
	s.Text = encText
	return nil
}

func (s *SecretText) Decrypt(key, iv string) error {
	decText, err := encrypt.AESDecrypted(s.Text, key, iv)
	if err != nil {
		return err
	}
	s.Text = decText
	return nil
}
