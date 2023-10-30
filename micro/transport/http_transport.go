package transport

import "github.com/wxc/micro/logger"

type httpTransport struct {
	opts Options
}

func NewHTTPTransport(opts ...Option) *httpTransport {
	options := Options{
		BuffSizeH2: DefaultBufSizeH2,
		Logger:     logger.DefaultLogger,
	}

	for _, o := range opts {
		o(&options)
	}

	return &httpTransport{opts: options}
}

func (h *httpTransport) Init(opts ...Option) error {
	for _, o := range opts {
		o(&h.opts)
	}

	return nil
}

func (h *httpTransport) Dial(addr string, opts ...DialOption) (Client, error) {
	panic(" in Dial")
}

func (h *httpTransport) Listen(addr string, opts ...ListenOption) (Listener, error) {
	panic(" in Listen")
}

func (h *httpTransport) Options() Options {
	return h.opts
}

func (h *httpTransport) String() string {
	return "http"
}
