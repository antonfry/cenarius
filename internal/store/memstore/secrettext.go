package memstore

import (
	"cenarius/internal/model"
	"cenarius/internal/store"
	"context"
)

type SecretTextRepository struct {
	store *Store
	m     map[int]map[int]*model.SecretText
}

func (r *SecretTextRepository) Ping() error {
	return nil
}

func (r *SecretTextRepository) SearchByName(ctx context.Context, name string, userID int) ([]*model.SecretText, error) {
	mm := make([]*model.SecretText, 0)
	for _, m := range r.m[userID] {
		if m.Name == name {
			mm = append(mm, m)
		}
	}
	return mm, nil
}

func (r *SecretTextRepository) GetByID(ctx context.Context, id, userID int) (*model.SecretText, error) {
	m, ok := r.m[userID][id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}
	return m, nil
}

func (r *SecretTextRepository) Add(ctx context.Context, m *model.SecretText) error {
	r.m[m.UserID][m.ID] = m
	return nil
}

func (r *SecretTextRepository) Update(ctx context.Context, m *model.SecretText) error {
	if err := r.Add(ctx, m); err != nil {
		return err
	}
	return nil
}

func (r *SecretTextRepository) Delete(ctx context.Context, id, userID int) error {
	delete(r.m[userID], id)
	return nil
}
