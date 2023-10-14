package sqlstore

import (
	"cenarius/internal/model"
	"context"
)

type LoginWithPasswordRepository struct {
	store *Store
}

func (r *LoginWithPasswordRepository) Ping() error {
	return r.store.db.Ping()
}

func (r *LoginWithPasswordRepository) Add(ctx context.Context, m *model.LoginWithPassword) error {
	if err := m.Validate(); err != nil {
		return err
	}
	if err := r.store.db.QueryRowContext(
		ctx, "INSERT INTO LoginWithPassword (user_id, name, meta, login, password) VALUES($1, $2, $3, $4, $5) RETURNING id",
		m.UserId,
		m.Name,
		m.Meta,
		m.Login,
		m.Password,
	).Scan(&m.ID); err != nil {
		return err
	}
	return nil
}

func (r *LoginWithPasswordRepository) Update(ctx context.Context, m *model.LoginWithPassword) error {
	if err := m.Validate(); err != nil {
		return err
	}
	if _, err := r.store.db.ExecContext(
		ctx, "UPDATE LoginWithPassword SET user_id=$1, name=$2, meta=$3, login=$4, password=$5 WHERE id=$6",
		m.UserId,
		m.Name,
		m.Meta,
		m.Login,
		m.Password,
		m.ID,
	); err != nil {
		return err
	}
	return nil
}

func (r *LoginWithPasswordRepository) Delete(ctx context.Context, id int) error {
	if _, err := r.store.db.ExecContext(ctx, "DELETE FROM LoginWithPassword WHERE id = $1", id); err != nil {
		return err
	}
	return nil
}

func (r *LoginWithPasswordRepository) GetByName(ctx context.Context, name string) (*model.LoginWithPassword, error) {
	m := &model.LoginWithPassword{}
	if err := r.store.db.QueryRowContext(
		ctx, "SELECT id, user_id, name, meta, login, password FROM LoginWithPassword WHERE name = $1", name,
	).Scan(&m.ID, &m.UserId, &m.Name, &m.Meta, &m.Login, &m.Password); err != nil {
		return nil, err
	}
	return m, nil
}

func (r *LoginWithPasswordRepository) GetByID(ctx context.Context, id int) (*model.LoginWithPassword, error) {
	m := &model.LoginWithPassword{}
	if err := r.store.db.QueryRowContext(
		ctx, "SELECT id, user_id, name, meta, login, password FROM LoginWithPassword WHERE id = $1", id,
	).Scan(&m.ID, &m.UserId, &m.Name, &m.Meta, &m.Login, &m.Password); err != nil {
		return nil, err
	}
	return m, nil
}
