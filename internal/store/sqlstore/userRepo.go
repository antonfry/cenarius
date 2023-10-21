package sqlstore

import (
	"cenarius/internal/model"
	"context"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Ping() error {
	return r.store.db.Ping()
}

func (r *UserRepository) FindByLogin(ctx context.Context, login string) (*model.User, error) {
	user := &model.User{}
	if err := r.store.db.QueryRowContext(
		ctx, "SELECT id, login, encrypted_password FROM users WHERE login = $1", login,
	).Scan(&user.ID, &user.Login, &user.EncryptedPassword); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id int) (*model.User, error) {
	user := &model.User{}
	if err := r.store.db.QueryRowContext(
		ctx, "SELECT id, login, encrypted_password FROM users WHERE id = $1", id,
	).Scan(&user.ID, &user.Login, &user.EncryptedPassword); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	if err := user.Validate(); err != nil {
		return err
	}
	if err := user.BeforeCreate(); err != nil {
		return err
	}
	if err := r.store.db.QueryRowContext(
		ctx, "INSERT INTO users (login, encrypted_password) VALUES($1, $2) RETURNING id",
		user.Login,
		user.EncryptedPassword,
	).Scan(&user.ID); err != nil {
		return err
	}
	return nil
}
