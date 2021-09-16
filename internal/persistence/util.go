package persistence

import "context"

type StorageFacade interface {
	Get(key string) (value []byte, err error)
	Set(key []byte, value []byte) error
	GetMetrics() interface{}
	Ready(cxt context.Context) bool
}
