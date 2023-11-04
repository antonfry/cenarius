package memstore

import (
	"cenarius/internal/model"
	"cenarius/internal/store"
	"context"
)

type SecretFileRepository struct {
	store *Store
	m     map[int]map[int]*model.SecretFile
}

func (r *SecretFileRepository) Ping() error {
	return nil
}

func (r *SecretFileRepository) SearchByName(ctx context.Context, name string, userID int) ([]*model.SecretFile, error) {
	mm := make([]*model.SecretFile, 0)
	for _, m := range r.m[userID] {
		if m.Name == name {
			mm = append(mm, m)
		}
	}
	return mm, nil
}

func (r *SecretFileRepository) GetByID(ctx context.Context, id, userID int) (*model.SecretFile, error) {
	m, ok := r.m[userID][id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}
	return m, nil
}

func (r *SecretFileRepository) Add(ctx context.Context, m *model.SecretFile) error {
	r.m[m.UserID][m.ID] = m
	return nil
}

func (r *SecretFileRepository) Update(ctx context.Context, m *model.SecretFile) error {
	if err := r.Add(ctx, m); err != nil {
		return err
	}
	return nil
}

func (r *SecretFileRepository) Delete(ctx context.Context, id, userID int) error {
	delete(r.m[userID], id)
	return nil
}
