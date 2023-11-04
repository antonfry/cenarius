package memstore_test

import (
	"cenarius/internal/model"
	"cenarius/internal/store/memstore"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreditCardRepository_Add(t *testing.T) {
	s := memstore.New()
	m := &model.CreditCard{
		OwnerName:     "Test",
		OwnerLastName: "Last Test",
		Number:        "335353",
		CVC:           "343",
	}
	m.ID = 1
	m.UserID = 1
	m.Name = "TestName"
	m.Meta = "TestMeta"
	err := s.CreditCardRepository.Add(context.Background(), m)
	assert.NoError(t, err)
}
