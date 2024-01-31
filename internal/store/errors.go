package store

import "errors"

var (
	ErrIncorrectPassword = errors.New("incorrect password")
	ErrNotAuthenticated  = errors.New("user is not authenticated")
	ErrUserAlredyExist   = errors.New("user already exist")
	ErrRecordNotFound    = errors.New("record not found")
	ErrUnableToGetRows   = errors.New("unable to get rows")
)
