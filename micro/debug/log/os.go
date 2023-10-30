package log

import (
	"sync"

	"github.com/google/uuid"
	"github.com/wxc/micro/util/ring"
)

type osLog struct {
	format FormatFunc
	buffer *ring.Buffer
	subs   map[string]*osStream

	sync.RWMutex
	once sync.Once
}

type osStream struct {
	stream chan Record
}

func (o *osLog) Read(...ReadOption) ([]Record, error) {
	var records []Record

	// read the last 100 records
	for _, v := range o.buffer.Get(100) {
		records = append(records, v.Value.(Record))
	}

	return records, nil
}

func (o *osLog) Write(r Record) error {
	o.buffer.Put(r)
	return nil
}

func (o *osLog) Stream() (Stream, error) {
	o.Lock()
	defer o.Unlock()

	st := &osStream{
		stream: make(chan Record, 128),
	}

	o.subs[uuid.New().String()] = st
	return st, nil
}

func (o *osStream) Chan() <-chan Record {
	return o.stream
}

func (o *osStream) Stop() error {
	return nil
}

func NewLog(opts ...Option) Log {
	options := Options{
		Format: DefaultFormat,
	}
	for _, o := range opts {
		o(&options)
	}

	l := &osLog{
		format: options.Format,
		buffer: ring.New(1024),
		subs:   make(map[string]*osStream),
	}

	return l
}
