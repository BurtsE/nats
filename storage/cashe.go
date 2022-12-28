package storage

import (
	"errors"
	"log"
	"main/service"
	"sync"
)

type cache struct {
	storage map[string]service.Message
	mutex   sync.RWMutex
}

var memory cache

func GetCache() *cache {
	if memory.storage == nil {
		memory = cache{
			storage: make(map[string]service.Message),
			mutex:   sync.RWMutex{},
		}
	}
	return &memory
}

func (c *cache) Set(key string, value service.Message) error {
	log.Println("adding message to cache", key)
	c.mutex.Lock()
	defer c.mutex.Unlock()
	_, exists := c.storage[key]
	if exists {
		return errors.New("message exists in cache")
	}
	c.storage[key] = value
	return nil
}
func (c *cache) Get(key string) (service.Message, error) {
	log.Println("getting message from cache")
	c.mutex.RLock()

	defer c.mutex.RUnlock()
	val, ok := c.storage[key]
	if !ok {
		return service.Message{}, errors.New("message not found")
	}
	return val, nil
}

func (c *cache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.storage, key)
}
