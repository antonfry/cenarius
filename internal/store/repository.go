package store

import (
	"cenarius/internal/model"
	"context"
)

type SecretDataDeleter interface {
	Delete(context.Context, int, int) error
}

type UserRepository interface {
	FindByID(context.Context, int) (*model.User, error)
	FindByLogin(context.Context, string) (*model.User, error)
	Create(context.Context, *model.User) error
}

type LoginWithPasswordRepository interface {
	SecretDataDeleter
	SearchByName(context.Context, string, int) ([]*model.LoginWithPassword, error)
	GetByID(context.Context, int, int) (*model.LoginWithPassword, error)
	Add(context.Context, *model.LoginWithPassword) error
	Update(context.Context, *model.LoginWithPassword) error
}

type CreditCardRepository interface {
	SecretDataDeleter
	SearchByName(context.Context, string, int) ([]*model.CreditCard, error)
	GetByID(context.Context, int, int) (*model.CreditCard, error)
	Add(context.Context, *model.CreditCard) error
	Update(context.Context, *model.CreditCard) error
}

type SecretTextRepository interface {
	SecretDataDeleter
	SearchByName(context.Context, string, int) ([]*model.SecretText, error)
	GetByID(context.Context, int, int) (*model.SecretText, error)
	Add(context.Context, *model.SecretText) error
	Update(context.Context, *model.SecretText) error
}

type SecretFileRepository interface {
	SecretDataDeleter
	SearchByName(context.Context, string, int) ([]*model.SecretFile, error)
	GetByID(context.Context, int, int) (*model.SecretFile, error)
	Add(context.Context, *model.SecretFile) error
	Update(context.Context, *model.SecretFile) error
}
