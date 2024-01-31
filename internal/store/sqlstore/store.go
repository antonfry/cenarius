package sqlstore

import (
	"cenarius/internal/store"
	"database/sql"

	_ "github.com/jackc/pgx"
)

type Store struct {
	db                          *sql.DB
	LoginWithPasswordRepository *LoginWithPasswordRepository
	CreditCardRepository        *CreditCardRepository
	SecretTextRepository        *SecretTextRepository
	SecretFileRepository        *SecretFileRepository
	UserRepository              *UserRepository
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) Close() {
	s.db.Close()
}

func (s *Store) LoginWithPassword() store.LoginWithPasswordRepository {
	if s.LoginWithPasswordRepository == nil {
		s.LoginWithPasswordRepository = &LoginWithPasswordRepository{
			store: s,
		}
	}
	return s.LoginWithPasswordRepository
}

func (s *Store) CreditCard() store.CreditCardRepository {
	if s.CreditCardRepository == nil {
		s.CreditCardRepository = &CreditCardRepository{
			store: s,
		}
	}
	return s.CreditCardRepository
}

func (s *Store) SecretText() store.SecretTextRepository {
	if s.SecretTextRepository == nil {
		s.SecretTextRepository = &SecretTextRepository{
			store: s,
		}
	}
	return s.SecretTextRepository
}

func (s *Store) SecretFile() store.SecretFileRepository {
	if s.SecretFileRepository == nil {
		s.SecretFileRepository = &SecretFileRepository{
			store: s,
		}
	}
	return s.SecretFileRepository
}

func (s *Store) User() store.UserRepository {
	if s.UserRepository == nil {
		s.UserRepository = &UserRepository{
			store: s,
		}
	}
	return s.UserRepository
}
