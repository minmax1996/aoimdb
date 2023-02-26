package database

import (
	"crypto/md5"
	"encoding/hex"
	"errors"

	"github.com/minmax1996/aoimdb/cmd/aoimd/datatypes"
	"github.com/minmax1996/aoimdb/cmd/aoimd/setter"
	"github.com/minmax1996/aoimdb/cmd/aoimd/table"
	"github.com/minmax1996/aoimdb/internal/pkg/filestorage"
	"github.com/minmax1996/aoimdb/internal/pkg/logger"
)

// databaseInstance to easy access to database from any part of program
var databaseInstance *DatabaseController

// DatabaseController Database structure
type DatabaseController struct {
	Databases    map[string]*Database
	Users        setter.Setter
	AccessTokens setter.Setter
}

// NewDatabaseController database constructir
func NewDatabaseController() *DatabaseController {
	dbc := &DatabaseController{
		Databases:    make(map[string]*Database),
		Users:        setter.NewSet(),
		AccessTokens: setter.NewSet(),
	}

	return dbc
}

// InitDatabaseController uses to restore data from backup if exists or create new database
func InitDatabaseController() {
	databaseInstance = NewDatabaseController()
	if err := filestorage.RestoreFromBackup(databaseInstance); err != nil {
		logger.Error("cant restore " + err.Error())
	}

	// backups only exported fields
	filestorage.StartBackups(databaseInstance, 30)
}

// AuthentificateByUserPass v
func AuthentificateByUserPass(user, pass string) error {
	val, err := databaseInstance.Users.Get(user)
	if err != nil {
		return err
	}

	if val.(string) != pass {
		return errors.New("pass not equals")
	}

	return nil
}

// AuthentificateByToken
func AuthentificateByToken(token string) error {
	_, err := databaseInstance.AccessTokens.Get(token)
	return err
}

// AddUser AddUser to auth with
func AddUser(user, pass string) error {
	if err := databaseInstance.Users.Set(user, pass); err != nil {
		return err
	}
	hasher := md5.New()
	hasher.Write([]byte(user))
	hasher.Write([]byte(pass))
	accessToken := hex.EncodeToString(hasher.Sum(nil))
	logger.Info(accessToken)
	return databaseInstance.AccessTokens.Set(accessToken, user)
}

// SelectDatabase checks if database exists and create it if not
func SelectDatabase(name string) {
	if _, ok := databaseInstance.Databases[name]; !ok {
		logger.Info("create database")
		databaseInstance.Databases[name] = NewDatabase(name)
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

// Set sets data to set by key
func GetKeys(dbName, keyspattern string) ([]string, error) {
	result := []string{}
	if len(dbName) > 0 {
		db, ok := databaseInstance.Databases[dbName]
		if !ok {
			return nil, errors.New("database with this name does not exist")
		}
		result = append(result, db.Keys(keyspattern)...)
	} else {
		for _, db := range databaseInstance.Databases {
			result = append(result, db.Keys(keyspattern)...)
		}
	}

	return result, nil
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

func CreateTable(dbName, tableName string, columNames []string, columTypes []datatypes.Datatype) error {
	db, ok := databaseInstance.Databases[dbName]
	if !ok {
		return errors.New("database with this name does not exist")
	}
	_, ok = db.Tables[tableName]
	if ok {
		return errors.New("cant create new table, table already exists")
	}

	if len(columNames) != len(columTypes) {
		return errors.New("cant create new table, number params not equal")
	}
	db.Tables[tableName] = table.NewTableSchema(tableName, columNames, columTypes)
	return nil
}

func InsertIntoTable(dbName, tableName string, columNames []string, values []interface{}) error {
	table, err := GetTable(dbName, tableName)
	if err != nil {
		return err
	}

	return table.Insert(columNames, values)
}

func SelectFromTable(dbName, tableName string, columNames []string) *table.Table {
	table, err := GetTable(dbName, tableName)
	if err != nil {
		return nil
	}
	if len(columNames) == 0 {
		columNames = table.ColumnNames
	}
	return table.Select(columNames)
}

func GetTable(dbName, tableName string) (*table.Table, error) {
	db, ok := databaseInstance.Databases[dbName]
	if !ok {
		return nil, errors.New("database with this name does not exist")
	}
	table, ok := db.Tables[tableName]
	if !ok {
		return nil, errors.New("table with this name does not exist")
	}

	return table, nil
}
