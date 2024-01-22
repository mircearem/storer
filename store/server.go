package store

import (
	"fmt"
	"time"

	_ "github.com/mircearem/storer/log"
	"github.com/sirupsen/logrus"
	"go.etcd.io/bbolt"
)

type StorageServer struct {
	DB       Store
	scantime time.Duration
	Closech  chan struct{}
	errch    chan error
}

func NewStorageServer(s Store) *StorageServer {
	return &StorageServer{
		DB:      s,
		Closech: make(chan struct{}),
		errch:   make(chan error),
	}
}

func (s *StorageServer) Run() error {
	if s.DB.cfg.autoclean {
		go s.autoclean()
	}

	ticker := time.NewTicker(s.scantime)
	for {
		select {
		case <-ticker.C:
			// scan the databases for expiring items
		case <-s.Closech:
			// gracefull shutdown
			return nil
		case err := <-s.errch:
			// error channel, sort and log the errors
			logrus.Warn(err)
		}
	}
}

func (s *StorageServer) autoclean() {}

// @TODO implement in production
func (s *Store) handleTTL(collection string) {
	ticker := time.NewTicker(15 * time.Second)
	for {
		<-ticker.C
		tx, err := s.db.Begin(true)
		if err != nil {
			continue
		}
		defer tx.Rollback()

		err = tx.ForEach(deleteOldRecords)
		if err != nil {
			logrus.Warn(err)
		}
	}
}

// @ TODO implement in production
func deleteOldRecords(name []byte, b *bbolt.Bucket) error {
	msg := fmt.Sprintf("cleaning outdated records for bucket: (%s)", string(name))
	logrus.Info(msg)

	return b.ForEach(func(key, value []byte) error {
		it := b.Cursor()
		for k, v := it.First(); k != nil; k, v = it.Next() {
			_ = v
		}

		return nil
	})
}
