package sqlstore

import (
	"cenarius/internal/model"
	"context"
)

type LoginWithPasswordRepository struct {
	store *Store
}

func (r *LoginWithPasswordRepository) Ping() error {
	return r.store.db.Ping()
}

func (r *LoginWithPasswordRepository) Add(ctx context.Context, lp *model.LoginWithPassword) error {
	if err := lp.Validate(); err != nil {
		return err
	}
	if err := r.store.db.QueryRowContext(
		ctx, "INSERT INTO lp (`user_id`, `name`, `meta`, `login`, `password`) VALUES($1, $2, $3, $4, $5) RETURNING id",
		lp.UserId,
		lp.Name,
		lp.Meta,
		lp.Login,
		lp.Password,
	).Scan(&lp.ID); err != nil {
		return err
	}
	return nil
}

func (r *LoginWithPasswordRepository) Delete(ctx context.Context, lp *model.LoginWithPassword) error {
	if _, err := r.store.db.ExecContext(ctx, "DELETE FROM lp WHERE id = $1", lp.ID); err != nil {
		return err
	}
	return nil
}

func (r *LoginWithPasswordRepository) GetByName(ctx context.Context, name string) (*model.LoginWithPassword, error) {
	lp := &model.LoginWithPassword{}
	if err := r.store.db.QueryRowContext(
		ctx, "SELECT `id`, `user_id`, `name`, `meta`, `login`, `password` FROM lp WHERE name = $1", name,
	).Scan(&lp.ID, &lp.UserId, &lp.Name, &lp.Meta, &lp.Login, &lp.Password); err != nil {
		return nil, err
	}
	return lp, nil
}

func (r *LoginWithPasswordRepository) GetByID(ctx context.Context, id int) (*model.LoginWithPassword, error) {
	lp := &model.LoginWithPassword{}
	if err := r.store.db.QueryRowContext(
		ctx, "SELECT `id`, `user_id`, `name`, `meta`, `login`, `password` FROM lp WHERE id = $1", id,
	).Scan(&lp.ID, &lp.UserId, &lp.Name, &lp.Meta, &lp.Login, &lp.Password); err != nil {
		return nil, err
	}
	return lp, nil
}
