package model

import (
	"cenarius/internal/encrypt"
	"fmt"
	"io/fs"
	"os"

	validation "github.com/go-ozzo/ozzo-validation"
	log "github.com/sirupsen/logrus"
)

type SecretFile struct {
	SecretData
	Path string `json:"path"`
}

func (s *SecretFile) String() string {
	return fmt.Sprintf("ID: %d, User ID: %d, Name: %s, Meta: %s, Path: %s", s.ID, s.UserID, s.Name, s.Meta, s.Path)
}

func (s *SecretFile) Validate() error {
	return validation.ValidateStruct(
		s,
		validation.Field(&s.Path, validation.Required, validation.Length(10, 200)),
		validation.Field(&s.UserID, validation.Required, validation.Min(1)),
	)
}

func (s *SecretFile) Encrypt(key, iv string) error {
	encPath, err := encrypt.AESEncrypted(s.Path, key, iv)
	if err != nil {
		return err
	}
	s.Path = encPath
	return nil
}

func (s *SecretFile) Decrypt(key, iv string) error {
	decPath, err := encrypt.AESDecrypted(s.Path, key, iv)
	if err != nil {
		return err
	}
	s.Path = decPath
	return nil
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

func (s *SecretFile) stat() (fs.FileInfo, error) {
	stat, err := os.Stat(s.Path)
	if err != nil {
		log.Errorf("Unable to get file stat %v: %v", s, err)
		return nil, err
	}
	return stat, nil
}
