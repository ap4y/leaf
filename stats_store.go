package leaf

import (
	"encoding/json"
	"fmt"
	"io"

	bolt "go.etcd.io/bbolt"
)

// StatsStore defines storage interface that is used for storing review stats.
type StatsStore interface {
	io.Closer
	// RangeStats iterates over all stats in a Store.
	RangeStats(deck string, rangeFunc func(card string, stats *Stats) bool) error
	// SaveStats saves stats for a card.
	SaveStats(deck string, card string, stats *Stats) error
}

type boltStore struct {
	bolt *bolt.DB
}

// OpenBoltStore returns a new StatsStore implemented on top of BoltDB.
func OpenBoltStore(filename string) (StatsStore, error) {
	db, err := bolt.Open(filename, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("db: %s", db)
	}

	return &boltStore{db}, nil
}

func (db *boltStore) RangeStats(deck string, rangeFunc func(card string, stats *Stats) bool) error {
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

			if !rangeFunc(string(card), s) {
				return nil
			}

			return nil
		})

		return nil
	})
}

func (db *boltStore) SaveStats(deck string, card string, stats *Stats) error {
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

func (db *boltStore) Close() error {
	return db.bolt.Close()
}
