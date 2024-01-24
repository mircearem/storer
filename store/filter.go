package store

import (
	"fmt"
	"time"

	_ "github.com/mircearem/storer/log"
	"github.com/sirupsen/logrus"
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

func (f *Filter) Set(k []byte, v []byte, ttl int32) (uint64, error) {
	tx, err := f.store.db.Begin(true)
	if err != nil {
		msg := fmt.Sprintf("failed to insert key: (%s), with value: (%s), reason: (%s)",
			string(k), string(v), err.Error())
		return 0, NewStoreError(ERR_PUT_FAIL, msg)
	}
	defer tx.Rollback()

	b, err := tx.CreateBucketIfNotExists([]byte(f.coll))
	if err != nil {
		msg := fmt.Sprintf("failed to insert key: (%s), with value: (%s), reason: (%s)",
			string(k), string(v), err.Error())
		return 0, NewStoreError(ERR_PUT_FAIL, msg)
	}
	if err := b.Put(k, v); err != nil {
		msg := fmt.Sprintf("failed to insert key: (%s), with value: (%s), reason: (%s)",
			string(k), string(v), err.Error())
		return 0, NewStoreError(ERR_PUT_FAIL, msg)
	}
	id, err := b.NextSequence()
	if err != nil {
		msg := fmt.Sprintf("failed to insert key: (%s), with value: (%s), reason: (%s)",
			string(k), string(v), err.Error())
		return 0, NewStoreError(ERR_PUT_FAIL, msg)
	}
	if err := tx.Commit(); err != nil {
		msg := fmt.Sprintf("failed to insert key: (%s), with value: (%s), reason: (%s)",
			string(k), string(v), err.Error())
		return 0, NewStoreError(ERR_PUT_FAIL, msg)
	}
	if ttl > 0 {
		go func() {
			expire := time.Duration(ttl) * time.Second
			<-time.After(expire)
			if err := f.remove(k); err != nil {
				logrus.Warn(err)
				return
			}
			msg := fmt.Sprintf("successfully deleted key: (%s), with value (%s) from collection: (%s)",
				string(k), string(v), f.coll)
			logrus.Warn(msg)
		}()
	}
	return id, nil
}

func (f *Filter) Get(k []byte) ([]byte, error) {
	tx, err := f.store.db.Begin(false)
	if err != nil {
		msg := fmt.Sprintf("failed to get value for key: (%s), reason: (%s)",
			string(k), err.Error())
		return nil, NewStoreError(ERR_GET_FAIL_UDEF, msg)
	}
	defer func() {
		if tx.DB() != nil {
			tx.Rollback()
		}
	}()

	b := tx.Bucket([]byte(f.coll))
	if b == nil {
		msg := fmt.Sprintf("collection: (%s) does not exist", f.coll)
		return nil, NewStoreError(ERR_COL_NOT_FOUND, msg)
	}
	bytes := b.Get(k)
	if bytes == nil {
		msg := fmt.Sprintf("key: (%s) not in collection: (%s)", string(k), f.coll)
		return nil, NewStoreError(ERR_GET_FAIL_NOTF, msg)
	}
	return bytes, tx.Rollback()
}

func (f *Filter) Delete(k []byte) error {
	tx, err := f.store.db.Begin(true)
	if err != nil {
		msg := fmt.Sprintf("failed to remove key: (%s) from collection: (%s), reason: (%s)",
			string(k), f.coll, err.Error())
		return NewStoreError(ERR_DEL_FAIL_UDEF, msg)
	}
	defer func() {
		if tx.DB() != nil {
			tx.Rollback()
		}
	}()

	b := tx.Bucket([]byte(f.coll))
	if b == nil {
		msg := fmt.Sprintf("collection: (%s) does not exist", f.coll)
		return NewStoreError(ERR_COL_NOT_FOUND, msg)
	}
	if err := b.Delete([]byte(k)); err != nil {
		msg := fmt.Sprintf("failed to remove key: (%s) from collection: (%s), reason: (%s)",
			string(k), f.coll, err.Error())
		return NewStoreError(ERR_DEL_FAIL_UDEF, msg)
	}
	if err := tx.Commit(); err != nil {
		msg := fmt.Sprintf("failed to remove key: (%s) from collection: (%s), reason: (%s)",
			string(k), f.coll, err.Error())
		return NewStoreError(ERR_DEL_FAIL_UDEF, msg)
	}
	return nil
}

func (f *Filter) remove(k []byte) error {
	// Start a new read write transaction
	tx, err := f.store.db.Begin(true)
	if err != nil {
		return err
	}
	// Make sure the transaction rolls back in the event of a panic.
	defer func() {
		if tx.DB() != nil {
			tx.Rollback()
		}
	}()
	// Get the bucket - since this is an autoclean function, the bucket is sure to exist
	b := tx.Bucket([]byte(f.coll))
	if err := b.Delete([]byte(k)); err != nil {
		return err
	}
	return tx.Commit()
}
