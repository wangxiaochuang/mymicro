package source

import (
	"context"

	"github.com/wxc/micro/client"
	"github.com/wxc/micro/config/encoder"
	"github.com/wxc/micro/config/encoder/json"
)

type Options struct {
	Encoder encoder.Encoder
	Context context.Context
	Client  client.Client
}

type Option func(o *Options)

func NewOptions(opts ...Option) Options {
	options := Options{
		Encoder: json.NewEncoder(),
		Context: context.Background(),
		Client:  client.DefaultClient,
	}

	for _, o := range opts {
		o(&options)
	}
	return options
}

func WithEncoder(e encoder.Encoder) Option {
	return func(o *Options) {
		o.Encoder = e
	}
}

func WithClient(c client.Client) Option {
	return func(o *Options) {
		o.Client = c
	}
}
