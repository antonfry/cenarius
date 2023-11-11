package cache

import "cenarius/internal/model"

type StoreCache interface {
	Cache() Casher
}

type Casher interface {
	GetAll() (*model.SecretCache, error)
	Save(*model.SecretCache) error
	Close() error
}
