package selector

import (
	"context"

	"github.com/wxc/micro/logger"
	"github.com/wxc/micro/registry"
)

type Options struct {
	Registry registry.Registry
	Strategy Strategy
	Context  context.Context
	Logger   logger.Logger
}

type SelectOptions struct {
	Context  context.Context
	Strategy Strategy

	Filters []Filter
}

type Option func(*Options)

type SelectOption func(*SelectOptions)

func Registry(r registry.Registry) Option {
	return func(o *Options) {
		o.Registry = r
	}
}

func SetStrategy(fn Strategy) Option {
	return func(o *Options) {
		o.Strategy = fn
	}
}

func WithFilter(fn ...Filter) SelectOption {
	return func(o *SelectOptions) {
		o.Filters = append(o.Filters, fn...)
	}
}

func WithStrategy(fn Strategy) SelectOption {
	return func(o *SelectOptions) {
		o.Strategy = fn
	}
}

func WithLogger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}
