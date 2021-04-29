package datatypes

type TableIndex interface {
	ColName() string
	AddToIndex(*Row) bool
	RemoveFromIndex(*Row) bool
	Find(func(value interface{}) bool) ([]Row, error)
	FindByValue(value interface{}) ([]Row, error)
}

type MapUnicIndex struct {
	colName  string
	colIndex int
	cache    map[interface{}]*Row
}

func NewMapUnicIndex(colIdx int, colName string) TableIndex {
	return MapUnicIndex{
		colIndex: colIdx,
		colName:  colName,
		cache:    make(map[interface{}]*Row),
	}
}

func (m MapUnicIndex) ColName() string {
	return m.colName
}

func (m MapUnicIndex) AddToIndex(r *Row) bool {
	m.cache[(*r)[m.colIndex]] = r
	return true
}

func (m MapUnicIndex) RemoveFromIndex(r *Row) bool {
	delete(m.cache, (*r)[m.colIndex])
	return true
}

func (m MapUnicIndex) Find(filterfunc func(value interface{}) bool) ([]Row, error) {
	result := make([]Row, 0)
	for k, v := range m.cache {
		if filterfunc(k) {
			result = append(result, *v)
		}
	}
	return result, nil
}

func (m MapUnicIndex) FindByValue(value interface{}) ([]Row, error) {
	result := make([]Row, 0)
	if v, ok := m.cache[value]; ok {
		result = append(result, *v)
	}
	return result, nil
}
