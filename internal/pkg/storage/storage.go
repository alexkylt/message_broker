package storage
// StrgInterface ... interface
type StrgInterface interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Delete(key string) error
	Keys(pattern string) ([]string, error)
}
