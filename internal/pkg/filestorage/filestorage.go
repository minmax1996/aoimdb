package filestorage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/minmax1996/aoimdb/internal/pkg/logger"
)

var databaseBackupFileName = "dbbackup.aoimdb"

// Backup writes to file db interface
func Backup(db interface{}) error {
	zipped, err := encodeAndCompress(db)
	if err != nil {
		return err
	}
	writeToFile(zipped, databaseBackupFileName)
	logger.Debug("Db Backup executed")
	return nil
}

// RestoreFromBackup reads file to db interface
func RestoreFromBackup(db interface{}) error {
	data, err := readFromFile(databaseBackupFileName)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return decompressAndDecore(&db, data)
}

// StartBackups start goroutine to backup db interface on interval
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

// readFromFile reads from file
func readFromFile(from string) ([]byte, error) {
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

// writeToFile writes to file
func writeToFile(data []byte, file string) {
	f, err := os.Create(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	_, _ = f.Write(data)
}

// encodeAndCompress encode p interface tand then gzip it
func encodeAndCompress(p interface{}) ([]byte, error) {
	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(p); err != nil {
		return nil, err
	}
	//TODO fix later
	// zipbuf := bytes.Buffer{}
	// zipped := gzip.NewWriter(&zipbuf)
	// defer zipped.Close()
	// if _, err := zipped.Write(buf.Bytes()); err != nil {
	// 	return nil, err
	// }
	return buf.Bytes(), nil
}

// decompressAndDecore decompresses []byte from s and decode to db interface{}
func decompressAndDecore(db interface{}, s []byte) error {
	//TODO fix later
	// rdr, err := gzip.NewReader(bytes.NewReader(s))
	// if err != nil {
	// 	return err
	// }
	// defer rdr.Close()
	if err := json.NewDecoder(bytes.NewReader(s)).Decode(&db); err != nil {
		return err
	}
	return nil
}
