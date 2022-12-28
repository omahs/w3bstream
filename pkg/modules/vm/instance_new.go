package vm

import (
	"context"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
	"os"

	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/modules/vm/wasmtime"
	"github.com/machinefi/w3bstream/pkg/types"
)

func NewInstance(ctx context.Context, path string, id types.SFID) error {
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "NewInstanceWithID")
	defer l.End()

	code, err := os.ReadFile(path)
	if err != nil {
		l.Error(err)
		return status.ExtractFileFailed.StatusErr().WithDesc(err.Error())
	}

	//TODO setting config
	instanceConfig := &wasm.InstanceConfig{
		KvType: 0,
	}

	i, err := wasmtime.NewInstanceByCode(ctx, id, code, instanceConfig)
	if err != nil {
		l.Error(err)
		return err
	}

	AddInstanceByID(ctx, id, i)
	return nil
}
