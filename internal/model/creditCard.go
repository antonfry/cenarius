package model

import (
	"cenarius/internal/encrypt"
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type CreditCard struct {
	SecretData
	OwnerName     string `json:"owner_name"`
	OwnerLastName string `json:"owner_last_name"`
	Number        string `json:"number"`
	CVC           string `json:"cvc"`
}

func (s *CreditCard) String() string {
	return fmt.Sprintf(
		"ID: %d, Name: %s, OwnerName: %s, OwnerLastName: %s, Number: %s, CVC: %s, Meta: %s",
		s.ID, s.Name, s.OwnerName, s.OwnerLastName, s.Number, s.CVC, s.Meta,
	)
}

func (s *CreditCard) Validate() error {
	return validation.ValidateStruct(
		s,
		validation.Field(&s.OwnerName, validation.Required, is.ASCII),
		validation.Field(&s.OwnerLastName, validation.Required, is.ASCII),
		validation.Field(&s.Number, validation.Required, is.CreditCard),
		validation.Field(&s.CVC, validation.Required),
	)
}

func (s *CreditCard) Encrypt(key, iv string) error {
	encOwnerName, err := encrypt.AESEncrypted(s.OwnerName, key, iv)
	if err != nil {
		return err
	}
	encOwnerLastName, err := encrypt.AESEncrypted(s.OwnerLastName, key, iv)
	if err != nil {
		return err
	}
	encNumber, err := encrypt.AESEncrypted(fmt.Sprintf("%s:%s:%s", s.Number, key, iv), key, iv)
	if err != nil {
		return err
	}
	encCVC, err := encrypt.AESEncrypted(s.CVC, key, iv)
	if err != nil {
		return err
	}
	s.OwnerName = encOwnerName
	s.OwnerLastName = encOwnerLastName
	s.Number = encNumber
	s.CVC = encCVC
	return nil
}

func (s *CreditCard) Decrypt(key, iv string) error {
	decOwnerName, err := encrypt.AESDecrypted(s.OwnerName, key, iv)
	if err != nil {
		return err
	}
	decOwnerLastName, err := encrypt.AESDecrypted(s.OwnerLastName, key, iv)
	if err != nil {
		return err
	}
	decNumber, err := encrypt.AESDecrypted(s.Number, key, iv)
	if err != nil {
		return err
	}
	decCVC, err := encrypt.AESDecrypted(s.CVC, key, iv)
	if err != nil {
		return err
	}
	s.OwnerName = decOwnerName
	s.OwnerLastName = decOwnerLastName
	s.Number = strings.Split(decNumber, ":")[0]
	s.CVC = decCVC
	return nil
}
