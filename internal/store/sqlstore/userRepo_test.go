package sqlstore_test

import (
	"cenarius/internal/model"
	"cenarius/internal/store/sqlstore"
	"context"
	"reflect"
	"testing"
)

type userTests struct {
	name     string
	login    string
	password string
	want     *model.User
	wantErr  bool
}

var utests = []*userTests{
	{
		name:     "Valid",
		login:    "validlogin",
		password: "valid_password",
		want:     &model.User{Login: "validlogin", Password: ""},
		wantErr:  false,
	},
	{
		name:     "InValidLogin",
		login:    "invalid@login",
		password: "valid_password",
		want:     nil,
		wantErr:  true,
	},
	{
		name:     "validLoginWithInvalidPassword",
		login:    "validLoginWithInvalidPassword",
		password: "1",
		want:     nil,
		wantErr:  true,
	},
}

func TestUserRepository_Create(t *testing.T) {
	s, teardown := sqlstore.TestStore(t, databaseURL)
	defer teardown("users")

	for _, tt := range utests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.User().Create(context.Background(), &model.User{Login: tt.login, Password: tt.password})
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.FindByLogin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUserRepository_FindByLogin(t *testing.T) {
	s, teardown := sqlstore.TestStore(t, databaseURL)
	defer teardown("users")

	for _, tt := range utests {
		t.Run(tt.name, func(t *testing.T) {
			_ = s.User().Create(context.Background(), &model.User{Login: tt.login, Password: tt.password})
			u, err := s.User().FindByLogin(context.Background(), tt.login)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.FindByLogin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if u != nil {
				tt.want.ID = u.ID
				tt.want.EncryptedPassword = u.EncryptedPassword
			}
			if !reflect.DeepEqual(u, tt.want) {
				t.Errorf("UserRepository.FindByLogin() = %v, want %v", u, tt.want)
			}
		})
	}
}

func TestUserRepository_FindByID(t *testing.T) {
	s, teardown := sqlstore.TestStore(t, databaseURL)
	defer teardown("users")

	for _, tt := range utests {
		t.Run(tt.name, func(t *testing.T) {
			u := &model.User{Login: tt.login, Password: tt.password}
			_ = s.User().Create(context.Background(), u)
			u, err := s.User().FindByID(context.Background(), u.ID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.FindByLogin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if u != nil {
				tt.want.ID = u.ID
				tt.want.EncryptedPassword = u.EncryptedPassword
			}
			if !reflect.DeepEqual(u, tt.want) {
				t.Errorf("UserRepository.FindByLogin() = %v, want %v", u, tt.want)
			}
		})
	}
}
