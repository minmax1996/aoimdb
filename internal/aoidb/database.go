package aoidb

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
