package store

import (
	"errors"
	"time"
)

var (
	ErrNotFound        = errors.New("not found")
	DefaultStore Store = NewStore()
)

type Store interface {
	Init(...Option) error
	Options() Options
	Read(key string, opts ...ReadOption) ([]*Record, error)
	Write(r *Record, opts ...WriteOption) error
	Delete(key string, opts ...DeleteOption) error
	List(opts ...ListOption) ([]string, error)
	Close() error
	String() string
}

type Record struct {
	Key      string                 `json:"key"`
	Value    []byte                 `json:"value"`
	Metadata map[string]interface{} `json:"metadata"`
	Expiry   time.Duration          `json:"expiry,omitempty"`
}

func NewStore(opts ...Option) Store {
	return NewMemoryStore(opts...)
}
