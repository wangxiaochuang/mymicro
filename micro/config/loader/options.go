package loader

import (
	"github.com/wxc/micro/config/reader"
	"github.com/wxc/micro/config/source"
)

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
