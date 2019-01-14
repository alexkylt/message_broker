package mapstorage

import (
	"errors"
	"fmt"
	"strings"
)
// MapStorage ...
type MapStorage struct {
	db map[string]string
}



// InitStorage ...
func InitStorage() *MapStorage {
	return &MapStorage{db: make(map[string]string)}
}
// Get ...
func (s *MapStorage) Get(key string) (string, error) {
	v, ok := s.db[key]
	if !ok {
		errMsg := fmt.Sprintf("Key:%q is not found", key)
		return "", errors.New(errMsg)
	}
	return v, nil
}
// Set ...
func (s *MapStorage) Set(key, value string) error {
	s.db[key] = value
	return nil
}
// Delete ...
func (s *MapStorage) Delete(key string) error {
	delete(s.db, key)
	return nil
}
// Keys ...
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

