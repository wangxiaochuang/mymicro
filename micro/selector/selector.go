package selector

import (
	"errors"

	"github.com/wxc/micro/registry"
)

type Selector interface {
	Init(opts ...Option) error
	Options() Options
	Select(service string, opts ...SelectOption) (Next, error)
	Mark(service string, node *registry.Node, err error)
	Reset(service string)
	Close() error
	String() string
}

type Next func() (*registry.Node, error)

type Filter func([]*registry.Service) []*registry.Service

type Strategy func([]*registry.Service) Next

var (
	DefaultSelector = NewSelector()

	ErrNotFound      = errors.New("not found")
	ErrNoneAvailable = errors.New("none available")
)
