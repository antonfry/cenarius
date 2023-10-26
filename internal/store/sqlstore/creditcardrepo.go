package sqlstore

import (
	"cenarius/internal/model"
	"cenarius/internal/store"
	"context"
	"database/sql"
)

type CreditCardRepository struct {
	store *Store
}

func (r *CreditCardRepository) Ping() error {
	return r.store.db.Ping()
}

func (r *CreditCardRepository) Add(ctx context.Context, m *model.CreditCard) error {
	if err := m.Validate(); err != nil {
		return err
	}
	if err := r.store.db.QueryRowContext(
		ctx, "INSERT INTO CreditCard (user_id, name, meta, owner_name, owner_last_name, number, cvc) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		m.UserId,
		m.Name,
		m.Meta,
		m.OwnerName,
		m.OwnerLastName,
		m.Number,
		m.CVC,
	).Scan(&m.ID); err != nil {
		return err
	}
	return nil
}

func (r *CreditCardRepository) Update(ctx context.Context, m *model.CreditCard) error {
	if err := m.Validate(); err != nil {
		return err
	}
	if _, err := r.store.db.ExecContext(
		ctx, "UPDATE CreditCard SET user_id=$1, name=$2, meta=$3, owner_name=$4, owner_last_name=$5, number=$6, cvc=$7  WHERE id=$8",
		m.UserId,
		m.Name,
		m.Meta,
		m.OwnerName,
		m.OwnerLastName,
		m.Number,
		m.CVC,
		m.ID,
	); err != nil {
		return err
	}
	return nil
}

func (r *CreditCardRepository) Delete(ctx context.Context, m *model.CreditCard) error {
	if _, err := r.store.db.ExecContext(ctx, "DELETE FROM CreditCard WHERE id = $1 AND user_id = $2", m.ID, m.UserId); err != nil {
		return err
	}
	return nil
}

func (r *CreditCardRepository) SearchByName(ctx context.Context, name string, id int) ([]*model.CreditCard, error) {
	mm := make([]*model.CreditCard, 0)
	sql_string := "SELECT id, name, meta, owner_name, owner_last_name, number, cvc FROM CreditCard WHERE user_id=$1"
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
		m := &model.CreditCard{}
		m.UserId = id
		err = rows.Scan(&m.ID, &m.Name, &m.Meta, &m.OwnerName, &m.OwnerLastName, &m.Number, &m.CVC)
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

func (r *CreditCardRepository) GetByID(ctx context.Context, m *model.CreditCard) (*model.CreditCard, error) {
	if err := r.store.db.QueryRowContext(
		ctx, "SELECT name, meta, owner_name, owner_last_name, number, cvc FROM CreditCard WHERE id = $1 AND user_id = $2", m.ID, m.UserId,
	).Scan(&m.Name, &m.Meta, &m.OwnerName, &m.OwnerLastName, &m.Number, &m.CVC); err != nil {
		return nil, err
	}
	return m, nil
}
