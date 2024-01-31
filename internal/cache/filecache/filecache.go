package filecache

import (
	"cenarius/internal/cache"
	"os"
)

type FileCache struct {
	file          *os.File
	FileCacheRepo *FileCacheRepo
}

func New(file *os.File) *FileCache {
	return &FileCache{
		file: file,
	}
}

func (c *FileCache) Cache() cache.Casher {
	if c.FileCacheRepo == nil {
		c.FileCacheRepo = &FileCacheRepo{
			store: c,
		}
	}
	return c.FileCacheRepo
}
