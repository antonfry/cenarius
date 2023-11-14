package userinput

import (
	"cenarius/internal/model"
	"errors"
	"os"

	log "github.com/sirupsen/logrus"
)

func InputSecretFile(askPath bool) *model.SecretFile {
	m := &model.SecretFile{}
	name := Input("Secret Name")
	if askPath {
		path := Input("File path")
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			log.Errorf("Path doesn't exist: %s", path)
			return nil
		}
		m.Path = path
	}
	meta := Input("Meta")
	m.Name = name
	m.Meta = meta
	return m
}
