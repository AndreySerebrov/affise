package executer

import (
	"context"
	"net/url"
)

//go:generate mockgen -source=interfaces.go -destination=mocks.go -package=executer

type Limiter interface {
	Take() chan bool
	Free()
}

type Requester interface {
	Do(ctx context.Context, url *url.URL) (string, error)
}
