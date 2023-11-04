package memstore

import (
	"cenarius/internal/model"
	"cenarius/internal/store"
	"context"
)

type UserRepository struct {
	store *Store
	users map[int]*model.User
}

func (r *UserRepository) Ping() error {
	return nil
}

func (r *UserRepository) FindByLogin(ctx context.Context, login string) (*model.User, error) {
	for _, value := range r.users {
		if value.Login == login {
			return value, nil
		}
	}
	return nil, store.ErrRecordNotFound
}

func (r *UserRepository) FindByID(ctx context.Context, id int) (*model.User, error) {
	user, ok := r.users[id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}
	return user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	r.users[user.ID] = user
	return nil
}
