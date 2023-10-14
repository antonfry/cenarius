package sqlstore

import (
	"cenarius/internal/model"
	"context"
)

type SecretBinaryRepository struct {
	store *Store
}

func (r *SecretBinaryRepository) Ping() error {
	return r.store.db.Ping()
}

func (r *SecretBinaryRepository) Add(ctx context.Context, m *model.SecretBinary) error {
	if err := m.Validate(); err != nil {
		return err
	}
	if err := r.store.db.QueryRowContext(
		ctx, "INSERT INTO SecretBinary (user_id, name, meta, binary) VALUES($1, $2, $3, $4) RETURNING id",
		m.UserId,
		m.Name,
		m.Meta,
		m.Binary,
	).Scan(&m.ID); err != nil {
		return err
	}
	return nil
}

func (r *SecretBinaryRepository) Update(ctx context.Context, m *model.SecretBinary) error {
	if err := m.Validate(); err != nil {
		return err
	}
	if _, err := r.store.db.ExecContext(
		ctx, "UPDATE SecretText SET user_id=$1, name=$2, meta=$3, binary=$4 WHERE id=$5",
		m.UserId,
		m.Name,
		m.Meta,
		m.Binary,
		m.ID,
	); err != nil {
		return err
	}
	return nil
}

func (r *SecretBinaryRepository) Delete(ctx context.Context, id int) error {
	if _, err := r.store.db.ExecContext(ctx, "DELETE FROM SecretBinary WHERE id = $1", id); err != nil {
		return err
	}
	return nil
}

func (r *SecretBinaryRepository) GetByName(ctx context.Context, name string) (*model.SecretBinary, error) {
	m := &model.SecretBinary{}
	if err := r.store.db.QueryRowContext(
		ctx, "SELECT id, user_id, name, meta, binary FROM SecretBinary WHERE name = $1", name,
	).Scan(&m.ID, &m.UserId, &m.Name, &m.Meta, &m.Binary); err != nil {
		return nil, err
	}
	return m, nil
}

func (r *SecretBinaryRepository) GetByID(ctx context.Context, id int) (*model.SecretBinary, error) {
	m := &model.SecretBinary{}
	if err := r.store.db.QueryRowContext(
		ctx, "SELECT id, user_id, name, meta, binary FROM SecretBinary WHERE id = $1", id,
	).Scan(&m.ID, &m.UserId, &m.Name, &m.Meta, &m.Binary); err != nil {
		return nil, err
	}
	return m, nil
}
