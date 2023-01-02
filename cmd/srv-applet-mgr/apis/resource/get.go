package resource

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/resource"
)

type ListResources struct {
	httpx.MethodGet
}

func (r *ListResources) Output(ctx context.Context) (interface{}, error) {
	return resource.ListResource(ctx)
}
