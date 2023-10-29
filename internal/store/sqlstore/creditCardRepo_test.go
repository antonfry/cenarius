package sqlstore_test

import (
	"cenarius/internal/model"
	"cenarius/internal/store/sqlstore"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type CCtest struct {
	name    string
	m       *model.CreditCard
	wantErr bool
}

var ccTests = []*CCtest{
	{
		name: "Valid",
		m: &model.CreditCard{
			OwnerName:     "someOwner",
			OwnerLastName: "someOwnerLastName",
			Number:        "1111111111111111",
			CVC:           "111",
		},
		wantErr: false,
	},
}

func TestCreditCardRepository_Add(t *testing.T) {
	s, teardown := sqlstore.TestStore(t, databaseURL)
	defer teardown("CreditCard")

	for _, tt := range ccTests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.CreditCard().Add(context.Background(), tt.m); err != nil {
				t.Errorf("CreditCardRepository.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreditCardRepository_Update(t *testing.T) {
	s, teardown := sqlstore.TestStore(t, databaseURL)
	defer teardown("CreditCard")

	for _, tt := range ccTests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.CreditCard().Update(context.Background(), tt.m); err != nil {
				t.Errorf("CreditCardRepository.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreditCardRepository_Delete(t *testing.T) {
	s, teardown := sqlstore.TestStore(t, databaseURL)
	defer teardown("CreditCard")
	for _, tt := range ccTests {
		t.Run(tt.name, func(t *testing.T) {
			_ = s.CreditCard().Add(context.Background(), tt.m)
			if err := s.CreditCard().Delete(context.Background(), tt.m); err != nil {
				t.Errorf("CreditCardRepository.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreditCardRepository_SearchByName(t *testing.T) {
	s, teardown := sqlstore.TestStore(t, databaseURL)
	defer teardown("CreditCard")
	for _, tt := range ccTests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Name = tt.name
			_ = s.CreditCard().Add(context.Background(), tt.m)
			l, err := s.CreditCard().SearchByName(context.Background(), tt.name, 0)
			if err != nil {
				t.Errorf("CreditCardRepository.SearchByName() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.NotEmpty(t, l)
		})
	}
}

func TestCreditCardRepository_GetByID(t *testing.T) {
	s, teardown := sqlstore.TestStore(t, databaseURL)
	defer teardown("CreditCard")

	for _, tt := range ccTests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Name = tt.name
			_ = s.CreditCard().Add(context.Background(), tt.m)
			lp, err := s.CreditCard().GetByID(context.Background(), tt.m)
			if err != nil {
				t.Errorf("CreditCardRepository.GetByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.NotNil(t, lp)
		})
	}
}
