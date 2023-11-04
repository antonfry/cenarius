package memstore

import (
	"cenarius/internal/model"
	"cenarius/internal/store"
	"context"
)

type LoginWithPasswordRepository struct {
	store *Store
	m     map[int]map[int]*model.LoginWithPassword
}

func (r *LoginWithPasswordRepository) Ping() error {
	return nil
}

func (r *LoginWithPasswordRepository) SearchByName(ctx context.Context, name string, userID int) ([]*model.LoginWithPassword, error) {
	mm := make([]*model.LoginWithPassword, 0)
	for _, m := range r.m[userID] {
		if m.Name == name {
			mm = append(mm, m)
		}
	}
	return mm, nil
}

func (r *LoginWithPasswordRepository) GetByID(ctx context.Context, id, userID int) (*model.LoginWithPassword, error) {
	m, ok := r.m[userID][id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}
	return m, nil
}

func (r *LoginWithPasswordRepository) Add(ctx context.Context, m *model.LoginWithPassword) error {
	r.m[m.UserID][m.ID] = m
	return nil
}

func (r *LoginWithPasswordRepository) Update(ctx context.Context, m *model.LoginWithPassword) error {
	if err := r.Add(ctx, m); err != nil {
		return err
	}
	return nil
}

func (r *LoginWithPasswordRepository) Delete(ctx context.Context, id, userID int) error {
	delete(r.m[userID], id)
	return nil
}
