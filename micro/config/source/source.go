package source

import (
	"errors"
	"time"
)

var (
	ErrWatcherStopped = errors.New("watcher stopped")
)

type Source interface {
	Read() (*ChangeSet, error)
	Write(*ChangeSet) error
	Watch() (Watcher, error)
	String() string
}

type ChangeSet struct {
	Timestamp time.Time
	Checksum  string
	Format    string
	Source    string
	Data      []byte
}

type Watcher interface {
	Next() (*ChangeSet, error)
	Stop() error
}
