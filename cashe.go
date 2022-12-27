package main

import (
	"errors"
	"log"
	"sync"
)

type Cache struct {
	storage map[string]message
	mutex   sync.RWMutex
}

func NewCashe() *Cache {
	return &Cache{
		storage: make(map[string]message),
		mutex:   sync.RWMutex{},
	}
}

func (c *Cache) Set(key string, value message) error {
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
func (c *Cache) Get(key string) (message, error) {
	log.Println("getting message from cache")
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	val, ok := c.storage[key]
	if !ok {
		return message{}, errors.New("message not found")
	}
	return val, nil
}

func (c *Cache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.storage, key)
}
