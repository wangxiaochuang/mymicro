package codec

import (
	"errors"
	"io"
)

const (
	Error MessageType = iota
	Request
	Response
	Event
)

var (
	ErrInvalidMessage = errors.New("invalid message")
)

type MessageType int

type NewCodec func(io.ReadWriteCloser) Codec

type Codec interface {
	Reader
	Writer
	Close() error
	String() string
}

type Reader interface {
	ReadHeader(*Message, MessageType) error
	ReadBody(interface{}) error
}

type Writer interface {
	Write(*Message, interface{}) error
}

type Marshaler interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
	String() string
}

type Message struct {
	Header   map[string]string
	Id       string
	Target   string
	Method   string
	Endpoint string
	Error    string

	Body []byte
	Type MessageType
}
