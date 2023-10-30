package registry

import "errors"

var (
	DefaultRegistry   = NewMemoryRegistry()
	ErrNotFound       = errors.New("service not found")
	ErrWatcherStopped = errors.New("watcher stopped")
)

type Registry interface {
	Init(...Option) error
	Options() Options
	Register(*Service, ...RegisterOption) error
	Deregister(*Service, ...DeregisterOption) error
	GetService(string, ...GetOption) ([]*Service, error)
	ListServices(...ListOption) ([]*Service, error)
	Watch(...WatchOption) (Watcher, error)
	String() string
}

type Service struct {
	Name      string            `json:"name"`
	Version   string            `json:"version"`
	Metadata  map[string]string `json:"metadata"`
	Endpoints []*Endpoint       `json:"endpoints"`
	Nodes     []*Node           `json:"nodes"`
}

type Node struct {
	Metadata map[string]string `json:"metadata"`
	Id       string            `json:"id"`
	Address  string            `json:"address"`
}

type Endpoint struct {
	Request  *Value            `json:"request"`
	Response *Value            `json:"response"`
	Metadata map[string]string `json:"metadata"`
	Name     string            `json:"name"`
}

type Value struct {
	Name   string   `json:"name"`
	Type   string   `json:"type"`
	Values []*Value `json:"values"`
}

type Option func(*Options)

type RegisterOption func(*RegisterOptions)

type WatchOption func(*WatchOptions)

type DeregisterOption func(*DeregisterOptions)

type GetOption func(*GetOptions)

type ListOption func(*ListOptions)

func Register(s *Service, opts ...RegisterOption) error {
	return DefaultRegistry.Register(s, opts...)
}

func Deregister(s *Service) error {
	return DefaultRegistry.Deregister(s)
}

func GetService(name string) ([]*Service, error) {
	return DefaultRegistry.GetService(name)
}

func ListServices() ([]*Service, error) {
	return DefaultRegistry.ListServices()
}

func Watch(opts ...WatchOption) (Watcher, error) {
	return DefaultRegistry.Watch(opts...)
}

func String() string {
	return DefaultRegistry.String()
}
