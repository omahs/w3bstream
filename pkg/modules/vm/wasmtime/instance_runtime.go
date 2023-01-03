package wasmtime

import (
	"encoding/binary"

	"github.com/bytecodealliance/wasmtime-go"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
	"github.com/pkg/errors"
)

type Runtime struct {
	engine   *wasmtime.Engine
	store    *wasmtime.Store
	linker   *wasmtime.Linker
	module   *wasmtime.Module
	instance *wasmtime.Instance
}

func NewRuntime(code []byte, abi wasm.ABI) (rt *Runtime, err error) {
	rt = &Runtime{}
	rt.engine = wasmtime.NewEngineWithConfig(wasmtime.NewConfig())
	rt.store = wasmtime.NewStore(rt.engine)
	rt.linker = wasmtime.NewLinker(rt.engine)

	_ = rt.linker.FuncWrap("env", "ws_log", abi.Log)
	_ = rt.linker.FuncWrap("env", "ws_get_data", abi.GetData)
	_ = rt.linker.FuncWrap("env", "ws_set_data", abi.SetData)
	_ = rt.linker.FuncWrap("env", "ws_get_db", abi.GetDB)
	_ = rt.linker.FuncWrap("env", "ws_set_db", abi.SetDB)
	_ = rt.linker.FuncWrap("env", "ws_send_tx", abi.SendTX)
	_ = rt.linker.FuncWrap("env", "ws_call_contract", abi.CallContract)
	_ = rt.linker.FuncWrap("env", "ws_set_sql_db", abi.SetSQLDB)
	_ = rt.linker.FuncWrap("env", "ws_get_sql_db", abi.GetSQLDB)
	_ = rt.linker.FuncWrap("env", "ws_get_env", abi.GetEnv)
	_ = rt.linker.DefineWasi()
	rt.store.SetWasi(wasmtime.NewWasiConfig())

	rt.module, err = wasmtime.NewModule(rt.engine, code)
	if err != nil {
		return nil, err
	}
	rt.instance, err = rt.linker.Instantiate(rt.store, rt.module)
	if err != nil {
		return nil, err
	}
	return
}

func (rt *Runtime) NewMemory() []byte {
	return rt.instance.GetExport(rt.store, "memory").Memory().UnsafeData(rt.store)
}

func (rt *Runtime) Alloc(size int32) (int32, []byte, error) {
	fn := rt.ExportFunc("alloc")
	if fn == nil {
		return 0, nil, errors.New("alloc is nil")
	}
	result, err := fn.Func().Call(rt.store, size)
	if err != nil {
		return 0, nil, err
	}
	return result.(int32), rt.NewMemory(), nil
}

func (rt *Runtime) ExportFunc(name string) *wasmtime.Extern {
	return rt.instance.GetExport(rt.store, name)
}

func (rt *Runtime) GetFunc(name string) *wasmtime.Func {
	return rt.instance.GetFunc(rt.store, name)
}

func (rt *Runtime) Call(name string, args ...interface{}) (interface{}, error) {
	fn := rt.GetFunc(name)
	if fn == nil {
		return nil, errors.Errorf("runtime: %s fn is not imported", name)
	}
	return fn.Call(rt.store, args...)
}

func (rt *Runtime) Read(addr, size int32) ([]byte, error) {
	mem := rt.NewMemory()
	if addr > int32(len(mem)) || addr+size > int32(len(mem)) {
		return nil, errors.New("overflow")
	}
	buf := make([]byte, size)
	if copied := copy(buf, mem[addr:addr+size]); int32(copied) != size {
		return nil, errors.New("overflow")
	}
	return buf, nil
}

func (rt *Runtime) Copy(hostData []byte, vmAddrPtr, vmSizePtr int32) error {
	size := len(hostData)
	addr, mem, err := rt.Alloc(int32(size))
	if copied := copy(mem[addr:], hostData); copied != size {
		return errors.New("fail to copy data")
	}
	if err = rt.PutUint32Le(mem, vmAddrPtr, uint32(addr)); err != nil {
		return err
	}
	if err = rt.PutUint32Le(mem, vmSizePtr, uint32(size)); err != nil {
		return err
	}

	return nil
}

func (rt *Runtime) PutUint32Le(buf []byte, vmAddr int32, val uint32) error {
	if int32(len(buf)) < vmAddr+4 {
		return errors.New("overflow")
	}
	binary.LittleEndian.PutUint32(buf[vmAddr:], val)
	return nil
}
