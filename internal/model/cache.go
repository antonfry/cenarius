package model

import "fmt"

type SecretCache struct {
	LoginWithPasswords []*LoginWithPassword `json:"login_and_passwords"`
	CreditCards        []*CreditCard        `json:"credit_cards"`
	SecretTexts        []*SecretText        `json:"secret_texts"`
	SecretFiles        []*SecretFile        `json:"secret_files"`
}

func (s *SecretCache) String() string {
	return fmt.Sprintf("LoginsPasswords:%v, CreditCards: %v, SecretTexts:%v, SecretFiles: %v", s.LoginWithPasswords, s.CreditCards, s.SecretTexts, s.SecretFiles)
}

func (c *SecretCache) Encrypt(key, iv string) error {
	for _, l := range c.LoginWithPasswords {
		if err := l.Encrypt(key, iv); err != nil {
			return err
		}
	}
	for _, c := range c.CreditCards {
		if err := c.Encrypt(key, iv); err != nil {
			return err
		}
	}
	for _, t := range c.SecretTexts {
		if err := t.Encrypt(key, iv); err != nil {
			return err
		}
	}
	for _, f := range c.SecretFiles {
		if err := f.Encrypt(key, iv); err != nil {
			return err
		}
	}
	return nil
}

func (c *SecretCache) Decrypt(key, iv string) error {
	for _, l := range c.LoginWithPasswords {
		if err := l.Decrypt(key, iv); err != nil {
			return err
		}
	}
	for _, c := range c.CreditCards {
		if err := c.Decrypt(key, iv); err != nil {
			return err
		}
	}
	for _, t := range c.SecretTexts {
		if err := t.Decrypt(key, iv); err != nil {
			return err
		}
	}
	for _, f := range c.SecretFiles {
		if err := f.Decrypt(key, iv); err != nil {
			return err
		}
	}
	return nil
}
