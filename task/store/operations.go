package store

import (
	"fmt"
	"time"
	"log"
	
	"github.com/boltdb/bolt"
)

const dbFile = "task.db"
const dbBucket = "taskBucket"

func Add(task string) {
	db := openDB()
	defer db.Close()
	err := db.Update(func (tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbBucket))
		if b == nil {
			return fmt.Errorf("Could not get store", dbBucket)
		}

		e := b.Put([]byte(task), []byte(""))
		if e != nil {
			return fmt.Errorf("%v", e)
		}
		return nil
	})
	if err != nil {
		log.Fatalln(err)
	}
}

func List() []string {
	db := openDB()
	defer db.Close()

	var tasks []string

	err := db.View(func (tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbBucket))
		if b == nil {
			return fmt.Errorf("Could not get store", dbBucket)
		}

		b.ForEach(func (k, v []byte) error {
			tasks = append(tasks, string(k))
			return nil
		})
		return nil
	})

	if err != nil {
		log.Fatalln(err)
	}

	return tasks

}

func Done(taskNum int) error {
	tasks := List()
	if taskNum<1 || taskNum > len(tasks) {
		return fmt.Errorf("Task %d not found", taskNum)
	}
	task := tasks[taskNum-1]
	db := openDB()
	defer db.Close()
	err := db.Update(func (tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbBucket))
		if b == nil {
			return fmt.Errorf("Could not get store", dbBucket)
		}
		e := b.Delete([]byte(task))
		if e != nil {
			return fmt.Errorf("%v", e)
		}

		return nil
	})

	if err != nil {
		log.Fatalln(err)
	}

	return nil
}

func openDB() *bolt.DB {
	db, err := bolt.Open(dbFile, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatalln(err)
	}

	err = db.Update(func (tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(dbBucket))
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Initializing store %s: %s\n", dbBucket, err)
	}
	return db
}