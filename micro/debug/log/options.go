package log

import "time"

type Option func(*Options)

type Options struct {
	Format FormatFunc
	Name   string
	Size   int
}

func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

func Size(s int) Option {
	return func(o *Options) {
		o.Size = s
	}
}

func Format(f FormatFunc) Option {
	return func(o *Options) {
		o.Format = f
	}
}

func DefaultOptions() Options {
	return Options{
		Size: DefaultSize,
	}
}

type ReadOptions struct {
	Since  time.Time
	Count  int
	Stream bool
}

type ReadOption func(*ReadOptions)

func Since(s time.Time) ReadOption {
	return func(o *ReadOptions) {
		o.Since = s
	}
}

func Count(c int) ReadOption {
	return func(o *ReadOptions) {
		o.Count = c
	}
}
