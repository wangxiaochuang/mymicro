package flag

import (
	"context"

	"github.com/wxc/micro/config/source"
)

type includeUnsetKey struct{}

func IncludeUnset(b bool) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, includeUnsetKey{}, true)
	}
}
