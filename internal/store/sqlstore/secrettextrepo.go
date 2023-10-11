package sqlstore

import (
	"cenarius/internal/model"
	"context"
)

type SecretTextRepository struct {
	store *Store
}

func (r *SecretTextRepository) Ping() error {
	return r.store.db.Ping()
}

func (r *SecretTextRepository) Add(ctx context.Context, m *model.SecretText) error {
	if err := m.Validate(); err != nil {
		return err
	}
	if err := r.store.db.QueryRowContext(
		ctx, "INSERT INTO credit_cards (`user_id`, `name`, `meta`, `text`) VALUES($1, $2, $3, $4) RETURNING id",
		m.UserId,
		m.Name,
		m.Meta,
		m.Text,
	).Scan(&m.ID); err != nil {
		return err
	}
	return nil
}

func (r *SecretTextRepository) Delete(ctx context.Context, m *model.SecretText) error {
	if _, err := r.store.db.ExecContext(ctx, "DELETE FROM credit_cards WHERE id = $1", m.ID); err != nil {
		return err
	}
	return nil
}

func (r *SecretTextRepository) GetByName(ctx context.Context, name string) (*model.SecretText, error) {
	m := &model.SecretText{}
	if err := r.store.db.QueryRowContext(
		ctx, "SELECT `id`, `user_id`, `name`, `meta`, `login`, `password` FROM credit_cards WHERE name = $1", name,
	).Scan(&m.ID, &m.UserId, &m.Name, &m.Meta, &m.Text); err != nil {
		return nil, err
	}
	return m, nil
}

func (r *SecretTextRepository) GetByID(ctx context.Context, id int) (*model.SecretText, error) {
	m := &model.SecretText{}
	if err := r.store.db.QueryRowContext(
		ctx, "SELECT `id`, `user_id`, `name`, `meta`, `login`, `password` FROM credit_cards WHERE id = $1", id,
	).Scan(&m.ID, &m.UserId, &m.Name, &m.Meta, &m.Text); err != nil {
		return nil, err
	}
	return m, nil
}
