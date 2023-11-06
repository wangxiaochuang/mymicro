package auth

import (
	"context"
	"errors"
	"time"
)

const (
	BearerScheme = "Bearer "
	ScopePublic  = ""
	ScopeAccount = "*"
)

var (
	ErrInvalidToken = errors.New("invalid token provided")
	ErrForbidden    = errors.New("resource forbidden")
)

type Auth interface {
	// Init the auth
	Init(opts ...Option)
	Options() Options
	Generate(id string, opts ...GenerateOption) (*Account, error)
	Inspect(token string) (*Account, error)
	Token(opts ...TokenOption) (*Token, error)
	String() string
}

type Rules interface {
	Verify(acc *Account, res *Resource, opts ...VerifyOption) error
	Grant(rule *Rule) error
	Revoke(rule *Rule) error
	List(...ListOption) ([]*Rule, error)
}

type Account struct {
	ID       string            `json:"id"`
	Type     string            `json:"type"`
	Issuer   string            `json:"issuer"`
	Metadata map[string]string `json:"metadata"`
	Scopes   []string          `json:"scopes"`
	Secret   string            `json:"secret"`
}

type Token struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	Created      time.Time `json:"created"`
	Expiry       time.Time `json:"expiry"`
}

func (t *Token) Expired() bool {
	return t.Expiry.Unix() < time.Now().Unix()
}

type Resource struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Endpoint string `json:"endpoint"`
}

type Access int

const (
	AccessGranted Access = iota
	AccessDenied
)

type Rule struct {
	ID       string
	Scope    string
	Resource *Resource
	Access   Access
	Priority int32
}

type accountKey struct{}

func AccountFromContext(ctx context.Context) (*Account, bool) {
	acc, ok := ctx.Value(accountKey{}).(*Account)
	return acc, ok
}

func ContextWithAccount(ctx context.Context, account *Account) context.Context {
	return context.WithValue(ctx, accountKey{}, account)
}
