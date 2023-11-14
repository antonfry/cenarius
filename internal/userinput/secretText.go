package userinput

import (
	"cenarius/internal/model"
)

func InputSecretText() *model.SecretText {
	name := Input("Secret Name")
	text := Input("Secret Text")
	meta := Input("Meta")
	m := &model.SecretText{Text: text}
	m.Name = name
	m.Meta = meta
	return m
}
