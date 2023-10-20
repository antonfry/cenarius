package sqlstore

import (
	"cenarius/internal/model"
	"cenarius/internal/store"
	"context"
	"database/sql"

	log "github.com/sirupsen/logrus"
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
		ctx, "UPDATE LoginWithPassword SET name=$1, meta=$2, login=$3, password=$4 WHERE user_id=$5 AND id=$6",
		m.Name,
		m.Meta,
		m.Login,
		m.Password,
		m.UserId,
		m.ID,
	); err != nil {
		return err
	}
	return nil
}

func (r *LoginWithPasswordRepository) Delete(ctx context.Context, m *model.LoginWithPassword) error {
	if _, err := r.store.db.ExecContext(ctx, "DELETE FROM LoginWithPassword WHERE id = $1 AND user_id=$2", m.ID, m.UserId); err != nil {
		return err
	}
	return nil
}

func (r *LoginWithPasswordRepository) SearchByName(ctx context.Context, name string, id int) ([]*model.LoginWithPassword, error) {
	mm := make([]*model.LoginWithPassword, 0)
	sql_string := "SELECT id, name, meta, login, password FROM LoginWithPassword WHERE user_id=$1"
	args := []any{id}
	if name != "" {
		sql_string += " AND name like $2"
		args = append(args, name)
	}
	log.Debugf(sql_string)
	rows, err := r.store.db.QueryContext(
		ctx, sql_string, args...,
	)
	if err != nil {
		log.Errorf("Unable to QueryContext in (r *LoginWithPasswordRepository) SearchByName: %v", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		m := &model.LoginWithPassword{}
		m.UserId = id
		err = rows.Scan(&m.ID, &m.Name, &m.Meta, &m.Login, &m.Password)
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

func (r *LoginWithPasswordRepository) GetByID(ctx context.Context, m *model.LoginWithPassword) (*model.LoginWithPassword, error) {
	if err := r.store.db.QueryRowContext(
		ctx, "SELECT name, meta, login, password FROM LoginWithPassword WHERE id = $1 AND user_id=$2", m.ID, m.UserId,
	).Scan(&m.Name, &m.Meta, &m.Login, &m.Password); err != nil {
		return nil, err
	}
	return m, nil
}
