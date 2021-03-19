package datatypes

//Database Database structure
type Database struct {
	Name  string
	Sets  *Set
	HSets *HSet
}

// NewDatabase database constructir
func NewDatabase(name string) *Database {
	return &Database{
		Name:  name,
		Sets:  NewSet(),
		HSets: NewHSet(),
	}
}

//Get Get
func (db *Database) Get(key string) (interface{}, error) {
	return db.Sets.Get(key)
}

//Set Set
func (db *Database) Set(key string, value interface{}) error {
	return db.Sets.Set(key, value)
}

//HSet HSet
func (db *Database) HSet(key string) (interface{}, error) {
	return db.Sets.Get(key)
}

//HGet HGet
func (db *Database) HGet(key string, value interface{}) error {
	return db.Sets.Set(key, value)
}
