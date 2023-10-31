package memory

import (
	"github.com/wxc/micro/config/loader"
	"github.com/wxc/micro/config/reader"
	"github.com/wxc/micro/config/source"
)

func WithSource(s source.Source) loader.Option {
	return func(o *loader.Options) {
		o.Source = append(o.Source, s)
	}
}

func WithReader(r reader.Reader) loader.Option {
	return func(o *loader.Options) {
		o.Reader = r
	}
}

func WithWatcherDisabled() loader.Option {
	return func(o *loader.Options) {
		o.WithWatcherDisabled = true
	}
}
