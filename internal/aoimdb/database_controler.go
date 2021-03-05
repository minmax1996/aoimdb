package aoimdb

import (
	"errors"

	"github.com/minmax1996/aoimdb/logger"
)

// DatabaseController Database structure
type DatabaseController struct {
	databases map[string]*Database
	users     *Set
}

// NewDatabaseController database constructir
func NewDatabaseController() *DatabaseController {
	return &DatabaseController{
		databases: make(map[string]*Database),
		users:     NewSet(),
	}
}

//AuthentificateByUserPass v
func (dbc *DatabaseController) AuthentificateByUserPass(user, pass string) error {
	val, err := dbc.users.Get(user)
	if err != nil {
		return err
	}

	if val.(string) != pass {
		return errors.New("pass not equals")
	}

	return nil
}

//AddUser AddUser
func (dbc *DatabaseController) AddUser(user, pass string) error {
	return dbc.users.Set(user, pass)
}

// SelectDatabase SelectDatabase
func (dbc *DatabaseController) SelectDatabase(name string) {
	if _, ok := dbc.databases[name]; !ok {
		logger.Info("create database")
		dbc.databases[name] = NewDatabase(name)
	}
}

// Get Get
func (dbc *DatabaseController) Get(dbName, key string) (interface{}, error) {
	db, ok := dbc.databases[dbName]
	if !ok {
		return nil, errors.New("database with this name does not exist")
	}

	return db.Get(key)
}

// Set Set
func (dbc *DatabaseController) Set(dbName, key string, value interface{}) error {
	db, ok := dbc.databases[dbName]
	if !ok {
		return errors.New("database with this name does not exist")
	}
	return db.Set(key, value)
}

// HSet hset
func (dbc *DatabaseController) HSet(dbName, key string) (interface{}, error) {
	db, ok := dbc.databases[dbName]
	if !ok {
		return nil, errors.New("database with this name does not exist")
	}
	return db.Get(key)
}

// HGet HGet
func (dbc *DatabaseController) HGet(dbName, key string, value interface{}) error {
	db, ok := dbc.databases[dbName]
	if !ok {
		return errors.New("database with this name does not exist")
	}
	return db.Set(key, value)
}
