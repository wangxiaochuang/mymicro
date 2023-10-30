package cache

import (
	"time"

	"github.com/wxc/micro/logger"
)

func WithTTL(t time.Duration) Option {
	return func(o *Options) {
		o.TTL = t
	}
}

func WithLogger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}
