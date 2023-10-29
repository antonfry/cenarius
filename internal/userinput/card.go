package userinput

import (
	"cenarius/internal/model"
	"log"
	"strconv"
)

func InputCreditCard() *model.CreditCard {
	name := Input("Secret Name")
	inputNumber := Input("Number")
	_, err := strconv.ParseInt(inputNumber, 10, 64)
	if err != nil {
		log.Fatalf("Wrong Number: %v", err)
	}
	ownerName := Input("Owner Name")
	ownerLastName := Input("Owner Last Name")
	inputCvc := Input("CVC number")
	_, err = strconv.Atoi(inputCvc)
	if err != nil {
		log.Fatalf("Wrong CVC: %v", err)
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
