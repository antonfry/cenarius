package mcache

import (
	"cenarius/internal/model"
	"sync"
)

type MCacheRepo struct {
	store *MCache
	cache *model.SecretCache
	mutex sync.RWMutex
}

func (r *MCacheRepo) Save(c *model.SecretCache) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.cache = c
	return nil
}

func (r *MCacheRepo) Get() (*model.SecretCache, error) {
	d := r.cache
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return d, nil
}

func (r *MCacheRepo) Close() error {
	return nil
}
