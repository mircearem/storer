package store

import (
	"fmt"

	"go.etcd.io/bbolt"
)

// Decode key - value json messages
type Map map[string]string

// JSON messages for insertion in the database
type Message struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	TTL   int32  `json:"ttl"`
}

type Store struct {
	name string
	cfg  StoreConfig
	db   *bbolt.DB
}

func NewStore(cfg StoreConfig) (*Store, error) {
	dbname := fmt.Sprintf("%s.db", cfg.DbName)
	var (
		db  *bbolt.DB
		err error
	)

	if cfg.timeout > 0 {
		db, err = bbolt.Open(dbname, 0666, &bbolt.Options{Timeout: cfg.timeout})
	} else {
		db, err = bbolt.Open(dbname, 0666, nil)
	}

	if err != nil {
		return nil, err
	}

	return &Store{
		name: dbname,
		db:   db,
		cfg:  cfg,
	}, nil
}

func (s *Store) CreateCollection(name string) (*bbolt.Bucket, error) {
	tx, err := s.db.Begin(true)
	if err != nil {
		msg := fmt.Sprintf("failed to create collection: (%s), reason: (%s)", name, err.Error())
		return nil, NewStoreError(ERR_COL_MAKE_FAIL, msg)
	}
	defer tx.Rollback()

	b, err := tx.CreateBucketIfNotExists([]byte(name))
	if err != nil {
		msg := fmt.Sprintf("failed to create collection: (%s), reason: (%s)", name, err.Error())
		return nil, NewStoreError(ERR_COL_MAKE_FAIL, msg)
	}
	return b, nil
}

func (s *Store) DeleteCollection(name string) error {
	tx, err := s.db.Begin(true)
	if err != nil {
		msg := fmt.Sprintf("failed to delete collection: (%s), reason: (%s)", name, err.Error())
		return NewStoreError(ERR_COL_DROP_FAIL, msg)
	}
	defer tx.Rollback()

	err = tx.DeleteBucket([]byte(name))
	if err != nil {
		msg := fmt.Sprintf("failed to delete collection: (%s), reason: (%s)", name, err.Error())
		return NewStoreError(ERR_COL_DROP_FAIL, msg)
	}
	// This transaction should be commited
	return nil
}

func (s *Store) Collection(name string) *Filter {
	return NewFilter(s, name)
}
