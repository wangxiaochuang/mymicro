package memory

import (
	"container/list"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/wxc/micro/config/loader"
	"github.com/wxc/micro/config/reader"
	"github.com/wxc/micro/config/source"
)

type memory struct {
	vals     reader.Values
	exit     chan bool
	snap     *loader.Snapshot
	watchers *list.List
	opts     loader.Options
	sets     []*source.ChangeSet
	sources  []source.Source
	sync.RWMutex
}

type updateValue struct {
	value   reader.Value
	version string
}

type watcher struct {
	value   reader.Value
	reader  reader.Reader
	version atomic.Value
	exit    chan bool
	updates chan updateValue
	path    []string
}

func (w *watcher) getVersion() string {
	return w.version.Load().(string)
}

func (m *memory) watch(idx int, s source.Source) {
	panic("in watch")
}

func (m *memory) loaded() bool {
	var loaded bool
	m.RLock()
	if m.vals != nil {
		loaded = true
	}
	m.RUnlock()
	return loaded
}

func (m *memory) reload() error {
	panic("in reload")
}

func (m *memory) update() {
	panic("in update")
}

func (m *memory) Snapshot() (*loader.Snapshot, error) {
	panic("in Snapshot")
}

func (m *memory) Sync() error {
	panic("in Sync")
}

func (m *memory) Close() error {
	select {
	case <-m.exit:
		return nil
	default:
		close(m.exit)
	}
	return nil
}

func (m *memory) Get(path ...string) (reader.Value, error) {
	panic("in Get")
}

func (m *memory) Load(sources ...source.Source) error {
	panic(" in Load")
}

func (m *memory) Watch(path ...string) (loader.Watcher, error) {
	panic(" in Watch")
}

func (m *memory) String() string {
	return "memory"
}

func (w *watcher) Next() (*loader.Snapshot, error) {
	panic(" in Next")
}

func (w *watcher) Stop() error {
	select {
	case <-w.exit:
	default:
		close(w.exit)
		close(w.updates)
	}

	return nil
}

func genVer() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func NewLoader(opts ...loader.Option) loader.Loader {
	panic("in NewLoader")
}
