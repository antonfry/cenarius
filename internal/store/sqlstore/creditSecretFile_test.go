package sqlstore_test

import (
	"cenarius/internal/model"
	"cenarius/internal/store/sqlstore"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type SFtest struct {
	name    string
	m       *model.SecretFile
	wantErr bool
}

var sFTests = []*SFtest{
	{
		name:    "Valid",
		m:       &model.SecretFile{Path: "/some/path"},
		wantErr: false,
	},
	{
		name:    "PathWithSpace",
		m:       &model.SecretFile{Path: "wrong path"},
		wantErr: true,
	},
	{
		name:    "InValidsymbol",
		m:       &model.SecretFile{Path: "Ж"},
		wantErr: true,
	},
}

func TestSecretFileRepository_Add(t *testing.T) {
	s, teardown := sqlstore.TestStore(t, databaseURL)
	defer teardown("SecretFile")

	for _, tt := range sFTests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.SecretFile().Add(context.Background(), tt.m); (err != nil) != tt.wantErr {
				t.Errorf("SecretFileRepository.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSecretFileRepository_Update(t *testing.T) {
	s, teardown := sqlstore.TestStore(t, databaseURL)
	defer teardown("SecretFile")

	for _, tt := range sFTests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.SecretFile().Update(context.Background(), tt.m); (err != nil) != tt.wantErr {
				t.Errorf("SecretFileRepository.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSecretFileRepository_Delete(t *testing.T) {
	s, teardown := sqlstore.TestStore(t, databaseURL)
	defer teardown("SecretFile")
	for _, tt := range sFTests {
		t.Run(tt.name, func(t *testing.T) {
			_ = s.SecretFile().Add(context.Background(), tt.m)
			if err := s.SecretFile().Delete(context.Background(), tt.m); err != nil {
				t.Errorf("SecretFileRepository.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSecretFileRepository_SearchByName(t *testing.T) {
	s, teardown := sqlstore.TestStore(t, databaseURL)
	defer teardown("SecretFile")
	for _, tt := range sFTests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Name = tt.name
			_ = s.SecretFile().Add(context.Background(), tt.m)
			l, err := s.SecretFile().SearchByName(context.Background(), tt.name, 0)
			if err != nil {
				t.Errorf("SecretFileRepository.SearchByName() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				assert.NotEmpty(t, l)
			} else {
				assert.Empty(t, l)
			}
		})
	}
}

func TestSecretFileRepository_GetByID(t *testing.T) {
	s, teardown := sqlstore.TestStore(t, databaseURL)
	defer teardown("SecretFile")

	for _, tt := range sFTests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Name = tt.name
			_ = s.SecretFile().Add(context.Background(), tt.m)
			lp, err := s.SecretFile().GetByID(context.Background(), tt.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("SecretFileRepository.GetByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				assert.NotEmpty(t, lp)
			} else {
				assert.Nil(t, lp)
			}
		})
	}
}
