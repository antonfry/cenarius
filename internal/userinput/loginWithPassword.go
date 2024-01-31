package userinput

import (
	"cenarius/internal/model"
)

func InputLoginWithPassword() *model.LoginWithPassword {
	name := Input("Secret Name")
	login := Input("Login")
	password := Input("Password")
	meta := Input("Meta")
	m := &model.LoginWithPassword{Login: login, Password: password}
	m.Name = name
	m.Password = password
	m.Meta = meta
	return m
}
