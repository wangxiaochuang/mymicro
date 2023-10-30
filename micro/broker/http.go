package broker

import (
	"crypto/tls"
	"errors"
	"io"
	"math/rand"
	"net"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/wxc/micro/codec/json"
	merr "github.com/wxc/micro/errors"
	"github.com/wxc/micro/registry"
	"github.com/wxc/micro/registry/cache"
	"go-micro.dev/v4/transport/headers"
	"golang.org/x/net/http2"
)

type httpBroker struct {
	opts Options

	r registry.Registry

	mux *http.ServeMux

	c           *http.Client
	subscribers map[string][]*httpSubscriber
	exit        chan chan error

	inbox   map[string][][]byte
	id      string
	address string

	sync.RWMutex

	// offline message inbox
	mtx     sync.RWMutex
	running bool
}

type httpSubscriber struct {
	opts  SubscribeOptions
	fn    Handler
	svc   *registry.Service
	hb    *httpBroker
	id    string
	topic string
}

type httpEvent struct {
	err error
	m   *Message
	t   string
}

var (
	DefaultPath      = "/"
	DefaultAddress   = "127.0.0.1:0"
	serviceName      = "micro.http.broker"
	broadcastVersion = "ff.http.broadcast"
	registerTTL      = time.Minute
	registerInterval = time.Second * 30
)

func init() {
	rand.Seed(time.Now().Unix())
}

func newTransport(config *tls.Config) *http.Transport {
	if config == nil {
		config = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	dialTLS := func(network string, addr string) (net.Conn, error) {
		return tls.Dial(network, addr, config)
	}

	t := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		DialTLS:             dialTLS,
	}
	runtime.SetFinalizer(&t, func(tr **http.Transport) {
		(*tr).CloseIdleConnections()
	})

	http2.ConfigureTransport(t)
	return t
}

func newHttpBroker(opts ...Option) Broker {
	options := *NewOptions(opts...)

	options.Registry = registry.DefaultRegistry
	options.Codec = json.Marshaler{}

	for _, o := range opts {
		o(&options)
	}

	addr := DefaultAddress

	if len(options.Addrs) > 0 && len(options.Addrs[0]) > 0 {
		addr = options.Addrs[0]
	}

	h := &httpBroker{
		id:          uuid.New().String(),
		address:     addr,
		opts:        options,
		r:           options.Registry,
		c:           &http.Client{Transport: newTransport(options.TLSConfig)},
		subscribers: make(map[string][]*httpSubscriber),
		exit:        make(chan chan error),
		mux:         http.NewServeMux(),
		inbox:       make(map[string][][]byte),
	}
	h.mux.Handle(DefaultPath, h)

	if h.opts.Context != nil {
		handlers, ok := h.opts.Context.Value("http_handlers").(map[string]http.Handler)
		if ok {
			for pattern, handler := range handlers {
				h.mux.Handle(pattern, handler)
			}
		}
	}
	return h
}

func (h *httpEvent) Ack() error {
	return nil
}

func (h *httpEvent) Error() error {
	return h.err
}

func (h *httpEvent) Message() *Message {
	return h.m
}

func (h *httpEvent) Topic() string {
	return h.t
}

func (h *httpSubscriber) Options() SubscribeOptions {
	return h.opts
}

func (h *httpSubscriber) Topic() string {
	return h.topic
}

func (h *httpSubscriber) Unsubscribe() error {
	return h.hb.unsubscribe(h)
}

func (h *httpBroker) saveMessage(topic string, msg []byte) {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	c := h.inbox[topic]
	c = append(c, msg)
	if len(c) > 64 {
		c = c[:64]
	}

	h.inbox[topic] = c
}

func (h *httpBroker) getMessage(topic string, num int) [][]byte {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	c, ok := h.inbox[topic]
	if !ok {
		return nil
	}

	if len(c) >= num {
		msg := c[:num]
		h.inbox[topic] = c[num:]
		return msg
	}

	h.inbox[topic] = nil
	return c
}

func (h *httpBroker) subscribe(s *httpSubscriber) error {
	h.Lock()
	defer h.Unlock()

	if err := h.r.Register(s.svc, registry.RegisterTTL(registerTTL)); err != nil {
		return err
	}

	h.subscribers[s.topic] = append(h.subscribers[s.topic], s)
	return nil
}

func (h *httpBroker) unsubscribe(s *httpSubscriber) error {
	h.Lock()
	defer h.Unlock()

	//nolint:prealloc
	var subscribers []*httpSubscriber

	for _, sub := range h.subscribers[s.topic] {
		if sub == s {
			_ = h.r.Deregister(sub.svc)
			continue
		}
		subscribers = append(subscribers, sub)
	}
	h.subscribers[s.topic] = subscribers

	return nil
}

func (h *httpBroker) run(l net.Listener) {
	t := time.NewTicker(registerInterval)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			h.RLock()
			for _, subs := range h.subscribers {
				for _, sub := range subs {
					_ = h.r.Register(sub.svc, registry.RegisterTTL(registerTTL))
				}
			}
			h.RUnlock()
		case ch := <-h.exit:
			ch <- l.Close()
			h.RLock()
			for _, subs := range h.subscribers {
				for _, sub := range subs {
					_ = h.r.Deregister(sub.svc)
				}
			}
			h.RUnlock()
			return
		}
	}
}

func (h *httpBroker) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		err := merr.BadRequest("go.micro.broker", "Method not allowed")
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}
	defer req.Body.Close()

	req.ParseForm()

	b, err := io.ReadAll(req.Body)
	if err != nil {
		errr := merr.InternalServerError("go.micro.broker", "Error reading request body: %v", err)
		w.WriteHeader(500)
		w.Write([]byte(errr.Error()))
		return
	}

	var m *Message
	if err = h.opts.Codec.Unmarshal(b, &m); err != nil {
		errr := merr.InternalServerError("go.micro.broker", "Error parsing request body: %v", err)
		w.WriteHeader(500)
		w.Write([]byte(errr.Error()))
		return
	}

	topic := m.Header[headers.Message]
	if len(topic) == 0 {
		errr := merr.InternalServerError("go.micro.broker", "Topic not found")
		w.WriteHeader(500)
		w.Write([]byte(errr.Error()))
		return
	}

	p := &httpEvent{m: m, t: topic}
	id := req.Form.Get("id")
	var subs []Handler

	h.RLock()
	for _, subscriber := range h.subscribers[topic] {
		if id != subscriber.id {
			continue
		}
		subs = append(subs, subscriber.fn)
	}
	h.RUnlock()

	for _, fn := range subs {
		p.err = fn(p)
	}
}

func (h *httpBroker) Address() string {
	h.RLock()
	defer h.RUnlock()
	return h.address
}

func (h *httpBroker) Connect() error {
	panic("in Connect")
}

func (h *httpBroker) Disconnect() error {
	panic("in Disconnect")
}

func (h *httpBroker) Init(opts ...Option) error {
	h.RLock()
	if h.running {
		h.RUnlock()
		return errors.New("cannot init while connected")
	}
	h.RUnlock()

	h.Lock()
	defer h.Unlock()

	for _, o := range opts {
		o(&h.opts)
	}

	if len(h.opts.Addrs) > 0 && len(h.opts.Addrs[0]) > 0 {
		h.address = h.opts.Addrs[0]
	}

	if len(h.id) == 0 {
		h.id = "go.micro.http.broker-" + uuid.New().String()
	}

	reg := h.opts.Registry
	if reg == nil {
		reg = registry.DefaultRegistry
	}

	if rc, ok := h.r.(cache.Cache); ok {
		rc.Stop()
	}

	h.r = cache.New(reg)

	if c := h.opts.TLSConfig; c != nil {
		h.c = &http.Client{
			Transport: newTransport(c),
		}
	}

	return nil
}

func (h *httpBroker) Options() Options {
	return h.opts
}

func (h *httpBroker) Publish(topic string, msg *Message, opts ...PublishOption) error {
	panic("in Publish")
}

func (h *httpBroker) Subscribe(topic string, handler Handler, opts ...SubscribeOption) (Subscriber, error) {
	panic("in subscribe")
}

func (h *httpBroker) String() string {
	return "http"
}

func NewBroker(opts ...Option) Broker {
	return newHttpBroker(opts...)
}
