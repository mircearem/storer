package store

import (
	"fmt"

	_ "github.com/mircearem/storer/log"
	"go.etcd.io/bbolt"
)

type dberror uint8

const (
	ERR_COL_MAKE_FAIL dberror = iota + 1
	ERR_COL_DROP_FAIL
	ERR_COL_NOT_FOUND // Collection not found
	ERR_PUT_FAIL_UDEF // Put fail, undefined
	ERR_PUT_FAIL_CONF // Put fail, conflict
	ERR_GET_FAIL_UDEF // Get fail, undefined
	ERR_GET_FAIL_NOTF // Get fail, not found
	ERR_UPD_FAIL_UDEF // Update fail, undefined
	ERR_UPD_FAIL_NOTF // Update fail, not found
	ERR_DEL_FAIL_UDEF
	ERR_DEL_FAIL_NOTF
)

type StoreConfig struct {
	DbName    string
	autoclean bool
}

func NewStoreConfig() StoreConfig {
	return StoreConfig{
		DbName:    "app",
		autoclean: false,
	}
}

func (s StoreConfig) WithDbName(name string) StoreConfig {
	s.DbName = name
	return s
}

func (s StoreConfig) WithAutoclean(b bool) StoreConfig {
	s.autoclean = b
	return s
}

type Map map[string]string

type Store struct {
	name string
	cfg  StoreConfig
	db   *bbolt.DB
}

type StoreError struct {
	errorType    dberror
	errorMessage string
}

func (s StoreError) Error() string {
	return s.errorMessage
}

func (s StoreError) Type() dberror {
	return s.errorType
}

func NewStoreError(typ dberror, msg string) StoreError {
	return StoreError{
		errorType:    typ,
		errorMessage: msg,
	}
}

func NewStore(cfg StoreConfig) (*Store, error) {
	dbname := fmt.Sprintf("%s.db", cfg.DbName)
	db, err := bbolt.Open(dbname, 0666, nil)
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
	return nil
}

func (s *Store) Collection(name string) *Filter {
	return NewFilter(s, name)
}
