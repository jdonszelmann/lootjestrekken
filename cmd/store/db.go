package store

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
	"lootjestrekken/pkg/lootjestrekken"
	"os"
)

const BucketName = "trekkingen"

type DbStore struct {
	Db *bolt.DB
}

func (i *DbStore) UpdateTrekking(trekking lootjestrekken.Trekking) error {
	log.Debugf("updating trekking with name %s in store", trekking.Name)

	jsont, err := json.Marshal(trekking)
	if err != nil {
		return err
	}

	err = i.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketName))
		err := b.Put([]byte(trekking.Name), jsont)
		return err
	})

	return err
}

func (i *DbStore) GetTrekking(name string) (lootjestrekken.Trekking, error) {
	log.Debugf("getting trekking with name %s from store", name)

	var t lootjestrekken.Trekking
	err := i.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketName))
		v := b.Get([]byte(name))
		if err := json.Unmarshal(v, &t); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return lootjestrekken.Trekking{}, err
	}

	return t, nil
}

func (i *DbStore) GetTrekkingNames() ([]string, error) {
	log.Debug("getting all trekking names from store")

	keys := make([]string, 0)

	err := i.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketName))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var t lootjestrekken.Trekking
			if err := json.Unmarshal(v, &t); err != nil {
				return err
			}

			keys = append(keys, t.Name)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (i *DbStore) GetTrekkingInfos() ([]string, error) {
	log.Debug("getting all trekking infos from store")

	keys := make([]string, 0)

	err := i.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketName))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var t lootjestrekken.Trekking
			if err := json.Unmarshal(v, &t); err != nil {
				return err
			}

			keys = append(keys, t.GetInfo())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (i *DbStore) AddTrekking(name string, trekking lootjestrekken.Trekking) error {
	log.Debugf("Adding trekking with name %s to store", name)

	trekking.Name = name

	jsont, err := json.Marshal(trekking)
	if err != nil {
		return err
	}

	err = i.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketName))
		if b.Get([]byte(name)) != nil {
			return ErrExists
		}

		err := b.Put([]byte(name), jsont)
		return err
	})

	return err
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func NewDbStore(location string) (*DbStore, error) {
	dbloc := location + "/lootjestrekken.db"
	if !exists(location) {
		return nil, fmt.Errorf("directory %s does not exist", location)
	}

	db, err := bolt.Open(dbloc, 0666, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(BucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	return &DbStore {
		Db: db,
	}, nil
}


