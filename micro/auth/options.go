package auth

import (
	"context"
	"time"

	"github.com/wxc/micro/logger"
)

func NewOptions(opts ...Option) Options {
	options := Options{
		Logger: logger.DefaultLogger,
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}

type Options struct {
	Namespace  string
	ID         string
	Secret     string
	Token      *Token
	PublicKey  string
	PrivateKey string
	Addrs      []string
	Logger     logger.Logger
}

type Option func(o *Options)

func Addrs(addrs ...string) Option {
	return func(o *Options) {
		o.Addrs = addrs
	}
}

func Namespace(n string) Option {
	return func(o *Options) {
		o.Namespace = n
	}
}

func PublicKey(key string) Option {
	return func(o *Options) {
		o.PublicKey = key
	}
}

func PrivateKey(key string) Option {
	return func(o *Options) {
		o.PrivateKey = key
	}
}

func WithLogger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

func Credentials(id, secret string) Option {
	return func(o *Options) {
		o.ID = id
		o.Secret = secret
	}
}

func ClientToken(token *Token) Option {
	return func(o *Options) {
		o.Token = token
	}
}

type GenerateOptions struct {
	Metadata map[string]string
	Scopes   []string
	Provider string
	Type     string
	Secret   string
}

type GenerateOption func(o *GenerateOptions)

func WithSecret(s string) GenerateOption {
	return func(o *GenerateOptions) {
		o.Secret = s
	}
}

func WithType(t string) GenerateOption {
	return func(o *GenerateOptions) {
		o.Type = t
	}
}

func WithMetadata(md map[string]string) GenerateOption {
	return func(o *GenerateOptions) {
		o.Metadata = md
	}
}

func WithProvider(p string) GenerateOption {
	return func(o *GenerateOptions) {
		o.Provider = p
	}
}

func WithScopes(s ...string) GenerateOption {
	return func(o *GenerateOptions) {
		o.Scopes = s
	}
}

func NewGenerateOptions(opts ...GenerateOption) GenerateOptions {
	var options GenerateOptions
	for _, o := range opts {
		o(&options)
	}
	return options
}

type TokenOptions struct {
	ID           string
	Secret       string
	RefreshToken string
	Expiry       time.Duration
}

type TokenOption func(o *TokenOptions)

func WithExpiry(ex time.Duration) TokenOption {
	return func(o *TokenOptions) {
		o.Expiry = ex
	}
}

func WithCredentials(id, secret string) TokenOption {
	return func(o *TokenOptions) {
		o.ID = id
		o.Secret = secret
	}
}

func WithToken(rt string) TokenOption {
	return func(o *TokenOptions) {
		o.RefreshToken = rt
	}
}

func NewTokenOptions(opts ...TokenOption) TokenOptions {
	var options TokenOptions
	for _, o := range opts {
		o(&options)
	}

	// set defualt expiry of token
	if options.Expiry == 0 {
		options.Expiry = time.Minute
	}

	return options
}

type VerifyOptions struct {
	Context context.Context
}

type VerifyOption func(o *VerifyOptions)

func VerifyContext(ctx context.Context) VerifyOption {
	return func(o *VerifyOptions) {
		o.Context = ctx
	}
}

type ListOptions struct {
	Context context.Context
}

type ListOption func(o *ListOptions)

func RulesContext(ctx context.Context) ListOption {
	return func(o *ListOptions) {
		o.Context = ctx
	}
}
