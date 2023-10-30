package config

import (
	"context"

	"github.com/wxc/micro/config/loader"
	"github.com/wxc/micro/config/reader"
	"github.com/wxc/micro/config/source"
)

type Config interface {
	reader.Values
	Init(opts ...Option) error
	Options() Options
	Close() error
	Load(source ...source.Source) error
	Sync() error
	Watch(path ...string) (Watcher, error)
}

type Watcher interface {
	Next() (reader.Value, error)
	Stop() error
}

type Options struct {
	Loader loader.Loader
	Reader reader.Reader

	// for alternative data
	Context context.Context

	Source []source.Source

	WithWatcherDisabled bool
}

type Option func(o *Options)

var (
	DefaultConfig, _ = NewConfig()
)

func NewConfig(opts ...Option) (Config, error) {
	return newConfig(opts...)
}

func Bytes() []byte {
	return DefaultConfig.Bytes()
}

func Map() map[string]interface{} {
	return DefaultConfig.Map()
}

func Scan(v interface{}) error {
	return DefaultConfig.Scan(v)
}

func Sync() error {
	return DefaultConfig.Sync()
}

func Get(path ...string) reader.Value {
	return DefaultConfig.Get(path...)
}

func Load(source ...source.Source) error {
	return DefaultConfig.Load(source...)
}

func Watch(path ...string) (Watcher, error) {
	return DefaultConfig.Watch(path...)
}

func LoadFile(path string) error {
	panic("in LoadFile")
}
