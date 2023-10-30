package pool

import (
	"sync"
	"time"

	"github.com/wxc/micro/transport"
)

type pool struct {
	tr transport.Transport

	conns map[string][]*poolConn
	size  int
	ttl   time.Duration

	sync.Mutex
}

type poolConn struct {
	created time.Time
	transport.Client
	id string
}

func newPool(options Options) *pool {
	return &pool{
		size:  options.Size,
		tr:    options.Transport,
		ttl:   options.TTL,
		conns: make(map[string][]*poolConn),
	}
}

func (p *pool) Close() error {
	p.Lock()
	defer p.Unlock()

	var err error

	for k, c := range p.conns {
		for _, conn := range c {
			if nerr := conn.Client.Close(); nerr != nil {
				err = nerr
			}
		}

		delete(p.conns, k)
	}

	return err
}

func (p *poolConn) Close() error {
	return nil
}

func (p *poolConn) Id() string {
	return p.id
}

func (p *poolConn) Created() time.Time {
	return p.created
}

func (p *pool) Get(addr string, opts ...transport.DialOption) (Conn, error) {
	panic("in Get")
}

func (p *pool) Release(conn Conn, err error) error {
	panic("in Release")
}
