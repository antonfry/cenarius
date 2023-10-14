package sqlstore

import (
	"cenarius/internal/model"
	"context"
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

func (r *CreditCardRepository) Delete(ctx context.Context, id int) error {
	if _, err := r.store.db.ExecContext(ctx, "DELETE FROM CreditCard WHERE id = $1", id); err != nil {
		return err
	}
	return nil
}

func (r *CreditCardRepository) GetByName(ctx context.Context, name string) (*model.CreditCard, error) {
	m := &model.CreditCard{}
	if err := r.store.db.QueryRowContext(
		ctx, "SELECT id, user_id, name, meta, login, password FROM CreditCard WHERE name = $1", name,
	).Scan(&m.ID, &m.UserId, &m.Name, &m.Meta, &m.OwnerName, &m.OwnerLastName, &m.Number, &m.CVC); err != nil {
		return nil, err
	}
	return m, nil
}

func (r *CreditCardRepository) GetByID(ctx context.Context, id int) (*model.CreditCard, error) {
	m := &model.CreditCard{}
	if err := r.store.db.QueryRowContext(
		ctx, "SELECT id, user_id, name, meta, login, password FROM CreditCard WHERE id = $1", id,
	).Scan(&m.ID, &m.UserId, &m.Name, &m.Meta, &m.OwnerName, &m.OwnerLastName, &m.Number, &m.CVC); err != nil {
		return nil, err
	}
	return m, nil
}
