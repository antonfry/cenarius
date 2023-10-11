package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type SecretData struct {
	ID     int    `json:"id"`
	UserId int    `json:"user_id"`
	Name   string `json:"name"`
	Meta   string `json:"meta"`
}

type LoginWithPassword struct {
	SecretData
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (s *LoginWithPassword) Validate() error {
	return validation.ValidateStruct(
		s,
		validation.Field(&s.Login, validation.Required, is.Alphanumeric),
		validation.Field(&s.Password, validation.Required, is.ASCII),
	)
}

type CreditCard struct {
	SecretData
	OwnerName     string `json:"owner_name"`
	OwnerLastName string `json:"owner_last_name"`
	Number        int    `json:"number"`
	CVC           int    `json:"cvc"`
}

func (s *CreditCard) Validate() error {
	return validation.ValidateStruct(
		s,
		validation.Field(&s.OwnerName, validation.Required, is.ASCII),
		validation.Field(&s.OwnerLastName, validation.Required, is.ASCII),
		validation.Field(&s.Number, validation.Required),
		validation.Field(&s.CVC, validation.Required, validation.Length(3, 3)),
	)
}

type SecretText struct {
	SecretData
	Text string `json:"text"`
}

func (s *SecretText) Validate() error {
	return validation.ValidateStruct(
		s,
		validation.Field(&s.Text, validation.Required, is.Alphanumeric),
	)
}

type SecretBinary struct {
	SecretData
	Binary []byte `json:"binary"`
}

func (s *SecretBinary) Validate() error {
	return validation.ValidateStruct(
		s,
		validation.Field(&s.Binary, validation.Required),
	)
}
