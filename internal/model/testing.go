package model

import "testing"

func TestUser(t *testing.T) *User {
	return &User{
		Login:    "user@example.org",
		Password: "password",
	}
}
