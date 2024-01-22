package store

import (
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
		msg := fmt.Sprintf("failed to insert key: (%s), with value: (%s), reason: (%s)",
			string(k), string(v), err.Error())
		return 0, NewStoreError(ERR_PUT_FAIL_UDEF, msg)
	}
	defer tx.Rollback()

	b, err := tx.CreateBucketIfNotExists([]byte(f.coll))
	if err != nil {
		return 0, err
	}
	val := b.Get(k)
	if val != nil {
		msg := fmt.Sprintf("key: (%s) already in collection: (%s)", string(k), f.coll)
		return 0, NewStoreError(ERR_PUT_FAIL_CONF, msg)
	}
	id, err := b.NextSequence()
	if err != nil {
		msg := fmt.Sprintf("failed to insert key: (%s), with value: (%s), reason: (%s)",
			string(k), string(v), err.Error())
		return 0, NewStoreError(ERR_PUT_FAIL_UDEF, msg)
	}
	if err := b.Put(k, v); err != nil {
		msg := fmt.Sprintf("failed to insert key: (%s), with value: (%s), reason: (%s)",
			string(k), string(v), err.Error())
		return 0, NewStoreError(ERR_PUT_FAIL_UDEF, msg)
	}
	if err := tx.Commit(); err != nil {
		msg := fmt.Sprintf("failed to insert key: (%s), with value: (%s), reason: (%s)",
			string(k), string(v), err.Error())
		return 0, NewStoreError(ERR_PUT_FAIL_UDEF, msg)
	}
	return id, nil
}

func (f *Filter) Update(k []byte, v []byte) error {
	tx, err := f.store.db.Begin(true)
	if err != nil {
		msg := fmt.Sprintf("failed to update key: (%s), with value: (%s), reason: (%s)",
			string(k), string(v), err.Error())
		return NewStoreError(ERR_UPD_FAIL_UDEF, msg)
	}
	defer tx.Rollback()

	b := tx.Bucket([]byte(f.coll))
	if b == nil {
		msg := fmt.Sprintf("collection: (%s) does not exist", f.coll)
		return NewStoreError(ERR_COL_NOT_FOUND, msg)
	}
	val := b.Get(k)
	if val == nil {
		msg := fmt.Sprintf("key: (%s) not in collection: (%s)", string(k), f.coll)
		return NewStoreError(ERR_UPD_FAIL_NOTF, msg)
	}
	if err := b.Put(k, v); err != nil {
		msg := fmt.Sprintf("failed to update key: (%s), with value: (%s), reason: (%s)",
			string(k), string(v), err.Error())
		return NewStoreError(ERR_UPD_FAIL_UDEF, msg)
	}
	if err := tx.Commit(); err != nil {
		msg := fmt.Sprintf("failed to update key: (%s), with value: (%s), reason: (%s)",
			string(k), string(v), err.Error())
		return NewStoreError(ERR_UPD_FAIL_UDEF, msg)
	}
	return nil
}

func (f *Filter) Get(k []byte) ([]byte, error) {
	tx, err := f.store.db.Begin(false)
	if err != nil {
		msg := fmt.Sprintf("failed to get value for key: (%s), reason: (%s)",
			string(k), err.Error())
		return nil, NewStoreError(ERR_GET_FAIL_UDEF, msg)
	}
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
	return bytes, nil
}

func (f *Filter) Delete(k []byte) error {
	tx, err := f.store.db.Begin(true)
	if err != nil {
		msg := fmt.Sprintf("failed to remove key: (%s) from collection: (%s), reason: (%s)",
			string(k), f.coll, err.Error())
		return NewStoreError(ERR_DEL_FAIL_UDEF, msg)
	}
	b := tx.Bucket([]byte(f.coll))
	if b == nil {
		msg := fmt.Sprintf("collection: (%s) does not exist", f.coll)
		return NewStoreError(ERR_COL_NOT_FOUND, msg)
	}
	bytes := b.Get(k)
	if bytes == nil {
		msg := fmt.Sprintf("key: (%s) not in collection: (%s)", string(k), f.coll)
		return NewStoreError(ERR_DEL_FAIL_NOTF, msg)

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
