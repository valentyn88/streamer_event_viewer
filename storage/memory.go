package storage

import "sync"

// InMemory in memory storage.
type InMemory struct {
	Mux    sync.Mutex
	Values []interface{}
}

// NewInMemory creates new instance of InMemory.
func NewInMemory() *InMemory {
	return &InMemory{Values: make([]interface{}, 0)}
}

// Save saves value to storage by key.
func (im InMemory) Save(val interface{}) {
	im.Mux.Lock()
	im.Values = append(im.Values, val)
	im.Mux.Unlock()
}

// Get gets values by key.
func (im InMemory) Last(l int) []interface{} {
	var res []interface{}

	im.Mux.Lock()
	for i := len(im.Values); i > 0; i-- {
		if len(im.Values) == l {
			break
		}
		res = append(res, im.Values[i])
	}
	im.Mux.Unlock()

	return res
}
