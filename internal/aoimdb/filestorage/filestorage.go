package filestorage

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/minmax1996/aoimdb/logger"
)

var databaseBackupFileName = "dbbackup.aoimdb"

func RestoreFromBackup(db interface{}) error {
	data, err := ReadFromFile(databaseBackupFileName)
	if err != nil {
		return err
	}
	return DecompressAndDecore(db, data)
}

func StartBackups(db interface{}) {
	go func() {
		for {
			time.Sleep(3 * time.Second)
			zipped, err := EncodeAndCompress(db)
			if err != nil {
				logger.Error(err.Error())
				continue
			}
			WriteToFile(zipped, databaseBackupFileName)
		}
	}()
}

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

func WriteToFile(s []byte, file string) {
	f, err := os.Create(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	f.Write(s)
}

func EncodeAndCompress(p interface{}) ([]byte, error) {
	buf := bytes.Buffer{}

	if err := gob.NewEncoder(&buf).Encode(p); err != nil {
		return nil, err
	}

	fmt.Println("uncompressed size (bytes): ", len(buf.Bytes()))

	zipbuf := bytes.Buffer{}
	zipped := gzip.NewWriter(&zipbuf)
	zipped.Write(buf.Bytes())
	zipped.Close()

	fmt.Println("compressed size (bytes): ", len(zipbuf.Bytes()))
	return zipbuf.Bytes(), nil
}

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

	fmt.Println("uncompressed size (bytes): ", len(data))

	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&db); err != nil {
		return err
	}

	return nil
}
