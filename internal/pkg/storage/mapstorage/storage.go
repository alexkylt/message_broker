package mapstorage

import (
	"errors"
	"fmt"
	"strings"
)

// MapStorage ... - structure for a map storage
type MapStorage struct {
	db map[string]string
}

// InitStorage ... - initiliaze map storage
func InitStorage() *MapStorage {
	return &MapStorage{db: make(map[string]string)}
}

// Get ... - method used to gather KV pair from a storage
func (s *MapStorage) Get(key string) (string, error) {
	v, ok := s.db[key]
	if !ok {
		errMsg := fmt.Sprintf("Key:%q is not found", key)
		return errMsg, errors.New(errMsg)
	}
	return v, nil
}

// Set ... - method used to insert/update KV pair to storage
func (s *MapStorage) Set(key, value string) error {
	s.db[key] = value
	return nil
}

// Delete ... - method used to delete KV pair from a storage
func (s *MapStorage) Delete(key string) error {
	delete(s.db, key)
	return nil
}

// Keys ... - method used to gather all KV pairs from a storage by pattern
func (s *MapStorage) Keys(key string) ([]string, error) {
	var keys []string
	var values []string

	for k, v := range s.db {
		//if k, found := s.db[key]; found {

		if strings.Contains(k, key) {
			fmt.Println(v)
			keys = append(keys, k)
			values = append(values, v)
		}
	}

	fmt.Println("KEYS - ", keys)
	fmt.Println("KEY - ", key)
	return keys, nil
}
