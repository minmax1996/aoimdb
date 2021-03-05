package aoimdb

//Database Database structure
type Database struct {
	name  string
	sets  *Set
	hsets *HSet
}

// NewDatabase database constructir
func NewDatabase(name string) *Database {
	return &Database{
		name:  name,
		sets:  NewSet(),
		hsets: NewHSet(),
	}
}

//Get Get
func (db *Database) Get(key string) (interface{}, error) {
	return db.sets.Get(key)
}

//Set Set
func (db *Database) Set(key string, value interface{}) error {
	return db.sets.Set(key, value)
}

//HSet HSet
func (db *Database) HSet(key string) (interface{}, error) {
	return db.sets.Get(key)
}

//HGet HGet
func (db *Database) HGet(key string, value interface{}) error {
	return db.sets.Set(key, value)
}
