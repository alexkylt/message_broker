package storage

type StorageInterface interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Delete(key string) error
	Keys(pattern string) ([]string, error)
}
