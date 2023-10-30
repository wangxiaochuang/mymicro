package config

import (
	"sync"

	"github.com/wxc/micro/config/loader"
	"github.com/wxc/micro/config/reader"
	"github.com/wxc/micro/config/reader/json"
	"github.com/wxc/micro/config/source"
	"go-micro.dev/v4/config/loader/memory"
)

type config struct {
	vals reader.Values
	exit chan bool
	snap *loader.Snapshot
	opts Options
	sync.RWMutex
}

type watcher struct {
	lw    loader.Watcher
	rd    reader.Reader
	value reader.Value
	path  []string
}

func newConfig(opts ...Option) (Config, error) {
	var c config
	err := c.Init(opts...)
	if err != nil {
		return nil, err
	}
	if !c.opts.WithWatcherDisabled {
		go c.run()
	}
	return &c, nil
}

func (c *config) Init(opts ...Option) error {
	c.opts = Options{
		Reader: json.NewReader(),
	}
	c.exit = make(chan bool)
	for _, o := range opts {
		o(&c.opts)
	}

	if c.opts.Loader == nil {
		loaderOpts := []loader.Option{memory.WithReader(c.opts.Reader)}
		if c.opts.WithWatcherDisabled {
			loaderOpts = append(loaderOpts, memory.WithWatcherDisabled())
		}

		c.opts.Loader = memory.NewLoader(loaderOpts...)
	}
	err := c.opts.Loader.Load(c.opts.Source...)
	if err != nil {
		return err
	}

	c.snap, err = c.opts.Loader.Snapshot()
	if err != nil {
		return err
	}

	c.vals, err = c.opts.Reader.Values(c.snap.ChangeSet)
	if err != nil {
		return err
	}

	return nil
}

func (c *config) Options() Options {
	return c.opts
}

func (c *config) run() {
	panic("in run")
}

func (c *config) Map() map[string]interface{} {
	c.RLock()
	defer c.RUnlock()
	return c.vals.Map()
}

func (c *config) Scan(v interface{}) error {
	c.RLock()
	defer c.RUnlock()
	return c.vals.Scan(v)
}

func (c *config) Sync() error {
	panic("in sync")
}

func (c *config) Close() error {
	select {
	case <-c.exit:
		return nil
	default:
		close(c.exit)
	}
	return nil
}

func (c *config) Get(path ...string) reader.Value {
	c.RLock()
	defer c.RUnlock()

	if c.vals != nil {
		return c.vals.Get(path...)
	}

	return newValue()
}

func (c *config) Set(val interface{}, path ...string) {
	c.Lock()
	defer c.Unlock()

	if c.vals != nil {
		c.vals.Set(val, path...)
	}

	return
}

func (c *config) Del(path ...string) {
	c.Lock()
	defer c.Unlock()

	if c.vals != nil {
		c.vals.Del(path...)
	}

	return
}

func (c *config) Bytes() []byte {
	c.RLock()
	defer c.RUnlock()

	if c.vals == nil {
		return []byte{}
	}

	return c.vals.Bytes()
}

func (c *config) Load(sources ...source.Source) error {
	panic("in Load")
}

func (c *config) Watch(path ...string) (Watcher, error) {
	panic("in Watch")
}

func (c *config) String() string {
	return "config"
}

func (w *watcher) Next() (reader.Value, error) {
	panic("in Next")
}

func (w *watcher) Stop() error {
	return w.lw.Stop()
}
