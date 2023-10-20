package model

import (
	"bufio"
	"io/fs"
	"os"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	log "github.com/sirupsen/logrus"
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
		validation.Field(&s.Password, validation.Required, is.Alphanumeric),
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

type SecretFile struct {
	SecretData
	Path string `json:"path"`
}

func (s *SecretFile) Validate() error {
	return validation.ValidateStruct(
		s,
		validation.Field(&s.Path, validation.Required),
	)
}

func (s *SecretFile) Remove() error {
	if _, err := s.stat(); err != nil {
		return err
	}
	if err := os.Remove(s.Path); err != nil {
		log.Errorf("Unable to remove file %v: %v", s, err)
		return err
	}
	return nil
}

func (s *SecretFile) Get() ([]byte, error) {
	if _, err := s.stat(); err != nil {
		return nil, err
	}
	file, err := os.Open(s.Path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	fileStat, err := file.Stat()
	if err != nil {
		log.Errorf("Unable to get file stats %v: %v", s, err)
		return nil, err
	}
	reader := bufio.NewReader(file)

	buf := make([]byte, fileStat.Size())
	_, err = reader.Read(buf)
	if err != nil {
		log.Errorf("Unable to read data into buffer: %s, %v", s.Path, err)
		return nil, err
	}
	return buf, nil
}

func (s *SecretFile) stat() (fs.FileInfo, error) {
	stat, err := os.Stat(s.Path)
	if err != nil {
		log.Errorf("Unable to get file stat %v: %v", s, err)
		return nil, err
	}
	return stat, nil
}
