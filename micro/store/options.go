package store

import (
	"context"
	"time"

	"github.com/wxc/micro/client"
	"github.com/wxc/micro/logger"
)

type Options struct {
	Nodes    []string
	Database string
	Table    string
	Context  context.Context
	Client   client.Client
	Logger   logger.Logger
}

type Option func(o *Options)

func Nodes(a ...string) Option {
	return func(o *Options) {
		o.Nodes = a
	}
}

func Database(db string) Option {
	return func(o *Options) {
		o.Database = db
	}
}

func Table(t string) Option {
	return func(o *Options) {
		o.Table = t
	}
}

func WithContext(c context.Context) Option {
	return func(o *Options) {
		o.Context = c
	}
}

func WithClient(c client.Client) Option {
	return func(o *Options) {
		o.Client = c
	}
}

func WithLogger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

type ReadOptions struct {
	Database, Table string
	Prefix          bool
	Suffix          bool
	Limit           uint
	Offset          uint
}

type ReadOption func(r *ReadOptions)

func ReadFrom(database, table string) ReadOption {
	return func(r *ReadOptions) {
		r.Database = database
		r.Table = table
	}
}

func ReadPrefix() ReadOption {
	return func(r *ReadOptions) {
		r.Prefix = true
	}
}

func ReadSuffix() ReadOption {
	return func(r *ReadOptions) {
		r.Suffix = true
	}
}

func ReadLimit(l uint) ReadOption {
	return func(r *ReadOptions) {
		r.Limit = l
	}
}

func ReadOffset(o uint) ReadOption {
	return func(r *ReadOptions) {
		r.Offset = o
	}
}

type WriteOptions struct {
	Database, Table string
	Expiry          time.Time
	TTL             time.Duration
}

type WriteOption func(w *WriteOptions)

func WriteTo(database, table string) WriteOption {
	return func(w *WriteOptions) {
		w.Database = database
		w.Table = table
	}
}

func WriteExpiry(t time.Time) WriteOption {
	return func(w *WriteOptions) {
		w.Expiry = t
	}
}

func WriteTTL(d time.Duration) WriteOption {
	return func(w *WriteOptions) {
		w.TTL = d
	}
}

type DeleteOptions struct {
	Database, Table string
}

type DeleteOption func(d *DeleteOptions)

func DeleteFrom(database, table string) DeleteOption {
	return func(d *DeleteOptions) {
		d.Database = database
		d.Table = table
	}
}

type ListOptions struct {
	Database, Table string
	Prefix          string
	Suffix          string
	Limit           uint
	Offset          uint
}

type ListOption func(l *ListOptions)

func ListFrom(database, table string) ListOption {
	return func(l *ListOptions) {
		l.Database = database
		l.Table = table
	}
}

func ListPrefix(p string) ListOption {
	return func(l *ListOptions) {
		l.Prefix = p
	}
}

func ListSuffix(s string) ListOption {
	return func(l *ListOptions) {
		l.Suffix = s
	}
}

func ListLimit(l uint) ListOption {
	return func(lo *ListOptions) {
		lo.Limit = l
	}
}

func ListOffset(o uint) ListOption {
	return func(l *ListOptions) {
		l.Offset = o
	}
}
