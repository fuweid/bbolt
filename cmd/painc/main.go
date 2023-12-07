package main

import (
	"log"
	"math/rand"
	"os"
	"path/filepath"

	bolt "go.etcd.io/bbolt"
)

var BC = make(chan []byte)

func main() {
	dir, err := os.MkdirTemp("", "testXXX")
	if err != nil {
		panic(err)
	}
	defer func() {
		os.RemoveAll(dir)
	}()

	go Add(BC, filepath.Join(dir, "test.db"))

	s := make([]byte, 30)
	for i := 0; i < 100000000; i++ {
		rand.Read(s)
		BC <- s
	}

	close(BC)

}

func Add(c chan []byte, bfile string) {
	const bname = `bc41`

	db, err := bolt.Open(bfile, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bname))
		if err != nil {
			return err //fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	more := true
	for more {
		err := db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bname))
			for i := 0; i < 100000; i++ {
				v := <-c
				if v == nil {
					more = false
					return nil
				}
				err = b.Put(v, []byte("0"))
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
			panic(err)
		}
	}
}
