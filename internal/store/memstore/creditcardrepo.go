package memstore

import (
	"cenarius/internal/model"
	"cenarius/internal/store"
	"context"
)

type CreditCardRepository struct {
	store *Store
	m     map[int]map[int]*model.CreditCard
}

func (r *CreditCardRepository) Ping() error {
	return nil
}

func (r *CreditCardRepository) SearchByName(ctx context.Context, name string, userID int) ([]*model.CreditCard, error) {
	mm := make([]*model.CreditCard, 0)
	for _, m := range r.m[userID] {
		if m.Name == name {
			mm = append(mm, m)
		}
	}
	return mm, nil
}

func (r *CreditCardRepository) GetByID(ctx context.Context, id, userID int) (*model.CreditCard, error) {
	m, ok := r.m[userID][id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}
	return m, nil
}

func (r *CreditCardRepository) Add(ctx context.Context, m *model.CreditCard) error {
	r.m[m.UserID][m.ID] = m
	return nil
}

func (r *CreditCardRepository) Update(ctx context.Context, m *model.CreditCard) error {
	if err := r.Add(ctx, m); err != nil {
		return err
	}
	return nil
}

func (r *CreditCardRepository) Delete(ctx context.Context, id, userID int) error {
	delete(r.m[userID], id)
	return nil
}
