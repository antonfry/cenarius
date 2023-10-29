package sqlstore_test

import (
	"cenarius/internal/model"
	"cenarius/internal/store/sqlstore"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginWithPasswordRepository_Add(t *testing.T) {
	s, teardown := sqlstore.TestStore(t, databaseURL)
	defer teardown("loginwithpassword")
	tests := []struct {
		name    string
		m       *model.LoginWithPassword
		wantErr bool
	}{
		{
			name:    "Valid",
			m:       &model.LoginWithPassword{Login: "somevalidlogin", Password: "somevalidPassword"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.LoginWithPassword().Add(context.Background(), tt.m); (err != nil) != tt.wantErr {
				t.Errorf("LoginWithPasswordRepository.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoginWithPasswordRepository_Update(t *testing.T) {
	s, teardown := sqlstore.TestStore(t, databaseURL)
	defer teardown("loginwithpassword")
	tests := []struct {
		name    string
		m       *model.LoginWithPassword
		wantErr bool
	}{
		{
			name:    "Valid",
			m:       &model.LoginWithPassword{Login: "somevalidlogin", Password: "somevalidPassword"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.LoginWithPassword().Update(context.Background(), tt.m); (err != nil) != tt.wantErr {
				t.Errorf("LoginWithPasswordRepository.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoginWithPasswordRepository_Delete(t *testing.T) {
	s, teardown := sqlstore.TestStore(t, databaseURL)
	defer teardown("loginwithpassword")
	tests := []struct {
		name    string
		m       *model.LoginWithPassword
		wantErr bool
	}{
		{
			name:    "Valid",
			m:       &model.LoginWithPassword{Login: "somevalidlogin", Password: "somevalidPassword"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = s.LoginWithPassword().Add(context.Background(), tt.m)
			if err := s.LoginWithPassword().Delete(context.Background(), tt.m); (err != nil) != tt.wantErr {
				t.Errorf("LoginWithPasswordRepository.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoginWithPasswordRepository_SearchByName(t *testing.T) {
	s, teardown := sqlstore.TestStore(t, databaseURL)
	defer teardown("loginwithpassword")
	tests := []struct {
		name    string
		m       *model.LoginWithPassword
		wantErr bool
	}{
		{
			name:    "Valid",
			m:       &model.LoginWithPassword{Login: "somevalidlogin", Password: "somevalidPassword"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Name = tt.name
			_ = s.LoginWithPassword().Add(context.Background(), tt.m)
			l, err := s.LoginWithPassword().SearchByName(context.Background(), tt.name, 0)
			if err != nil {
				t.Errorf("LoginWithPasswordRepository.SearchByName() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.NotEmpty(t, l)
		})
	}
}

func TestLoginWithPasswordRepository_GetByID(t *testing.T) {
	s, teardown := sqlstore.TestStore(t, databaseURL)
	defer teardown("loginwithpassword")

	tests := []struct {
		name    string
		m       *model.LoginWithPassword
		wantErr bool
	}{
		{
			name:    "Valid",
			m:       &model.LoginWithPassword{Login: "somevalidlogin", Password: "somevalidPassword"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Name = tt.name
			_ = s.LoginWithPassword().Add(context.Background(), tt.m)
			lp, err := s.LoginWithPassword().GetByID(context.Background(), tt.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoginWithPasswordRepository.GetByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				assert.NotEmpty(t, lp)
			} else {
				assert.Nil(t, lp)
			}
		})
	}
}
