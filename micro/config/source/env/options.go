package env

import (
	"context"
	"strings"

	"github.com/wxc/micro/config/source"
)

type strippedPrefixKey struct{}
type prefixKey struct{}

func WithStrippedPrefix(p ...string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}

		o.Context = context.WithValue(o.Context, strippedPrefixKey{}, appendUnderscore(p))
	}
}

func WithPrefix(p ...string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, prefixKey{}, appendUnderscore(p))
	}
}

func appendUnderscore(prefixes []string) []string {
	var result []string
	for _, p := range prefixes {
		if !strings.HasSuffix(p, "_") {
			result = append(result, p+"_")
		}
		result = append(result, p)
	}
	return result
}
