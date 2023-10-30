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

var (
	// DefaultBackoff is the default backoff function for retries.
	DefaultBackoff = exponentialBackoff
	// DefaultRetry is the default check-for-retry function for retries.
	DefaultRetry = RetryOnError
	// DefaultRetries is the default number of times a request is tried.
	DefaultRetries = 5
	// DefaultRequestTimeout is the default request timeout.
	DefaultRequestTimeout = time.Second * 30
	// DefaultConnectionTimeout is the default connection timeout.
	DefaultConnectionTimeout = time.Second * 5
	// DefaultPoolSize sets the connection pool size.
	DefaultPoolSize = 100
	// DefaultPoolTTL sets the connection pool ttl.
	DefaultPoolTTL = time.Minute
)

type Options struct {
	CallOptions CallOptions

	Router Router

	Registry  registry.Registry
	Selector  selector.Selector
	Transport transport.Transport

	Broker broker.Broker

	Logger logger.Logger

	Context context.Context
	Codecs  map[string]codec.NewCodec

	Cache *Cache

	ContentType string

	Wrappers []Wrapper

	PoolSize int
	PoolTTL  time.Duration
}

type CallOptions struct {
	Context       context.Context
	Backoff       BackoffFunc
	Retry         RetryFunc
	SelectOptions []selector.SelectOption

	Address      []string
	CallWrappers []CallWrapper

	ConnectionTimeout time.Duration
	RequestTimeout    time.Duration
	StreamTimeout     time.Duration
	CacheExpiry       time.Duration
	DialTimeout       time.Duration
	Retries           int
	ServiceToken      bool
	ConnClose         bool
}

type PublishOptions struct {
	Context  context.Context
	Exchange string
}

type MessageOptions struct {
	ContentType string
}

type RequestOptions struct {
	Context     context.Context
	ContentType string
	Stream      bool
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

func WithConnClose() CallOption {
	return func(o *CallOptions) {
		o.ConnClose = true
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
