package memory

import (
	"context"

	"github.com/wxc/micro/config/source"
)

type changeSetKey struct{}

func withData(d []byte, f string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, changeSetKey{}, &source.ChangeSet{
			Data:   d,
			Format: f,
		})
	}
}

func WithChangeSet(cs *source.ChangeSet) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, changeSetKey{}, cs)
	}
}

func WithJSON(d []byte) source.Option {
	return withData(d, "json")
}

func WithYAML(d []byte) source.Option {
	return withData(d, "yaml")
}
