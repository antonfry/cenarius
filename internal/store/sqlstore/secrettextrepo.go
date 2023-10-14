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
		ctx, "INSERT INTO SecretText (user_id, name, meta, text) VALUES($1, $2, $3, $4) RETURNING id",
		m.UserId,
		m.Name,
		m.Meta,
		m.Text,
	).Scan(&m.ID); err != nil {
		return err
	}
	return nil
}

func (r *SecretTextRepository) Update(ctx context.Context, m *model.SecretText) error {
	if err := m.Validate(); err != nil {
		return err
	}
	if _, err := r.store.db.ExecContext(
		ctx, "UPDATE SecretText SET user_id=$1, name=$2, meta=$3, text=$4 WHERE id=$5",
		m.UserId,
		m.Name,
		m.Meta,
		m.Text,
		m.ID,
	); err != nil {
		return err
	}
	return nil
}

func (r *SecretTextRepository) Delete(ctx context.Context, id int) error {
	if _, err := r.store.db.ExecContext(ctx, "DELETE FROM SecretText WHERE id = $1", id); err != nil {
		return err
	}
	return nil
}

func (r *SecretTextRepository) GetByName(ctx context.Context, name string) (*model.SecretText, error) {
	m := &model.SecretText{}
	if err := r.store.db.QueryRowContext(
		ctx, "SELECT id, user_id, name, meta, text FROM SecretText WHERE name = $1", name,
	).Scan(&m.ID, &m.UserId, &m.Name, &m.Meta, &m.Text); err != nil {
		return nil, err
	}
	return m, nil
}

func (r *SecretTextRepository) GetByID(ctx context.Context, id int) (*model.SecretText, error) {
	m := &model.SecretText{}
	if err := r.store.db.QueryRowContext(
		ctx, "SELECT id, user_id, name, meta, text FROM SecretText WHERE id = $1", id,
	).Scan(&m.ID, &m.UserId, &m.Name, &m.Meta, &m.Text); err != nil {
		return nil, err
	}
	return m, nil
}
