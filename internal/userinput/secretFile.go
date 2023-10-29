package userinput

import (
	"cenarius/internal/model"
	"errors"
	"log"
	"os"
)

func InputSecretFile(askPath bool) *model.SecretFile {
	m := &model.SecretFile{}
	name := Input("Secret Name")
	if askPath {
		path := Input("File path")
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			log.Fatalf("Path doesn't exist: %s", path)
		}
		m.Path = path
	}
	meta := Input("Meta")
	m.Name = name
	m.Meta = meta
	return m
}
