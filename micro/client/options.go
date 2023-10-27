package client

import (
	"context"
	"time"

	"github.com/wxc/micro/broker"
	"github.com/wxc/micro/codec"
	"github.com/wxc/micro/logger"
	"github.com/wxc/micro/registry"
	"github.com/wxc/micro/selector"
	"github.com/wxc/micro/transport"
)

type Options struct {
	ContentType string

	Broker    broker.Broker
	Codecs    map[string]codec.NewCodec
	Registry  registry.Registry
	Selector  selector.Selector
	Transport transport.Transport

	Router Router

	PoolSize int
	PoolTTL  time.Duration

	Cache *Cache

	Wrappers []Wrapper

	CallOptions CallOptions

	Logger logger.Logger

	Context context.Context
}

type CallOptions struct {
	SelectOptions []selector.SelectOption

	Address        []string
	Backoff        BackoffFunc
	Retry          RetryFunc
	DialTimeout    time.Duration
	Retries        int
	RequestTimeout time.Duration
	StreamTimeout  time.Duration
	ServiceToken   bool
	CacheExpiry    time.Duration

	CallWrappers []CallWrapper

	Context context.Context
}

type PublishOptions struct {
	Exchange string
	Context  context.Context
}

type MessageOptions struct {
	ContentType string
}

type RequestOptions struct {
	ContentType string
	Stream      bool
	Context     context.Context
}

func NewOptions(options ...Option) Options {
	opts := Options{
		Cache:       NewCache(),
		Context:     context.Background(),
		ContentType: DefaultContentType,
		Codecs:      make(map[string]codec.NewCodec),
		CallOptions: CallOptions{
			Backoff:        DefaultBackoff,
			Retry:          DefaultRetry,
			Retries:        DefaultRetries,
			RequestTimeout: DefaultRequestTimeout,
			DialTimeout:    transport.DefaultDialTimeout,
		},
		PoolSize:  DefaultPoolSize,
		PoolTTL:   DefaultPoolTTL,
		Broker:    broker.DefaultBroker,
		Selector:  selector.DefaultSelector,
		Registry:  registry.DefaultRegistry,
		Transport: transport.DefaultTransport,
		Logger:    logger.DefaultLogger,
	}

	for _, o := range options {
		o(&opts)
	}

	return opts
}

func Broker(b broker.Broker) Option {
	return func(o *Options) {
		o.Broker = b
	}
}

func Codec(contentType string, c codec.NewCodec) Option {
	return func(o *Options) {
		o.Codecs[contentType] = c
	}
}

func ContentType(ct string) Option {
	return func(o *Options) {
		o.ContentType = ct
	}
}

func PoolSize(d int) Option {
	return func(o *Options) {
		o.PoolSize = d
	}
}

func PoolTTL(d time.Duration) Option {
	return func(o *Options) {
		o.PoolTTL = d
	}
}

func Registry(r registry.Registry) Option {
	return func(o *Options) {
		o.Registry = r
		o.Selector.Init(selector.Registry(r))
	}
}

func Transport(t transport.Transport) Option {
	return func(o *Options) {
		o.Transport = t
	}
}

func Selector(s selector.Selector) Option {
	return func(o *Options) {
		o.Selector = s
	}
}

func Wrap(w Wrapper) Option {
	return func(o *Options) {
		o.Wrappers = append(o.Wrappers, w)
	}
}

func WrapCall(cw ...CallWrapper) Option {
	return func(o *Options) {
		o.CallOptions.CallWrappers = append(o.CallOptions.CallWrappers, cw...)
	}
}

func Backoff(fn BackoffFunc) Option {
	return func(o *Options) {
		o.CallOptions.Backoff = fn
	}
}

func Retries(i int) Option {
	return func(o *Options) {
		o.CallOptions.Retries = i
	}
}

func Retry(fn RetryFunc) Option {
	return func(o *Options) {
		o.CallOptions.Retry = fn
	}
}

func RequestTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.CallOptions.RequestTimeout = d
	}
}

func StreamTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.CallOptions.StreamTimeout = d
	}
}

func DialTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.CallOptions.DialTimeout = d
	}
}

func WithExchange(e string) PublishOption {
	return func(o *PublishOptions) {
		o.Exchange = e
	}
}

func PublishContext(ctx context.Context) PublishOption {
	return func(o *PublishOptions) {
		o.Context = ctx
	}
}

func WithAddress(a ...string) CallOption {
	return func(o *CallOptions) {
		o.Address = a
	}
}

func WithSelectOption(so ...selector.SelectOption) CallOption {
	return func(o *CallOptions) {
		o.SelectOptions = append(o.SelectOptions, so...)
	}
}

func WithCallWrapper(cw ...CallWrapper) CallOption {
	return func(o *CallOptions) {
		o.CallWrappers = append(o.CallWrappers, cw...)
	}
}

func WithBackoff(fn BackoffFunc) CallOption {
	return func(o *CallOptions) {
		o.Backoff = fn
	}
}

func WithRetry(fn RetryFunc) CallOption {
	return func(o *CallOptions) {
		o.Retry = fn
	}
}

func WithRetries(i int) CallOption {
	return func(o *CallOptions) {
		o.Retries = i
	}
}

func WithRequestTimeout(d time.Duration) CallOption {
	return func(o *CallOptions) {
		o.RequestTimeout = d
	}
}

func WithStreamTimeout(d time.Duration) CallOption {
	return func(o *CallOptions) {
		o.StreamTimeout = d
	}
}

func WithDialTimeout(d time.Duration) CallOption {
	return func(o *CallOptions) {
		o.DialTimeout = d
	}
}

func WithServiceToken() CallOption {
	return func(o *CallOptions) {
		o.ServiceToken = true
	}
}

func WithCache(c time.Duration) CallOption {
	return func(o *CallOptions) {
		o.CacheExpiry = c
	}
}

func WithMessageContentType(ct string) MessageOption {
	return func(o *MessageOptions) {
		o.ContentType = ct
	}
}

func WithContentType(ct string) RequestOption {
	return func(o *RequestOptions) {
		o.ContentType = ct
	}
}

func StreamingRequest() RequestOption {
	return func(o *RequestOptions) {
		o.Stream = true
	}
}

func WithRouter(r Router) Option {
	return func(o *Options) {
		o.Router = r
	}
}

func WithLogger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}
