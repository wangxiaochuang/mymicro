package grpc

import (
	"io"

	"github.com/wxc/micro/codec"
)

type Codec struct {
	Conn        io.ReadWriteCloser
	ContentType string
}

func (c *Codec) ReadHeader(m *codec.Message, t codec.MessageType) error {
	panic("ReadHeader")
}

func (c *Codec) ReadBody(b interface{}) error {
	panic("in ReadBody")
}

func (c *Codec) Write(m *codec.Message, b interface{}) error {
	panic("in Write")
}

func (c *Codec) Close() error {
	return c.Conn.Close()
}

func (c *Codec) String() string {
	return "grpc"
}

func NewCodec(c io.ReadWriteCloser) codec.Codec {
	return &Codec{
		Conn:        c,
		ContentType: "application/grpc",
	}
}
