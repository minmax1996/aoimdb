package aoimdb

import (
	"errors"

	"github.com/minmax1996/aoimdb/internal/aoimdb/datatypes"
	"github.com/minmax1996/aoimdb/internal/aoimdb/filestorage"
	"github.com/minmax1996/aoimdb/logger"
)

//databaseInstance to easy access to database from any part of program
var databaseInstance *DatabaseController

// DatabaseController Database structure
type DatabaseController struct {
	Databases map[string]*datatypes.Database
	users     *datatypes.Set
}

// NewDatabaseController database constructir
func NewDatabaseController() *DatabaseController {
	dbc := &DatabaseController{
		Databases: make(map[string]*datatypes.Database),
		users:     datatypes.NewSet(),
	}

	return dbc
}

//InitDatabaseController uses to restore data from backup if exists or create new database
func InitDatabaseController() {
	databaseInstance = NewDatabaseController()
	if err := filestorage.RestoreFromBackup(databaseInstance); err != nil {
		logger.Error("cant restore " + err.Error())
	}

	// backups only exported fields
	filestorage.StartBackups(databaseInstance)
}

//AuthentificateByUserPass v
func AuthentificateByUserPass(user, pass string) error {
	val, err := databaseInstance.users.Get(user)
	if err != nil {
		return err
	}

	if val.(string) != pass {
		return errors.New("pass not equals")
	}

	return nil
}

//AddUser AddUser to auth with
func AddUser(user, pass string) error {
	return databaseInstance.users.Set(user, pass)
}

// SelectDatabase checks if database exists and create it if not
func SelectDatabase(name string) {
	if _, ok := databaseInstance.Databases[name]; !ok {
		logger.Info("create database")
		databaseInstance.Databases[name] = datatypes.NewDatabase(name)
	}
}

// Get gets data from set by key
func Get(dbName, key string) (interface{}, error) {
	db, ok := databaseInstance.Databases[dbName]
	if !ok {
		return nil, errors.New("database with this name does not exist")
	}

	return db.Get(key)
}

// Set sets data to set by key
func Set(dbName, key string, value interface{}) error {
	db, ok := databaseInstance.Databases[dbName]
	if !ok {
		return errors.New("database with this name does not exist")
	}
	return db.Set(key, value)
}

// HSet sets data in hashset <NOT IMPLEMENTED YET>
func HSet(dbName, key string, value interface{}) error {
	db, ok := databaseInstance.Databases[dbName]
	if !ok {
		return errors.New("database with this name does not exist")
	}
	return db.Set(key, value)
}

// HGet gets data from hashset <NOT IMPLEMENTED YET>
func HGet(dbName, key string) (interface{}, error) {
	db, ok := databaseInstance.Databases[dbName]
	if !ok {
		return nil, errors.New("database with this name does not exist")
	}
	return db.Get(key)
}
