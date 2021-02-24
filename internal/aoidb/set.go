package aoidb

//Set hset
type Set struct {
	cache map[string]interface{}
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
		cache: make(map[string]interface{}),
	}
}

// Set : set key v1
func (hs *Set) Set(key string, value interface{}) error {
	return nil
}

// Get : get key
func (hs *Set) Get(key string) (interface{}, error) {
	return nil, nil
}
