package client

import (
	"bytes"
	errs "errors"

	"github.com/wxc/micro/codec"
	raw "github.com/wxc/micro/codec/bytes"
	"github.com/wxc/micro/codec/grpc"
	"github.com/wxc/micro/codec/json"
	"github.com/wxc/micro/codec/jsonrpc"
	"github.com/wxc/micro/codec/proto"
	"github.com/wxc/micro/codec/protorpc"
	"github.com/wxc/micro/errors"
	"github.com/wxc/micro/registry"
	"github.com/wxc/micro/transport"
	"github.com/wxc/micro/transport/headers"
)

const (
	lastStreamResponseError = "EOS"
)

type serverError string

func (e serverError) Error() string {
	return string(e)
}

var (
	errShutdown = errs.New("connection is shut down")
)

type rpcCodec struct {
	client transport.Client
	codec  codec.Codec

	req *transport.Message
	buf *readWriteCloser

	// signify if its a stream
	stream string
}

type readWriteCloser struct {
	wbuf *bytes.Buffer
	rbuf *bytes.Buffer
}

var (
	// DefaultContentType header.
	DefaultContentType = "application/json"

	// DefaultCodecs map.
	DefaultCodecs = map[string]codec.NewCodec{
		"application/grpc":         grpc.NewCodec,
		"application/grpc+json":    grpc.NewCodec,
		"application/grpc+proto":   grpc.NewCodec,
		"application/protobuf":     proto.NewCodec,
		"application/json":         json.NewCodec,
		"application/json-rpc":     jsonrpc.NewCodec,
		"application/proto-rpc":    protorpc.NewCodec,
		"application/octet-stream": raw.NewCodec,
	}

	// TODO: remove legacy codec list.
	defaultCodecs = map[string]codec.NewCodec{
		"application/json":         jsonrpc.NewCodec,
		"application/json-rpc":     jsonrpc.NewCodec,
		"application/protobuf":     protorpc.NewCodec,
		"application/proto-rpc":    protorpc.NewCodec,
		"application/octet-stream": protorpc.NewCodec,
	}
)

func (rwc *readWriteCloser) Read(p []byte) (n int, err error) {
	return rwc.rbuf.Read(p)
}

func (rwc *readWriteCloser) Write(p []byte) (n int, err error) {
	return rwc.wbuf.Write(p)
}

func (rwc *readWriteCloser) Close() error {
	rwc.rbuf.Reset()
	rwc.wbuf.Reset()

	return nil
}

func getHeaders(m *codec.Message) {
	set := func(v, hdr string) string {
		if len(v) > 0 {
			return v
		}

		return m.Header[hdr]
	}

	m.Error = set(m.Error, headers.Error)

	m.Endpoint = set(m.Endpoint, headers.Endpoint)

	m.Method = set(m.Method, headers.Method)

	m.Id = set(m.Id, headers.ID)
}

func setHeaders(m *codec.Message, stream string) {
	set := func(hdr, v string) {
		if len(v) == 0 {
			return
		}

		m.Header[hdr] = v
	}

	set(headers.ID, m.Id)
	set(headers.Request, m.Target)
	set(headers.Method, m.Method)
	set(headers.Endpoint, m.Endpoint)
	set(headers.Error, m.Error)

	if len(stream) > 0 {
		set(headers.Stream, stream)
	}
}

func setupProtocol(msg *transport.Message, node *registry.Node) codec.NewCodec {
	protocol := node.Metadata["protocol"]

	// got protocol
	if len(protocol) > 0 {
		return nil
	}

	// processing topic publishing
	if len(msg.Header[headers.Message]) > 0 {
		return nil
	}

	// no protocol use old codecs
	switch msg.Header["Content-Type"] {
	case "application/json":
		msg.Header["Content-Type"] = "application/json-rpc"
	case "application/protobuf":
		msg.Header["Content-Type"] = "application/proto-rpc"
	}

	return defaultCodecs[msg.Header["Content-Type"]]
}

func newRPCCodec(req *transport.Message, client transport.Client, c codec.NewCodec, stream string) codec.Codec {
	rwc := &readWriteCloser{
		wbuf: bytes.NewBuffer(nil),
		rbuf: bytes.NewBuffer(nil),
	}

	return &rpcCodec{
		buf:    rwc,
		client: client,
		codec:  c(rwc),
		req:    req,
		stream: stream,
	}
}

func (c *rpcCodec) Write(message *codec.Message, body interface{}) error {
	panic("in Write")
}

func (c *rpcCodec) ReadHeader(msg *codec.Message, r codec.MessageType) error {
	panic("in ReadHeader")
}

func (c *rpcCodec) ReadBody(b interface{}) error {
	panic("in ReadBody")
}

func (c *rpcCodec) Close() error {
	if err := c.buf.Close(); err != nil {
		return err
	}

	if err := c.codec.Close(); err != nil {
		return err
	}

	if err := c.client.Close(); err != nil {
		return errors.InternalServerError("go.micro.client.transport", err.Error())
	}

	return nil
}

func (c *rpcCodec) String() string {
	return "rpc"
}
