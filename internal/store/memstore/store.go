package memstore

import (
	"cenarius/internal/model"
	"cenarius/internal/store"
)

type Store struct {
	LoginWithPasswordRepository *LoginWithPasswordRepository
	CreditCardRepository        *CreditCardRepository
	SecretTextRepository        *SecretTextRepository
	SecretFileRepository        *SecretFileRepository
	UserRepository              *UserRepository
}

func New() *Store {
	return &Store{}
}

func (s *Store) Close() {
}

func (s *Store) User() store.UserRepository {
	if s.UserRepository == nil {
		s.UserRepository = &UserRepository{
			store: s,
			users: make(map[int]*model.User),
		}
	}
	return s.UserRepository
}

func (s *Store) LoginWithPassword() store.LoginWithPasswordRepository {
	if s.LoginWithPasswordRepository == nil {
		s.LoginWithPasswordRepository = &LoginWithPasswordRepository{
			store: s,
			m:     make(map[int]map[int]*model.LoginWithPassword),
		}
	}
	return s.LoginWithPasswordRepository
}

func (s *Store) CreditCard() store.CreditCardRepository {
	if s.CreditCardRepository == nil {
		s.CreditCardRepository = &CreditCardRepository{
			store: s,
			m:     make(map[int]map[int]*model.CreditCard),
		}
	}
	return s.CreditCardRepository
}

func (s *Store) SecretText() store.SecretTextRepository {
	if s.SecretTextRepository == nil {
		s.SecretTextRepository = &SecretTextRepository{
			store: s,
			m:     make(map[int]map[int]*model.SecretText),
		}
	}
	return s.SecretTextRepository
}

func (s *Store) SecretFile() store.SecretFileRepository {
	if s.SecretFileRepository == nil {
		s.SecretFileRepository = &SecretFileRepository{
			store: s,
			m:     make(map[int]map[int]*model.SecretFile),
		}
	}
	return s.SecretFileRepository
}
