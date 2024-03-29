package hsetter

// HSetter hsetter interface
type HSetter interface {
	AddOrUpdate(string, ...interface{}) error
	// add key field1 value1 field2 value2
	Add(string, ...interface{}) error
	// set whole structure in map[string]interface
	Set(string, interface{}) error
	Keys() []string
}

//HSet hset
type HSet struct {
	Cache map[string]map[string]interface{}
}

// NewHSet sets constructor
func NewHSet() *HSet {
	return &HSet{
		Cache: make(map[string]map[string]interface{}),
	}
}

// AddOrUpdate : addorupdate key f1 v1
func (hs *HSet) AddOrUpdate(key string, fieldValues ...interface{}) error {
	return nil
}

// Add : add key f1 v1 f2 v2
func (hs *HSet) Add(key string, fieldValues ...interface{}) error {
	return nil
}

// Set : set key somestruct
func (hs *HSet) Set(key string, value interface{}) error {
	return nil
}

func (hs *HSet) Keys() []string {
	result := make([]string, 0, len(hs.Cache))
	for k := range hs.Cache {
		result = append(result, k)
	}
	return result
}
