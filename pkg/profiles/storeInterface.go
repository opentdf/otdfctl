package profiles

type NewStoreInterface func(namespace string, key string) StoreInterface

type StoreInterface interface {
	Exists() bool
	Get(value interface{}) error
	Set(value interface{}) error
	Delete() error
}
