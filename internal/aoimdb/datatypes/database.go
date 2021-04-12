package datatypes

import (
	"fmt"
	"regexp"
)

//Database Database structure
type Database struct {
	Name   string
	Sets   *Set
	HSets  *HSet
	Tables map[string]*Table
}

// NewDatabase database constructir
func NewDatabase(name string) *Database {
	return &Database{
		Name:   name,
		Sets:   NewSet(),
		HSets:  NewHSet(),
		Tables: make(map[string]*Table),
	}
}

func (db *Database) Keys(keysPattern string) []string {
	result := []string{}
	for k := range db.Sets.Cache {
		if matched, _ := regexp.MatchString(keysPattern, k); matched {
			result = append(result, fmt.Sprintf("SET %s.%s", db.Name, k))
		}
	}
	for k := range db.HSets.Cache {
		if matched, _ := regexp.MatchString(keysPattern, k); matched {
			result = append(result, fmt.Sprintf("HSET %s.%s", db.Name, k))
		}
	}
	return result
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
