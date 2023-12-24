package store

import (
	"testing"
)

// func TestCreate(t *testing.T) {
// 	db, err := NewStore(WithDBName("ips"))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	_ = db
// }

func TestInsert(t *testing.T) {
	values := map[string]string{
		"5.8.198.125": "Bucharest",
		"5.8.198.139": "London",
		"5.8.198.165": "Berlin",
		"5.8.198.184": "Paris",
		"5.8.198.169": "New York",
	}

	db, err := NewStore(WithDBName("ips"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.DropDatabase("ips")
	// Insert key value pairs in store
	for k, v := range values {
		id, err := db.Collection("ips").Put([]byte(k), []byte(v))
		_ = id
		if err != nil {
			t.Fatal(err)
		}
	}
	// Check key value pairs
	for k := range values {
		v, err := db.Collection("ips").Get([]byte(k))
		if string(v) != values[k] {
			t.Fatal(err)
		}
	}
}

func TestDelete(t *testing.T) {
	values := map[string]string{
		"5.8.198.125": "Bucharest",
		"5.8.198.139": "London",
		"5.8.198.165": "Berlin",
		"5.8.198.184": "Paris",
		"5.8.198.169": "New York",
	}

	db, err := NewStore(WithDBName("ips"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.DropDatabase("ips")
	// Insert key value pairs in store
	for k, v := range values {
		id, err := db.Collection("ips").Put([]byte(k), []byte(v))
		_ = id
		if err != nil {
			t.Fatal(err)
		}
	}
	// Check key value pairs
	for k := range values {
		err := db.Collection("ips").Delete([]byte(k))
		if err != nil {
			t.Fatal(err)
		}
	}
	for k := range values {
		v, err := db.Collection("ips").Get([]byte(k))
		_ = err
		if v != nil {
			t.Fatalf("not correctly deleted key: %s; expected nil, found value: %s", k, string(v))
		}
	}
}
