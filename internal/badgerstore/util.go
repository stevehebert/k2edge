package badgerstore

import (
	"context"
	"fmt"
	"strings"

	"github.com/dgraph-io/badger/v3"
	"github.com/stevehebert/k2edge/internal/persistence"
)

type badgerStore struct {
	badger *badger.DB
}

func (bi *badgerStore) Get(key string) (value []byte, err error) {
	err = bi.badger.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))

		if err != nil {
			if strings.Contains(err.Error(), "Key not found") {
				return nil
			}
			return err
		}
		tmpValue, err := item.ValueCopy(nil)

		if err != nil {
			return err
		}

		value = make([]byte, len(tmpValue))
		copy(value, tmpValue)
		return nil
	})
	return

}

func (bi *badgerStore) Set(key []byte, value []byte) error {
	fmt.Println("setting")
	if key == nil {
		return nil
	}

	return bi.badger.Update(func(txn *badger.Txn) error {
		entry := badger.NewEntry(key, value)

		err := txn.SetEntry(entry)
		return err
	})
}

func (bi *badgerStore) GetMetrics() interface{} {
	return nil

}

func (bi *badgerStore) Ready(cxt context.Context) bool {
	return true
}

func Connect() (persistence.StorageFacade, error) {
	b, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	bin := &badgerStore{
		badger: b,
	}

	return bin, err
}
