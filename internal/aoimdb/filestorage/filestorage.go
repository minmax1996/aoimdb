package filestorage

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/minmax1996/aoimdb/logger"
)

var databaseBackupFileName = "dbbackup.aoimdb"

//RestoreFromBackup reads file to db interface
func RestoreFromBackup(db interface{}) error {
	data, err := ReadFromFile(databaseBackupFileName)
	if err != nil {
		return err
	}
	return DecompressAndDecore(db, data)
}

//StartBackups start goroutine to backup db interface on interval
func StartBackups(db interface{}, seconds int) {
	go func() {
		for {
			time.Sleep(time.Duration(seconds) * time.Second)
			if err := Backup(db); err != nil {
				logger.Error(err)
			}
		}
	}()
}

//Backup writes to file db interface
func Backup(db interface{}) error {
	zipped, err := EncodeAndCompress(db)
	if err != nil {
		return err
	}
	WriteToFile(zipped, databaseBackupFileName)
	logger.Debug("Db Backup executed")
	return nil
}

//ReadFromFile reads from file
func ReadFromFile(from string) ([]byte, error) {
	f, err := os.Open(from)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return data, nil
}

//WriteToFile writes to file
func WriteToFile(s []byte, file string) {
	f, err := os.Create(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	f.Write(s)
}

//EncodeAndCompress ecode p interface tand then gzip it
func EncodeAndCompress(p interface{}) ([]byte, error) {
	buf := bytes.Buffer{}

	if err := gob.NewEncoder(&buf).Encode(p); err != nil {
		return nil, err
	}
	zipbuf := bytes.Buffer{}
	zipped := gzip.NewWriter(&zipbuf)
	zipped.Write(buf.Bytes())
	zipped.Close()
	return zipbuf.Bytes(), nil
}

//DecompressAndDecore decompresses []byte from s and decode to db interface{}
func DecompressAndDecore(db interface{}, s []byte) error {
	rdr, err := gzip.NewReader(bytes.NewReader(s))
	if err != nil {
		return err
	}
	defer rdr.Close()
	data, err := ioutil.ReadAll(rdr)
	if err != nil {
		return err
	}

	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(db); err != nil {
		return err
	}

	return nil
}
