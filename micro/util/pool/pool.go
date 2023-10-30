package pool

import (
	"time"

	"github.com/wxc/micro/transport"
)

type Pool interface {
	Close() error
	Get(addr string, opts ...transport.DialOption) (Conn, error)
	Release(c Conn, status error) error
}

type Conn interface {
	Id() string
	Created() time.Time
	transport.Client
}

func NewPool(opts ...Option) Pool {
	var options Options
	for _, o := range opts {
		o(&options)
	}

	return newPool(options)
}
