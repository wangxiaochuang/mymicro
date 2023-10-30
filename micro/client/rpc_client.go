package client

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"
	"github.com/wxc/micro/codec"
	"github.com/wxc/micro/registry"
	"github.com/wxc/micro/selector"
	"github.com/wxc/micro/util/pool"
)

const (
	packageID = "go.micro.client"
)

type rpcClient struct {
	opts Options
	once atomic.Value
	pool pool.Pool

	seq uint64

	mu sync.RWMutex
}

func newRPCClient(opt ...Option) Client {
	opts := NewOptions(opt...)

	p := pool.NewPool(
		pool.Size(opts.PoolSize),
		pool.TTL(opts.PoolTTL),
		pool.Transport(opts.Transport),
	)

	rc := &rpcClient{
		opts: opts,
		pool: p,
		seq:  0,
	}
	rc.once.Store(false)

	c := Client(rc)

	// wrap in reverse
	for i := len(opts.Wrappers); i > 0; i-- {
		c = opts.Wrappers[i-1](c)
	}

	return c
}

func (r *rpcClient) newCodec(contentType string) (codec.NewCodec, error) {
	if c, ok := r.opts.Codecs[contentType]; ok {
		return c, nil
	}

	if cf, ok := DefaultCodecs[contentType]; ok {
		return cf, nil
	}

	return nil, fmt.Errorf("unsupported Content-Type: %s", contentType)
}

func (r *rpcClient) call(ctx context.Context, node *registry.Node, req Request, resp interface{}, opts CallOptions) error {
	panic("in call")
}

func (r *rpcClient) stream(ctx context.Context, node *registry.Node, req Request, opts CallOptions) (Stream, error) {
	panic("in stream")
}

func (r *rpcClient) Init(opts ...Option) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	size := r.opts.PoolSize
	ttl := r.opts.PoolTTL
	tr := r.opts.Transport

	for _, o := range opts {
		o(&r.opts)
	}

	// update pool configuration if the options changed
	if size != r.opts.PoolSize || ttl != r.opts.PoolTTL || tr != r.opts.Transport {
		// close existing pool
		if err := r.pool.Close(); err != nil {
			return errors.Wrap(err, "failed to close pool")
		}

		// create new pool
		r.pool = pool.NewPool(
			pool.Size(r.opts.PoolSize),
			pool.TTL(r.opts.PoolTTL),
			pool.Transport(r.opts.Transport),
		)
	}

	return nil
}

func (r *rpcClient) Options() Options {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.opts
}

func (r *rpcClient) next(request Request, opts CallOptions) (selector.Next, error) {
	panic("in next")
}

func (r *rpcClient) Call(ctx context.Context, request Request, response interface{}, opts ...CallOption) error {
	panic("in Call")
}

func (r *rpcClient) Stream(ctx context.Context, request Request, opts ...CallOption) (Stream, error) {
	panic("in Stream")
}

func (r *rpcClient) Publish(ctx context.Context, msg Message, opts ...PublishOption) error {
	panic("in Publish")
}

func (r *rpcClient) NewMessage(topic string, message interface{}, opts ...MessageOption) Message {
	return newMessage(topic, message, r.opts.ContentType, opts...)
}

func (r *rpcClient) NewRequest(service, method string, request interface{}, reqOpts ...RequestOption) Request {
	return newRequest(service, method, request, r.opts.ContentType, reqOpts...)
}

func (r *rpcClient) String() string {
	return "mucp"
}
