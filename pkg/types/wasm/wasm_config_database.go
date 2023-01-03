package wasm

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/schema"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/types"
)

type Schema struct{ schema.Schema }

func (s *Schema) ConfigType() enums.ConfigType {
	return enums.CONFIG_TYPE__PROJECT_SCHEMA
}

func (s *Schema) WithContext(ctx context.Context) context.Context {
	return WithDBExecutor(ctx, s.DBExecutor(types.MustWasmDBExecutorFromContext(ctx)))
}
