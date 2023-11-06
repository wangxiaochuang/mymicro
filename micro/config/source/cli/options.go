package cli

import (
	"context"

	"github.com/urfave/cli/v2"
	"github.com/wxc/micro/config/source"
)

type contextKey struct{}

func Context(c *cli.Context) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, contextKey{}, c)
	}
}
