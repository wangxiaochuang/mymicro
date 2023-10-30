package ring

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type Buffer struct {
	streams map[string]*Stream
	vals    []*Entry
	size    int

	sync.RWMutex
}

type Entry struct {
	Value     interface{}
	Timestamp time.Time
}

type Stream struct {
	Entries chan *Entry
	Stop    chan bool
	Id      string
}

func (b *Buffer) Put(v interface{}) {
	b.Lock()
	defer b.Unlock()

	entry := &Entry{
		Value:     v,
		Timestamp: time.Now(),
	}
	b.vals = append(b.vals, entry)

	if len(b.vals) > b.size {
		b.vals = b.vals[1:]
	}

	for _, stream := range b.streams {
		select {
		case <-stream.Stop:
			delete(b.streams, stream.Id)
			close(stream.Entries)
		case stream.Entries <- entry:
		}
	}
}

func (b *Buffer) Get(n int) []*Entry {
	b.RLock()
	defer b.RUnlock()
	if n > len(b.vals) || n < 0 {
		n = len(b.vals)
	}

	delta := len(b.vals) - n
	return b.vals[delta:]
}

func (b *Buffer) Since(t time.Time) []*Entry {
	b.RLock()
	defer b.RUnlock()

	if t.IsZero() {
		return b.vals
	}

	if time.Since(t).Seconds() < 0.0 {
		return nil
	}

	for i, v := range b.vals {
		d := v.Timestamp.Sub(t)

		if d.Seconds() > 0.0 {
			return b.vals[i:]
		}
	}

	return nil
}

func (b *Buffer) Stream() (<-chan *Entry, chan bool) {
	b.Lock()
	defer b.Unlock()

	entries := make(chan *Entry, 128)
	id := uuid.New().String()
	stop := make(chan bool)

	b.streams[id] = &Stream{
		Id:      id,
		Entries: entries,
		Stop:    stop,
	}

	return entries, stop
}

func (b *Buffer) Size() int {
	return b.size
}

func New(i int) *Buffer {
	return &Buffer{
		size:    i,
		streams: make(map[string]*Stream),
	}
}
