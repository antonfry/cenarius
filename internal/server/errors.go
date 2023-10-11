package server

import "errors"

var (
	errNotAuthenticated = errors.New("user is not authenticated")
	errUserAlredyExist  = errors.New("user already exist")
)
