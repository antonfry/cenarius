package memstore

import (
	"cenarius/internal/model"
	"cenarius/internal/store"
)

type Store struct {
	UserRepository *UserRepository
}

func New() *Store {
	return &Store{}
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
