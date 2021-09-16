package badgerstore

import (
	"context"

	"github.com/dgraph-io/badger/v3"
)

type badgerInteraction struct {
	badger *badger.DB
}

func (bi *badgerInteraction) Get(key string) (value []byte, err error) {
	return nil, nil

}

func (bi *badgerInteraction) Set(key []byte, value []byte) error {
	return nil

}

func (bi *badgerInteraction) GetMetrics() interface{} {
	return nil

}

func (bi *badgerInteraction) Ready(cxt context.Context) bool {
	return true
}

func Connect() (Interaction, error) {
	b, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	bin := &badgerInteraction{
		badger: b,
	}

	return bin, err
}

type Interaction interface {
	Get(key string) (value []byte, err error)
	Set(key []byte, value []byte) error
	GetMetrics() interface{}
	Ready(cxt context.Context) bool
}
