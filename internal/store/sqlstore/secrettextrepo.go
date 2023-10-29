package sqlstore

import (
	"cenarius/internal/model"
	"cenarius/internal/store"
	"context"
	"database/sql"
)

type SecretTextRepository struct {
	store *Store
}

func (r *SecretTextRepository) Ping() error {
	return r.store.db.Ping()
}

func (r *SecretTextRepository) Add(ctx context.Context, m *model.SecretText) error {
	if err := r.store.db.QueryRowContext(
		ctx, "INSERT INTO SecretText (user_id, name, meta, text) VALUES($1, $2, $3, $4) RETURNING id",
		m.UserID,
		m.Name,
		m.Meta,
		m.Text,
	).Scan(&m.ID); err != nil {
		return err
	}
	return nil
}

func (r *SecretTextRepository) Update(ctx context.Context, m *model.SecretText) error {
	if _, err := r.store.db.ExecContext(
		ctx, "UPDATE SecretText SET user_id=$1, name=$2, meta=$3, text=$4 WHERE id=$5",
		m.UserID,
		m.Name,
		m.Meta,
		m.Text,
		m.ID,
	); err != nil {
		return err
	}
	return nil
}

func (r *SecretTextRepository) Delete(ctx context.Context, m *model.SecretText) error {
	if _, err := r.store.db.ExecContext(ctx, "DELETE FROM SecretText WHERE id = $1 AND user_id = $2", m.ID, m.UserID); err != nil {
		return err
	}
	return nil
}

func (r *SecretTextRepository) SearchByName(ctx context.Context, name string, id int) ([]*model.SecretText, error) {
	mm := make([]*model.SecretText, 0)
	sqlString := "SELECT id, name, meta, text FROM SecretText WHERE user_id=$1"
	args := []any{id}
	if name != "" {
		sqlString += " AND name like $2"
		args = append(args, name)
	}
	rows, err := r.store.db.QueryContext(
		ctx, sqlString, args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		m := &model.SecretText{}
		m.UserID = id
		err = rows.Scan(&m.ID, &m.Name, &m.Meta, &m.Text)
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

func (r *SecretTextRepository) GetByID(ctx context.Context, m *model.SecretText) (*model.SecretText, error) {
	if err := r.store.db.QueryRowContext(
		ctx, "SELECT name, meta, text FROM SecretText WHERE id = $1 AND user_id = $2", m.ID, m.UserID,
	).Scan(&m.Name, &m.Meta, &m.Text); err != nil {
		return nil, err
	}
	return m, nil
}
