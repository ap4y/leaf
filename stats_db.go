package leaf

import (
	"encoding/json"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

type StatsDB struct {
	bolt *bolt.DB
}

func OpenStatsDB(filename string) (*StatsDB, error) {
	db, err := bolt.Open(filename, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("db: %s", db)
	}

	return &StatsDB{db}, nil
}

func (db *StatsDB) GetStats(deck string, rangeFunc func(card string, stats *Stats)) error {
	return db.bolt.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(deck))
		if err != nil {
			return err
		}

		err = b.ForEach(func(card, stats []byte) error {
			s := new(Stats)
			if err := json.Unmarshal(stats, s); err != nil {
				return fmt.Errorf("json: %s", err)
			}

			rangeFunc(string(card), s)
			return nil
		})

		return nil
	})
}

func (db *StatsDB) SaveStats(deck string, card string, stats *Stats) error {
	return db.bolt.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(deck))
		if err != nil {
			return err
		}

		data, err := json.Marshal(stats)
		if err != nil {
			return fmt.Errorf("json: %s", err)
		}

		return b.Put([]byte(card), data)
	})
}

func (db *StatsDB) Close() error {
	return db.bolt.Close()
}
