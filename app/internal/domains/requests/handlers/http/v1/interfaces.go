package v1

import (
	"affice/internal/domains/requests/stories/executer"
	"context"
)

type Requester interface {
	MakeRequest(ctx context.Context, urls []string) (responseList []executer.Response, errFun error)
}
