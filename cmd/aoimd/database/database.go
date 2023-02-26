package database

import (
	"fmt"
	"regexp"

	"github.com/minmax1996/aoimdb/cmd/aoimd/hsetter"
	"github.com/minmax1996/aoimdb/cmd/aoimd/setter"
	"github.com/minmax1996/aoimdb/cmd/aoimd/table"
)

// Database Database structure
type Database struct {
	Name   string
	Sets   setter.Setter
	HSets  hsetter.HSetter
	Tables map[string]*table.Table
}

// NewDatabase database constructir
func NewDatabase(name string) *Database {
	return &Database{
		Name:   name,
		Sets:   setter.NewSet(),
		HSets:  hsetter.NewHSet(),
		Tables: make(map[string]*table.Table),
	}
}

func (db *Database) Keys(keysPattern string) []string {
	result := []string{}
	for _, k := range db.Sets.Keys() {
		if matched, _ := regexp.MatchString(keysPattern, k); matched {
			result = append(result, fmt.Sprintf("SET %s.%s", db.Name, k))
		}
	}
	for _, k := range db.HSets.Keys() {
		if matched, _ := regexp.MatchString(keysPattern, k); matched {
			result = append(result, fmt.Sprintf("HSET %s.%s", db.Name, k))
		}
	}

	for k := range db.Tables {
		if matched, _ := regexp.MatchString(keysPattern, k); matched {
			result = append(result, fmt.Sprintf("TABLE %s.%s", db.Name, k))
		}
	}
	return result
}

// Get Get
func (db *Database) Get(key string) (interface{}, error) {
	return db.Sets.Get(key)
}

// Set Set
func (db *Database) Set(key string, value interface{}) error {
	return db.Sets.Set(key, value)
}

// HSet HSet
func (db *Database) HSet(key string) (interface{}, error) {
	return db.Sets.Get(key)
}

// HGet HGet
func (db *Database) HGet(key string, value interface{}) error {
	return db.Sets.Set(key, value)
}
