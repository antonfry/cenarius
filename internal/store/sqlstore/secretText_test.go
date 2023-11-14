package sqlstore_test

import (
	"cenarius/internal/model"
	"cenarius/internal/store/sqlstore"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type STtest struct {
	name    string
	m       *model.SecretText
	wantErr bool
}

var stTests = []*STtest{
	{
		name:    "Valid",
		m:       &model.SecretText{Text: "somesecretText"},
		wantErr: false,
	},
}

func TestSecretTextRepository_Add(t *testing.T) {
	s, teardown := sqlstore.TestStore(t, databaseURL)
	defer teardown("SecretText")

	for _, tt := range stTests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.SecretText().Add(context.Background(), tt.m); (err != nil) != tt.wantErr {
				t.Errorf("SecretTextRepository.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSecretTextRepository_Update(t *testing.T) {
	s, teardown := sqlstore.TestStore(t, databaseURL)
	defer teardown("SecretText")

	for _, tt := range stTests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.SecretText().Update(context.Background(), tt.m); (err != nil) != tt.wantErr {
				t.Errorf("SecretTextRepository.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSecretTextRepository_Delete(t *testing.T) {
	s, teardown := sqlstore.TestStore(t, databaseURL)
	defer teardown("SecretText")
	for _, tt := range stTests {
		t.Run(tt.name, func(t *testing.T) {
			_ = s.SecretText().Add(context.Background(), tt.m)
			if err := s.SecretText().Delete(context.Background(), tt.m.ID, tt.m.UserID); err != nil {
				t.Errorf("SecretTextRepository.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSecretTextRepository_SearchByName(t *testing.T) {
	s, teardown := sqlstore.TestStore(t, databaseURL)
	defer teardown("SecretText")
	for _, tt := range stTests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Name = tt.name
			_ = s.SecretText().Add(context.Background(), tt.m)
			l, err := s.SecretText().SearchByName(context.Background(), tt.name, 0)
			if err != nil {
				t.Errorf("SecretTextRepository.SearchByName() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				assert.NotEmpty(t, l)
			} else {
				assert.Empty(t, l)
			}
		})
	}
}

func TestSecretTextRepository_GetByID(t *testing.T) {
	s, teardown := sqlstore.TestStore(t, databaseURL)
	defer teardown("SecretText")

	for _, tt := range stTests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Name = tt.name
			_ = s.SecretText().Add(context.Background(), tt.m)
			lp, err := s.SecretText().GetByID(context.Background(), tt.m.ID, tt.m.UserID)
			if (err != nil) != tt.wantErr {
				t.Errorf("SecretTextRepository.GetByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				assert.NotEmpty(t, lp)
			} else {
				assert.Nil(t, lp)
			}
		})
	}
}
