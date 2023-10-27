package client

import (
	"context"

	"github.com/wxc/micro/errors"
)

type RetryFunc func(ctx context.Context, req Request, retryCount int, err error) (bool, error)

func RetryAlways(ctx context.Context, req Request, retryCount int, err error) (bool, error) {
	return true, nil
}

func RetryOnError(ctx context.Context, req Request, retryCount int, err error) (bool, error) {
	if err == nil {
		return false, nil
	}

	e := errors.Parse(err.Error())
	if e == nil {
		return false, nil
	}

	switch e.Code {
	case 408, 500:
		return true, nil
	default:
		return false, nil
	}
}
