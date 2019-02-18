package storage

// Storager commom interface for all storages.
type Storager interface {
	Save(val interface{})
	Last(l int) []interface{}
}
