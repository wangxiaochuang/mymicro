package cmd

import (
	"runtime/pprof"

	"github.com/urfave/cli/v2"
	"github.com/wxc/micro/auth"
	"github.com/wxc/micro/broker"
	"github.com/wxc/micro/client"
	"github.com/wxc/micro/config"
	"github.com/wxc/micro/debug/profile"
	"github.com/wxc/micro/debug/profile/http"
	"github.com/wxc/micro/debug/trace"
	"github.com/wxc/micro/registry"
	"github.com/wxc/micro/registry/cache"
	"github.com/wxc/micro/runtime"
	"github.com/wxc/micro/selector"
	"github.com/wxc/micro/server"
	"github.com/wxc/micro/store"
	"github.com/wxc/micro/transport"
)

var (
	DefaultFlags = []cli.Flag{
		&cli.StringFlag{
			Name:    "client",
			EnvVars: []string{"MICRO_CLIENT"},
			Usage:   "Client for go-micro; rpc",
		},
		&cli.StringFlag{
			Name:    "client_request_timeout",
			EnvVars: []string{"MICRO_CLIENT_REQUEST_TIMEOUT"},
			Usage:   "Sets the client request timeout. e.g 500ms, 5s, 1m. Default: 5s",
		},
		&cli.IntFlag{
			Name:    "client_retries",
			EnvVars: []string{"MICRO_CLIENT_RETRIES"},
			Value:   client.DefaultRetries,
			Usage:   "Sets the client retries. Default: 1",
		},
		&cli.IntFlag{
			Name:    "client_pool_size",
			EnvVars: []string{"MICRO_CLIENT_POOL_SIZE"},
			Usage:   "Sets the client connection pool size. Default: 1",
		},
		&cli.StringFlag{
			Name:    "client_pool_ttl",
			EnvVars: []string{"MICRO_CLIENT_POOL_TTL"},
			Usage:   "Sets the client connection pool ttl. e.g 500ms, 5s, 1m. Default: 1m",
		},
		&cli.IntFlag{
			Name:    "register_ttl",
			EnvVars: []string{"MICRO_REGISTER_TTL"},
			Value:   60,
			Usage:   "Register TTL in seconds",
		},
		&cli.IntFlag{
			Name:    "register_interval",
			EnvVars: []string{"MICRO_REGISTER_INTERVAL"},
			Value:   30,
			Usage:   "Register interval in seconds",
		},
		&cli.StringFlag{
			Name:    "server",
			EnvVars: []string{"MICRO_SERVER"},
			Usage:   "Server for go-micro; rpc",
		},
		&cli.StringFlag{
			Name:    "server_name",
			EnvVars: []string{"MICRO_SERVER_NAME"},
			Usage:   "Name of the server. go.micro.srv.example",
		},
		&cli.StringFlag{
			Name:    "server_version",
			EnvVars: []string{"MICRO_SERVER_VERSION"},
			Usage:   "Version of the server. 1.1.0",
		},
		&cli.StringFlag{
			Name:    "server_id",
			EnvVars: []string{"MICRO_SERVER_ID"},
			Usage:   "Id of the server. Auto-generated if not specified",
		},
		&cli.StringFlag{
			Name:    "server_address",
			EnvVars: []string{"MICRO_SERVER_ADDRESS"},
			Usage:   "Bind address for the server. 127.0.0.1:8080",
		},
		&cli.StringFlag{
			Name:    "server_advertise",
			EnvVars: []string{"MICRO_SERVER_ADVERTISE"},
			Usage:   "Used instead of the server_address when registering with discovery. 127.0.0.1:8080",
		},
		&cli.StringSliceFlag{
			Name:    "server_metadata",
			EnvVars: []string{"MICRO_SERVER_METADATA"},
			Value:   &cli.StringSlice{},
			Usage:   "A list of key-value pairs defining metadata. version=1.0.0",
		},
		&cli.StringFlag{
			Name:    "broker",
			EnvVars: []string{"MICRO_BROKER"},
			Usage:   "Broker for pub/sub. http, nats, rabbitmq",
		},
		&cli.StringFlag{
			Name:    "broker_address",
			EnvVars: []string{"MICRO_BROKER_ADDRESS"},
			Usage:   "Comma-separated list of broker addresses",
		},
		&cli.StringFlag{
			Name:    "profile",
			Usage:   "Debug profiler for cpu and memory stats",
			EnvVars: []string{"MICRO_DEBUG_PROFILE"},
		},
		&cli.StringFlag{
			Name:    "registry",
			EnvVars: []string{"MICRO_REGISTRY"},
			Usage:   "Registry for discovery. etcd, mdns",
		},
		&cli.StringFlag{
			Name:    "registry_address",
			EnvVars: []string{"MICRO_REGISTRY_ADDRESS"},
			Usage:   "Comma-separated list of registry addresses",
		},
		&cli.StringFlag{
			Name:    "runtime",
			Usage:   "Runtime for building and running services e.g local, kubernetes",
			EnvVars: []string{"MICRO_RUNTIME"},
		},
		&cli.StringFlag{
			Name:    "runtime_source",
			Usage:   "Runtime source for building and running services e.g github.com/micro/service",
			EnvVars: []string{"MICRO_RUNTIME_SOURCE"},
			Value:   "github.com/micro/services",
		},
		&cli.StringFlag{
			Name:    "selector",
			EnvVars: []string{"MICRO_SELECTOR"},
			Usage:   "Selector used to pick nodes for querying",
		},
		&cli.StringFlag{
			Name:    "store",
			EnvVars: []string{"MICRO_STORE"},
			Usage:   "Store used for key-value storage",
		},
		&cli.StringFlag{
			Name:    "store_address",
			EnvVars: []string{"MICRO_STORE_ADDRESS"},
			Usage:   "Comma-separated list of store addresses",
		},
		&cli.StringFlag{
			Name:    "store_database",
			EnvVars: []string{"MICRO_STORE_DATABASE"},
			Usage:   "Database option for the underlying store",
		},
		&cli.StringFlag{
			Name:    "store_table",
			EnvVars: []string{"MICRO_STORE_TABLE"},
			Usage:   "Table option for the underlying store",
		},
		&cli.StringFlag{
			Name:    "transport",
			EnvVars: []string{"MICRO_TRANSPORT"},
			Usage:   "Transport mechanism used; http",
		},
		&cli.StringFlag{
			Name:    "transport_address",
			EnvVars: []string{"MICRO_TRANSPORT_ADDRESS"},
			Usage:   "Comma-separated list of transport addresses",
		},
		&cli.StringFlag{
			Name:    "tracer",
			EnvVars: []string{"MICRO_TRACER"},
			Usage:   "Tracer for distributed tracing, e.g. memory, jaeger",
		},
		&cli.StringFlag{
			Name:    "tracer_address",
			EnvVars: []string{"MICRO_TRACER_ADDRESS"},
			Usage:   "Comma-separated list of tracer addresses",
		},
		&cli.StringFlag{
			Name:    "auth",
			EnvVars: []string{"MICRO_AUTH"},
			Usage:   "Auth for role based access control, e.g. service",
		},
		&cli.StringFlag{
			Name:    "auth_id",
			EnvVars: []string{"MICRO_AUTH_ID"},
			Usage:   "Account ID used for client authentication",
		},
		&cli.StringFlag{
			Name:    "auth_secret",
			EnvVars: []string{"MICRO_AUTH_SECRET"},
			Usage:   "Account secret used for client authentication",
		},
		&cli.StringFlag{
			Name:    "auth_namespace",
			EnvVars: []string{"MICRO_AUTH_NAMESPACE"},
			Usage:   "Namespace for the services auth account",
			Value:   "go.micro",
		},
		&cli.StringFlag{
			Name:    "auth_public_key",
			EnvVars: []string{"MICRO_AUTH_PUBLIC_KEY"},
			Usage:   "Public key for JWT auth (base64 encoded PEM)",
		},
		&cli.StringFlag{
			Name:    "auth_private_key",
			EnvVars: []string{"MICRO_AUTH_PRIVATE_KEY"},
			Usage:   "Private key for JWT auth (base64 encoded PEM)",
		},
		&cli.StringFlag{
			Name:    "config",
			EnvVars: []string{"MICRO_CONFIG"},
			Usage:   "The source of the config to be used to get configuration",
		},
	}

	DefaultBrokers = map[string]func(...broker.Option) broker.Broker{}

	DefaultClients = map[string]func(...client.Option) client.Client{}

	DefaultRegistries = map[string]func(...registry.Option) registry.Registry{}

	DefaultSelectors = map[string]func(...selector.Option) selector.Selector{}

	DefaultServers = map[string]func(...server.Option) server.Server{}

	DefaultTransports = map[string]func(...transport.Option) transport.Transport{}

	DefaultRuntimes = map[string]func(...runtime.Option) runtime.Runtime{}

	DefaultStores = map[string]func(...store.Option) store.Store{}

	DefaultTracers = map[string]func(...trace.Option) trace.Tracer{}

	DefaultAuths = map[string]func(...auth.Option) auth.Auth{}

	DefaultProfiles = map[string]func(...profile.Option) profile.Profile{
		"http":  http.NewProfile,
		"pprof": pprof.NewProfile,
	}

	DefaultConfigs = map[string]func(...config.Option) (config.Config, error){}

	DefaultCaches = map[string]func(...cache.Option) cache.Cache{}
)
