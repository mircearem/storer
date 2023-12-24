package store

import (
	"fmt"
	"os"

	"go.etcd.io/bbolt"
)

const (
	defaultDBName = "app"
	ext           = "db"
)

type Map map[string]string

type Store struct {
	name string
	*Options
	db *bbolt.DB
}

func NewStore(options ...OptFunc) (*Store, error) {
	opts := &Options{
		DBName: defaultDBName,
	}
	for _, fn := range options {
		fn(opts)
	}
	dbname := fmt.Sprintf("%s.%s", opts.DBName, ext)
	db, err := bbolt.Open(dbname, 0666, nil)
	if err != nil {
		return nil, err
	}

	return &Store{
		name:    dbname,
		db:      db,
		Options: opts,
	}, nil
}

func (r *Store) DropDatabase(name string) error {
	dbname := fmt.Sprintf("%s.%s", name, ext)
	return os.Remove(dbname)
}

func (r *Store) CreateCollection(name string) (*bbolt.Bucket, error) {
	tx, err := r.db.Begin(true)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	b, err := tx.CreateBucketIfNotExists([]byte(name))
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (r *Store) Collection(name string) *Filter {
	return NewFilter(r, name)
}
