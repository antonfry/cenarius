package mcache

import (
	"cenarius/internal/cache"
	"cenarius/internal/model"
)

type MCache struct {
	MCacheRepo *MCacheRepo
}

func New() *MCache {
	return &MCache{}
}

func (c *MCache) Cache() cache.Casher {
	if c.MCacheRepo == nil {
		c.MCacheRepo = &MCacheRepo{
			store: c,
			cache: &model.SecretCache{},
		}
	}
	return c.MCacheRepo
}
