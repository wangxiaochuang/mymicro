package transport

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	"github.com/wxc/micro/codec"
	"github.com/wxc/micro/logger"
)

var (
	DefaultBufSizeH2 = 4 * 1024 * 1024
)

type Options struct {
	Codec      codec.Marshaler
	Context    context.Context
	Logger     logger.Logger
	TLSConfig  *tls.Config
	Addrs      []string
	Timeout    time.Duration
	BuffSizeH2 int
	Secure     bool
}

type DialOptions struct {
	Context            context.Context
	Timeout            time.Duration
	Stream             bool
	ConnClose          bool
	InsecureSkipVerify bool
}

type ListenOptions struct {
	Context context.Context
}

func Addrs(addrs ...string) Option {
	return func(o *Options) {
		o.Addrs = addrs
	}
}

func Codec(c codec.Marshaler) Option {
	return func(o *Options) {
		o.Codec = c
	}
}

func Timeout(t time.Duration) Option {
	return func(o *Options) {
		o.Timeout = t
	}
}

func Secure(b bool) Option {
	return func(o *Options) {
		o.Secure = b
	}
}

func TLSConfig(t *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = t
	}
}

func WithStream() DialOption {
	return func(o *DialOptions) {
		o.Stream = true
	}
}

func WithTimeout(d time.Duration) DialOption {
	return func(o *DialOptions) {
		o.Timeout = d
	}
}

func WithConnClose() DialOption {
	return func(o *DialOptions) {
		o.ConnClose = true
	}
}

func WithInsecureSkipVerify(b bool) DialOption {
	return func(o *DialOptions) {
		o.InsecureSkipVerify = b
	}
}

func Logger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

func BuffSizeH2(size int) Option {
	return func(o *Options) {
		o.BuffSizeH2 = size
	}
}

func NetListener(customListener net.Listener) ListenOption {
	return func(o *ListenOptions) {
		if customListener == nil {
			return
		}

		if o.Context == nil {
			o.Context = context.TODO()
		}

		o.Context = context.WithValue(o.Context, netListener{}, customListener)
	}
}
