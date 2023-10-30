package client

import (
	"context"

	"github.com/wxc/micro/codec"
)

var (
	NewClient     func(...Option) Client = newRPCClient
	DefaultClient Client                 = newRPCClient()
)

type Client interface {
	Init(...Option) error
	Options() Options
	NewMessage(topic string, msg interface{}, opts ...MessageOption) Message
	NewRequest(service, endpoint string, req interface{}, reqOpts ...RequestOption) Request
	Call(ctx context.Context, req Request, rsp interface{}, opts ...CallOption) error
	Stream(ctx context.Context, req Request, opts ...CallOption) (Stream, error)
	Publish(ctx context.Context, msg Message, opts ...PublishOption) error
	String() string
}

type Router interface {
	SendRequest(context.Context, Request) (Response, error)
}

type Message interface {
	Topic() string
	Payload() interface{}
	ContentType() string
}

type Request interface {
	Service() string
	Method() string
	Endpoint() string
	ContentType() string
	Body() interface{}
	Codec() codec.Writer
	Stream() bool
}

type Response interface {
	Codec() codec.Reader
	Header() map[string]string
	Read() ([]byte, error)
}

type Stream interface {
	Closer
	Content() context.Context
	Request() Request
	Response() Response
	Send(interface{}) error
	Recv(interface{}) error
	Error() error
	Close() error
}

type Closer interface {
	CloseSend() error
}

type Option func(*Options)

type CallOption func(*CallOptions)

type PublishOption func(*PublishOptions)

type MessageOption func(*MessageOptions)

type RequestOption func(*RequestOptions)

func Call(ctx context.Context, request Request, response interface{}, opts ...CallOption) error {
	return DefaultClient.Call(ctx, request, response, opts...)
}

func Publish(ctx context.Context, msg Message, opts ...PublishOption) error {
	return DefaultClient.Publish(ctx, msg, opts...)
}

func NewMessage(topic string, payload interface{}, opts ...MessageOption) Message {
	return DefaultClient.NewMessage(topic, payload, opts...)
}

func NewRequest(service, endpoint string, request interface{}, reqOpts ...RequestOption) Request {
	return DefaultClient.NewRequest(service, endpoint, request, reqOpts...)
}

func NewStream(ctx context.Context, request Request, opts ...CallOption) (Stream, error) {
	return DefaultClient.Stream(ctx, request, opts...)
}

func String() string {
	return DefaultClient.String()
}
