package model

import (
	"bufio"
	"cenarius/internal/encrypt"
	"io/fs"
	"os"
	"strconv"

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

type CreditCard struct {
	SecretData
	OwnerName     string `json:"owner_name"`
	OwnerLastName string `json:"owner_last_name"`
	EncNumber     string `json:"encnumber"`
	Number        int    `json:"number"`
	CVC           int    `json:"cvc"`
	EncCVC        string `json:"enccvc"`
}

func (s *CreditCard) Validate() error {
	return validation.ValidateStruct(
		s,
		validation.Field(&s.OwnerName, validation.Required, is.ASCII),
		validation.Field(&s.OwnerLastName, validation.Required, is.ASCII),
		validation.Field(&s.Number, validation.Required, validation.Min(1000000000000000), validation.Max(9999999999999999)),
		validation.Field(&s.CVC, validation.Required, validation.Min(100), validation.Max(999)),
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
	encNumber, err := encrypt.AESEncrypted(strconv.Itoa(s.Number), key, iv)
	if err != nil {
		return err
	}
	encCVC, err := encrypt.AESEncrypted(strconv.Itoa(s.CVC), key, iv)
	if err != nil {
		return err
	}
	s.OwnerName = encOwnerName
	s.OwnerLastName = encOwnerLastName
	s.EncNumber = encNumber
	s.Number = 0
	s.EncCVC = encCVC
	s.CVC = 0
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
	decNumber, err := encrypt.AESDecrypted(s.EncNumber, key, iv)
	if err != nil {
		return err
	}
	decCVC, err := encrypt.AESDecrypted(s.EncCVC, key, iv)
	if err != nil {
		return err
	}
	s.OwnerName = decOwnerName
	s.OwnerLastName = decOwnerLastName
	s.Number, err = strconv.Atoi(decNumber)
	if err != nil {
		return err
	}
	s.CVC, err = strconv.Atoi(decCVC)
	if err != nil {
		return err
	}
	return nil
}

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

type SecretFile struct {
	SecretData
	Path string `json:"path"`
}

func (s *SecretFile) Validate() error {
	return validation.ValidateStruct(
		s,
		validation.Field(&s.Path, validation.Required, validation.Length(1, 30)),
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
