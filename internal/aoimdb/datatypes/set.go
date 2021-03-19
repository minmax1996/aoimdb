package datatypes

import "errors"

//Set hset
type Set struct {
	Cache map[string]interface{}
}

// Setter hsetter interface
type Setter interface {
	// set whole structure in interface{}
	Set(string, interface{}) error
	Get(string) (interface{}, error)
}

// NewSet sets constructor
func NewSet() *Set {
	return &Set{
		Cache: make(map[string]interface{}),
	}
}

// Set : set key v1
func (hs *Set) Set(key string, value interface{}) error {
	hs.Cache[key] = value
	return nil
}

// Get : get key
func (hs *Set) Get(key string) (interface{}, error) {
	val, ok := hs.Cache[key]
	if !ok {
		return nil, errors.New("element with this name not found")
	}
	return val, nil
}
