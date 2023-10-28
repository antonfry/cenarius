package model

import (
	"bufio"
	"io/fs"
	"os"

	validation "github.com/go-ozzo/ozzo-validation"
	log "github.com/sirupsen/logrus"
)

type SecretFile struct {
	SecretData
	Path string `json:"path"`
}

func (s *SecretFile) Validate() error {
	return validation.ValidateStruct(
		s,
		validation.Field(&s.Path, validation.Required, validation.Length(10, 200)),
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

func (s *SecretFile) Encrypt(key, iv string) error {
	return nil
}

func (s *SecretFile) Decrypt(key, iv string) error {
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
