package registry

import "time"

type Watcher interface {
	Next() (*Result, error)
	Stop()
}

type Result struct {
	Service *Service
	Action  string
}

type EventType int

const (
	Create EventType = iota
	Delete
	Update
)

func (t EventType) String() string {
	switch t {
	case Create:
		return "create"
	case Delete:
		return "delete"
	case Update:
		return "update"
	default:
		return "unknown"
	}
}

type Event struct {
	Timestamp time.Time
	Service   *Service
	Id        string
	Type      EventType
}
