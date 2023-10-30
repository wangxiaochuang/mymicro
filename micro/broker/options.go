package broker

import (
	"context"
	"crypto/tls"

	"github.com/wxc/micro/registry"
	"go-micro.dev/v4/codec"
	"go-micro.dev/v4/logger"
)

type Options struct {
	Codec        codec.Marshaler
	Logger       logger.Logger
	Registry     registry.Registry
	Context      context.Context
	ErrorHandler Handler
	TLSConfig    *tls.Config
	Addrs        []string
	Secure       bool
}

type PublishOptions struct {
	Context context.Context
}

type SubscribeOptions struct {
	Context context.Context
	Queue   string
	AutoAck bool
}

type Option func(*Options)

type PublishOption func(*PublishOptions)

func PublishContext(ctx context.Context) PublishOption {
	return func(o *PublishOptions) {
		o.Context = ctx
	}
}

type SubscribeOption func(*SubscribeOptions)

func NewOptions(opts ...Option) *Options {
	options := Options{
		Context: context.Background(),
		Logger:  logger.DefaultLogger,
	}

	for _, o := range opts {
		o(&options)
	}

	return &options
}

func NewSubscribeOptions(opts ...SubscribeOption) SubscribeOptions {
	opt := SubscribeOptions{
		AutoAck: true,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
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

func DisableAutoAck() SubscribeOption {
	return func(o *SubscribeOptions) {
		o.AutoAck = false
	}
}

func ErrorHandler(h Handler) Option {
	return func(o *Options) {
		o.ErrorHandler = h
	}
}

func Queue(name string) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.Queue = name
	}
}

func Registry(r registry.Registry) Option {
	return func(o *Options) {
		o.Registry = r
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

func Logger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

func SubscribeContext(ctx context.Context) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.Context = ctx
	}
}
