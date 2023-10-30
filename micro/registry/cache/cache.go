package cache

import (
	"math"
	"math/rand"
	"sync"
	"time"

	log "github.com/wxc/micro/logger"
	"github.com/wxc/micro/registry"
	"golang.org/x/sync/singleflight"
)

type Cache interface {
	// embed the registry interface
	registry.Registry
	// stop the cache watcher
	Stop()
}

type Options struct {
	Logger log.Logger
	TTL    time.Duration
}

type Option func(o *Options)

type cache struct {
	opts Options
	registry.Registry
	status         error
	sg             singleflight.Group
	cache          map[string][]*registry.Service
	ttls           map[string]time.Time
	watched        map[string]bool
	exit           chan bool
	watchedRunning map[string]bool
	sync.RWMutex
}

var (
	DefaultTTL = time.Minute
)

func backoff(attempts int) time.Duration {
	if attempts == 0 {
		return time.Duration(0)
	}
	return time.Duration(math.Pow(10, float64(attempts))) * time.Millisecond
}

func (c *cache) getStatus() error {
	c.RLock()
	defer c.RUnlock()
	return c.status
}

func (c *cache) setStatus(err error) {
	c.Lock()
	c.status = err
	c.Unlock()
}

func (c *cache) isValid(services []*registry.Service, ttl time.Time) bool {
	// no services exist
	if len(services) == 0 {
		return false
	}

	// ttl is invalid
	if ttl.IsZero() {
		return false
	}

	// time since ttl is longer than timeout
	if time.Since(ttl) > 0 {
		return false
	}

	// ok
	return true
}

func (c *cache) quit() bool {
	select {
	case <-c.exit:
		return true
	default:
		return false
	}
}

func (c *cache) del(service string) {
	// don't blow away cache in error state
	if err := c.status; err != nil {
		return
	}
	delete(c.cache, service)
	delete(c.ttls, service)
}

func (c *cache) get(service string) ([]*registry.Service, error) {
	panic("in get")
}

func (c *cache) set(service string, services []*registry.Service) {
	c.cache[service] = services
	c.ttls[service] = time.Now().Add(c.opts.TTL)
}

func (c *cache) update(res *registry.Result) {
	panic("in update")
}

func (c *cache) run(service string) {
	panic("in run")
}

func (c *cache) watch(w registry.Watcher) error {
	panic(" in watch")
}

func (c *cache) GetService(service string, opts ...registry.GetOption) ([]*registry.Service, error) {
	panic(" in GetService")
}

func (c *cache) Stop() {
	c.Lock()
	defer c.Unlock()

	select {
	case <-c.exit:
		return
	default:
		close(c.exit)
	}
}

func (c *cache) String() string {
	return "cache"
}

func New(r registry.Registry, opts ...Option) Cache {
	rand.Seed(time.Now().UnixNano())

	options := Options{
		TTL:    DefaultTTL,
		Logger: log.DefaultLogger,
	}

	for _, o := range opts {
		o(&options)
	}

	return &cache{
		Registry:       r,
		opts:           options,
		watched:        make(map[string]bool),
		watchedRunning: make(map[string]bool),
		cache:          make(map[string][]*registry.Service),
		ttls:           make(map[string]time.Time),
		exit:           make(chan bool),
	}
}
