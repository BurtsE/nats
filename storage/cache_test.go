package storage

import (
	"main/service"
	"strconv"
	"testing"
)

func TestCreate(t *testing.T) {
	m := GetCache()
	for i := 0; i < 5; i++ {
		key := strconv.Itoa(i)
		m.storage[key] = service.Message{Order_uid: key}
	}

	m = GetCache()
	for i := 0; i < 5; i++ {
		key := strconv.Itoa(i)
		val, ok := m.storage[key]
		if !ok || val.Order_uid != key {
			t.Errorf("cache changed or lost")
		}
	}
}

func TestSet(t *testing.T) {
	result := make(map[string]service.Message)
	c := GetCache()
	for i := 0; i < 5; i++ {
		key := strconv.Itoa(i)
		result[key] = service.Message{Order_uid: key}
	}
	for key, value := range result {
		c.Set(key, value)
	}
	for key, value := range result {
		err := c.Set(key, value)
		if err == nil {
			t.Errorf("existing value was reset")
		}
	}
	for i := 0; i < 5; i++ {
		key := strconv.Itoa(i)
		val, ok := c.storage[key]
		if !ok || val.Order_uid != result[key].Order_uid {
			t.Errorf("cache changed or lost")
		}
	}

}

func TestGet(t *testing.T) {
	result := make(map[string]service.Message)
	c := GetCache()
	for i := 0; i < 5; i++ {
		key := strconv.Itoa(i)
		result[key] = service.Message{Order_uid: key}
		c.storage[key] = service.Message{Order_uid: key}
	}
	for i := 0; i < 5; i++ {
		key := strconv.Itoa(i)
		val, err := c.Get(key)
		if err != nil || val.Order_uid != result[key].Order_uid {
			t.Errorf("cache changed or lost")
		}
	}
}
