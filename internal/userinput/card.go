package userinput

import (
	"cenarius/internal/model"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func InputCreditCard() *model.CreditCard {
	name := Input("Secret Name")
	inputNumber := Input("Number")
	_, err := strconv.ParseInt(inputNumber, 10, 64)
	if err != nil {
		log.Errorf("Wrong Number: %v", err)
		return nil
	}
	ownerName := Input("Owner Name")
	ownerLastName := Input("Owner Last Name")
	inputCvc := Input("CVC number")
	_, err = strconv.Atoi(inputCvc)
	if err != nil {
		log.Errorf("Wrong CVC: %v", err)
		return nil
	}
	meta := Input("Meta")
	m := &model.CreditCard{
		OwnerName:     ownerName,
		OwnerLastName: ownerLastName,
		Number:        inputNumber,
		CVC:           inputCvc,
	}
	m.Name = name
	m.Meta = meta
	return m
}
