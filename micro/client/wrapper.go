package client

import (
	"context"

	"github.com/wxc/micro/registry"
)

type CallFunc func(ctx context.Context, node *registry.Node, req Request, rsp interface{}, opts CallOptions) error

type CallWrapper func(CallFunc) CallFunc

type Wrapper func(Client) Client

type StreamWrapper func(Stream) Stream
