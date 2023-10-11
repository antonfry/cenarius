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
		ctx, "INSERT INTO credit_cards (`user_id`, `name`, `meta`, `owner_name`, `owner_last_name`, `number`, `cvc`) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id",
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

func (r *CreditCardRepository) Delete(ctx context.Context, m *model.CreditCard) error {
	if _, err := r.store.db.ExecContext(ctx, "DELETE FROM credit_cards WHERE id = $1", m.ID); err != nil {
		return err
	}
	return nil
}

func (r *CreditCardRepository) GetByName(ctx context.Context, name string) (*model.CreditCard, error) {
	m := &model.CreditCard{}
	if err := r.store.db.QueryRowContext(
		ctx, "SELECT `id`, `user_id`, `name`, `meta`, `login`, `password` FROM credit_cards WHERE name = $1", name,
	).Scan(&m.ID, &m.UserId, &m.Name, &m.Meta, &m.OwnerName, &m.OwnerLastName, &m.Number, &m.CVC); err != nil {
		return nil, err
	}
	return m, nil
}

func (r *CreditCardRepository) GetByID(ctx context.Context, id int) (*model.CreditCard, error) {
	m := &model.CreditCard{}
	if err := r.store.db.QueryRowContext(
		ctx, "SELECT `id`, `user_id`, `name`, `meta`, `login`, `password` FROM credit_cards WHERE id = $1", id,
	).Scan(&m.ID, &m.UserId, &m.Name, &m.Meta, &m.OwnerName, &m.OwnerLastName, &m.Number, &m.CVC); err != nil {
		return nil, err
	}
	return m, nil
}
