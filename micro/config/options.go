package config

import (
	"github.com/wxc/micro/config/loader"
	"github.com/wxc/micro/config/reader"
	"github.com/wxc/micro/config/source"
)

func WithLoader(l loader.Loader) Option {
	return func(o *Options) {
		o.Loader = l
	}
}

func WithSource(s source.Source) Option {
	return func(o *Options) {
		o.Source = append(o.Source, s)
	}
}

func WithReader(r reader.Reader) Option {
	return func(o *Options) {
		o.Reader = r
	}
}

func WithWatcherDisabled() Option {
	return func(o *Options) {
		o.WithWatcherDisabled = true
	}
}
