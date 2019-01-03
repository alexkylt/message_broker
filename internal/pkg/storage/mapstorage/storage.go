package mapstorage

import (
	"errors"
	"fmt"
)

type mapStorage struct {
	db map[string]string
}

func InitStorage() *mapStorage {
	return &mapStorage{db: make(map[string]string)}
}

func (s *mapStorage) Get(key string) (string, error) {
	v, ok := s.db[key]
	if !ok {
		errMsg := fmt.Sprintf("Key:%q is not found", key)
		return "", errors.New(errMsg)
	}
	return v, nil
}

func (s *mapStorage) Set(key, value string) error {
	s.db[key] = value
	return nil
}

func (s *mapStorage) Delete(key string) error {
	delete(s.db, key)
	return nil
}
