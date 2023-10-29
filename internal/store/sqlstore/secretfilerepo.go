package sqlstore

import (
	"cenarius/internal/model"
	"cenarius/internal/store"
	"context"
	"database/sql"
)

type SecretFileRepository struct {
	store *Store
}

func (r *SecretFileRepository) Ping() error {
	return r.store.db.Ping()
}

func (r *SecretFileRepository) Add(ctx context.Context, m *model.SecretFile) error {
	if err := r.store.db.QueryRowContext(
		ctx, "INSERT INTO SecretFile (user_id, name, meta, path) VALUES($1, $2, $3, $4) RETURNING id",
		m.UserId,
		m.Name,
		m.Meta,
		m.Path,
	).Scan(&m.ID); err != nil {
		return err
	}
	return nil
}

func (r *SecretFileRepository) Update(ctx context.Context, m *model.SecretFile) error {
	if _, err := r.store.db.ExecContext(
		ctx, "UPDATE SecretFile SET name=$1, meta=$2 WHERE id=$3 AND user_id=$4",
		m.Name,
		m.Meta,
		m.ID,
		m.UserId,
	); err != nil {
		return err
	}
	return nil
}

func (r *SecretFileRepository) Delete(ctx context.Context, m *model.SecretFile) error {
	if _, err := r.store.db.ExecContext(ctx, "DELETE FROM SecretFile WHERE id = $1 AND user_id = $2", m.ID, m.UserId); err != nil {
		return err
	}
	return nil
}

func (r *SecretFileRepository) SearchByName(ctx context.Context, name string, id int) ([]*model.SecretFile, error) {
	mm := make([]*model.SecretFile, 0)
	sql_string := "SELECT id, name, meta, path FROM SecretFile WHERE user_id=$1"
	args := []any{id}
	if name != "" {
		sql_string += " AND name like $2"
		args = append(args, name)
	}
	rows, err := r.store.db.QueryContext(
		ctx, sql_string, args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		m := &model.SecretFile{}
		m.UserId = id
		err = rows.Scan(&m.ID, &m.Name, &m.Meta, &m.Path)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, store.ErrRecordNotFound
			}
			return nil, err
		}
		mm = append(mm, m)
	}
	if rows.Err() != nil {
		return nil, store.ErrUnableToGetRows
	}
	return mm, nil
}

func (r *SecretFileRepository) GetByID(ctx context.Context, id, user_id int) (*model.SecretFile, error) {
	m := &model.SecretFile{}
	if err := r.store.db.QueryRowContext(
		ctx, "SELECT name, meta, path FROM SecretFile WHERE id = $1 AND user_id = $2", id, user_id,
	).Scan(&m.Name, &m.Meta, &m.Path); err != nil {
		return nil, err
	}
	m.ID = id
	m.UserId = user_id
	return m, nil
}
