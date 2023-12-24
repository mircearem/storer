package store

import (
	"errors"
	"fmt"
)

type Filter struct {
	store *Store
	coll  string
}

func NewFilter(db *Store, coll string) *Filter {
	return &Filter{
		store: db,
		coll:  coll,
	}
}

func (f *Filter) Put(k []byte, v []byte) (uint64, error) {
	tx, err := f.store.db.Begin(true)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	b, err := tx.CreateBucketIfNotExists([]byte(f.coll))
	if err != nil {
		return 0, err
	}
	// Check if the key is already in the db and
	// return an error if the key is found
	val := b.Get(k)
	if val != nil {
		return 0, fmt.Errorf("key: (%s) is already in the collection", string(k))
	}

	id, err := b.NextSequence()
	if err != nil {
		return 0, err
	}
	if err := b.Put(k, v); err != nil {
		return 0, err
	}
	return id, tx.Commit()
}

func (f *Filter) Update(k []byte, v []byte) error {
	tx, err := f.store.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	b := tx.Bucket([]byte(f.coll))
	if b == nil {
		msg := fmt.Sprintf("collection (%s) does not exist", f.coll)
		err := errors.New(msg)
		return err
	}
	// Check if the key is already in the db and
	// return an error if the key is found
	val := b.Get(k)
	if val == nil {
		return fmt.Errorf("key: (%s) is not in collection (%s)", string(k), f.coll)
	}

	if err := b.Put(k, v); err != nil {
		return err
	}
	return tx.Commit()
}

func (f *Filter) Get(k []byte) ([]byte, error) {
	tx, err := f.store.db.Begin(false)
	if err != nil {
		return nil, err
	}
	b := tx.Bucket([]byte(f.coll))
	if b == nil {
		return nil, fmt.Errorf("collection (%s) not found", f.coll)
	}
	return b.Get(k), nil
}

func (f *Filter) Delete(k []byte) error {
	tx, err := f.store.db.Begin(true)
	if err != nil {
		return fmt.Errorf("collection %s does not exist", f.coll)
	}
	b := tx.Bucket([]byte(f.coll))
	if err := b.Delete([]byte(k)); err != nil {
		return err
	}
	return tx.Commit()
}
