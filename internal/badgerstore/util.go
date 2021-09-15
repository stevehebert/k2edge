package badgerstore

import "github.com/dgraph-io/badger/v3"

var Badger *badger.DB

func Connect() (*badger.DB, error) {
	b, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	Badger = b

	return b, err
}
