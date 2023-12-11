package main

import (
	"fmt"
	"os"

	"go.etcd.io/bbolt"
)

func main() {
	dbFile := "/tmp/xx.db"

	db, err := bbolt.Open(dbFile, 0600, nil)
	if err != nil {
		panic(err)
	}
	defer func() {
		db.Close()
		os.Remove(dbFile)
	}()

	err = db.Update(func(tx *bbolt.Tx) error {
		_, txerr := tx.CreateBucketIfNotExists([]byte("mock"))
		return txerr
	})
	if err != nil {
		panic(err)
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		bk := tx.Bucket([]byte("mock"))

		if err := bk.Put(idToBytes(1), make([]byte, 1000)); err != nil {
			return err
		}

		if err := bk.Put(idToBytes(2), make([]byte, 1000)); err != nil {
			return err
		}

		if err := bk.Put(idToBytes(3), make([]byte, 1000)); err != nil {
			return err
		}

		if err := bk.Put(idToBytes(4), make([]byte, 1000)); err != nil {
			return err
		}

		if err := bk.Put(idToBytes(10), make([]byte, 1000)); err != nil {
			return err
		}

		if err := bk.Put(idToBytes(11), make([]byte, 1000)); err != nil {
			return err
		}

		if err := bk.Put(idToBytes(12), make([]byte, 1000)); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		bk := tx.Bucket([]byte("mock"))

		bbolt.ModifyKeyBeforeStore = func(_ []byte) []byte {
			return idToBytes(1)
		}

		if err := bk.Put(idToBytes(13), make([]byte, 100)); err != nil {
			return err
		}

		bbolt.ModifyKeyBeforeStore = nil

		return nil
	})
	if err != nil {
		panic(err)
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		bk := tx.Bucket([]byte("mock"))

		if err := bk.Put(idToBytes(100), make([]byte, 100)); err != nil {
			return err
		}
		return nil

	})
	if err != nil {
		panic(err)
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		bk := tx.Bucket([]byte("mock"))

		for i := 101; i < 105; i++ {
			if err := bk.Put(idToBytes(i), make([]byte, 1000)); err != nil {
				return err
			}
		}
		return nil

	})
	if err != nil {
		panic(err)
	}
}

func idToBytes(id int) []byte {
	return []byte(fmt.Sprintf("%010d", id))
}
