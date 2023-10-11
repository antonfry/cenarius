package store

import (
	"cenarius/internal/model"
	"context"
)

type UserRepository interface {
	FindByID(context.Context, int) (*model.User, error)
	FindByLogin(context.Context, string) (*model.User, error)
	Create(context.Context, *model.User) error
}

type LoginWithPasswordRepository interface {
	GetByName(context.Context, string) (*model.LoginWithPassword, error)
	GetByID(context.Context, int) (*model.LoginWithPassword, error)
	Add(context.Context, *model.LoginWithPassword) error
	Delete(context.Context, *model.LoginWithPassword) error
}

type CreditCardRepository interface {
	GetByName(context.Context, string) (*model.CreditCard, error)
	GetByID(context.Context, int) (*model.CreditCard, error)
	Add(context.Context, *model.CreditCard) error
	Delete(context.Context, *model.CreditCard) error
}

type SecretTextRepository interface {
	GetByName(context.Context, string) (*model.SecretText, error)
	GetByID(context.Context, int) (*model.SecretText, error)
	Add(context.Context, *model.SecretText) error
	Delete(context.Context, *model.SecretText) error
}

type SecretBinaryRepository interface {
	GetByName(context.Context, string) (*model.SecretBinary, error)
	GetByID(context.Context, int) (*model.SecretBinary, error)
	Add(context.Context, *model.SecretBinary) error
	Delete(context.Context, *model.SecretBinary) error
}

type SecretDataRepository interface {
	GetByName(context.Context, string) (*model.SecretData, error)
	GetByID(context.Context, int) (*model.SecretData, error)
	Add(context.Context, *model.SecretData) error
	Delete(context.Context, *model.SecretData) error
}
