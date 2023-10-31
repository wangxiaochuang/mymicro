package loader

import (
	"context"

	"github.com/wxc/micro/config/reader"
	"github.com/wxc/micro/config/source"
)

type Loader interface {
	Close() error
	Load(...source.Source) error
	Snapshot() (*Snapshot, error)
	Sync() error
	Watch(...string) (Watcher, error)
	String() string
}

type Watcher interface {
	Next() (*Snapshot, error)
	Stop() error
}

type Snapshot struct {
	ChangeSet *source.ChangeSet
	Version   string
}

type Options struct {
	Reader              reader.Reader
	Context             context.Context
	Source              []source.Source
	WithWatcherDisabled bool
}

type Option func(o *Options)

func Copy(s *Snapshot) *Snapshot {
	cs := *(s.ChangeSet)

	return &Snapshot{
		ChangeSet: &cs,
		Version:   s.Version,
	}
}
