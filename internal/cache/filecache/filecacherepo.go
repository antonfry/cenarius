package filecache

import (
	"bufio"
	"cenarius/internal/model"
	"encoding/json"
	"sync"

	log "github.com/sirupsen/logrus"
)

type FileCacheRepo struct {
	store *FileCache
	mutex sync.RWMutex
}

func (r *FileCacheRepo) Save(c *model.SecretCache) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if err := r.store.file.Truncate(0); err != nil {
		return err
	}
	if _, err := r.store.file.Seek(0, 0); err != nil {
		return err
	}
	jsonSting, err := json.Marshal(c)
	if err != nil {
		log.Errorf("Save: failed to marshal metrics %v", c)
		return err
	}
	_, err = r.store.file.Write([]byte(jsonSting))
	if err != nil {
		log.Errorf("Save: failed to save metric %v", err)
		return err
	}
	_, err = r.store.file.WriteString("\n")
	if err != nil {
		log.Errorf("Save: failed to save new line")
		return err
	}

	return nil
}

func (r *FileCacheRepo) Get() (*model.SecretCache, error) {
	_, _ = r.store.file.Seek(0, 0)
	var d *model.SecretCache
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	s := bufio.NewScanner(r.store.file)
	for s.Scan() {
		log.Debug("FileCacheRepo.Get Scan: ", s.Text())
		if err := json.Unmarshal([]byte(s.Text()), &d); err != nil {
			log.Printf("GetAll failed to Unmarshal json: %+v", s.Text())
			return nil, err
		}
	}
	if err := s.Err(); err != nil {
		log.Errorf("GetAll: Failed to Scan file")
		return nil, err
	}
	log.Debug("FileCacheRepo.Get returning: ", d)
	return d, nil
}

func (r *FileCacheRepo) Close() error {
	if err := r.store.file.Sync(); err != nil {
		return err
	}
	if err := r.store.file.Close(); err != nil {
		return err
	}
	return nil
}
