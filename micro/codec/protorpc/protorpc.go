package protorpc

import (
	"bytes"
	"io"
	"strconv"
	"sync"

	"github.com/wxc/micro/codec"
)

type flusher interface {
	Flush() error
}

type protoCodec struct {
	rwc io.ReadWriteCloser
	buf *bytes.Buffer
	mt  codec.MessageType
	sync.Mutex
}

func (c *protoCodec) Close() error {
	c.buf.Reset()
	return c.rwc.Close()
}

func (c *protoCodec) String() string {
	return "proto-rpc"
}

func id(id string) uint64 {
	p, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		p = 0
	}
	i := uint64(p)
	return i
}

func (c *protoCodec) Write(m *codec.Message, b interface{}) error {
	panic("in Write")
}

func (c *protoCodec) ReadHeader(m *codec.Message, mt codec.MessageType) error {
	panic("in ReadHeader")
}

func (c *protoCodec) ReadBody(b interface{}) error {
	panic("in ReadBody")
}

func NewCodec(rwc io.ReadWriteCloser) codec.Codec {
	return &protoCodec{
		buf: bytes.NewBuffer(nil),
		rwc: rwc,
	}
}
