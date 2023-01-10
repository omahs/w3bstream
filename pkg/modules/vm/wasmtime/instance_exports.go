package wasmtime

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/types/wasm/sql_util"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"

	conflog "github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/x/mapx"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

const (
	logTraceLevel uint32 = iota + 1
	logDebugLevel
	logInfoLevel
	logWarnLevel
	logErrorLevel
)

type (
	Import func(module, name string, f interface{}) error

	ABILinker interface {
		LinkABI(Import) error
	}

	ExportFuncs struct {
		rt  *Runtime
		res *mapx.Map[uint32, []byte]
		env *wasm.Env
		kvs wasm.KVStore
		db  sqlx.DBExecutor
		log conflog.Logger
		cl  *wasm.ChainClient
	}
)

func NewExportFuncs(ctx context.Context, rt *Runtime) (*ExportFuncs, error) {
	ef := &ExportFuncs{
		res: wasm.MustRuntimeResourceFromContext(ctx),
		kvs: wasm.MustKVStoreFromContext(ctx),
		db:  wasm.MustDBExecutorFromContext(ctx),
		log: wasm.MustLoggerFromContext(ctx),
		env: wasm.MustEnvFromContext(ctx),
	}
	ef.cl, _ = wasm.ChainClientFromContext(ctx)
	ef.rt = rt

	return ef, nil
}

var _ wasm.ABI = (*ExportFuncs)(nil)

func (ef *ExportFuncs) LinkABI(impt Import) error {
	for name, ff := range map[string]interface{}{
		"ws_log":           ef.Log,
		"ws_get_data":      ef.GetData,
		"ws_set_data":      ef.SetData,
		"ws_get_db":        ef.GetDB,
		"ws_set_db":        ef.SetDB,
		"ws_send_tx":       ef.SendTX,
		"ws_call_contract": ef.CallContract,
		"ws_set_sql_db":    ef.SetSQLDB,
		"ws_get_sql_db":    ef.GetSQLDB,
		"ws_get_env":       ef.GetEnv,
	} {
		if err := impt("env", name, ff); err != nil {
			return err
		}
	}

	return nil
}

func (ef *ExportFuncs) Log(logLevel, ptr, size int32) int32 {
	buf, err := ef.rt.Read(ptr, size)
	if err != nil {
		ef.log.Error(err)
		return wasm.ResultStatusCode_Failed
	}
	switch uint32(logLevel) {
	case logTraceLevel:
		ef.log.Trace(string(buf))
	case logDebugLevel:
		ef.log.Debug(string(buf))
	case logInfoLevel:
		ef.log.Info(string(buf))
	case logWarnLevel:
		ef.log.Warn(errors.New(string(buf)))
	case logErrorLevel:
		ef.log.Error(errors.New(string(buf)))
	default:
		return wasm.ResultStatusCode_Failed
	}
	return int32(wasm.ResultStatusCode_OK)
}

func (ef *ExportFuncs) GetData(rid, vmAddrPtr, vmSizePtr int32) int32 {
	data, ok := ef.res.Load(uint32(rid))
	if !ok {
		return int32(wasm.ResultStatusCode_ResourceNotFound)
	}

	if err := ef.rt.Copy(data, vmAddrPtr, vmSizePtr); err != nil {
		ef.log.Error(err)
		return int32(wasm.ResultStatusCode_TransDataToVMFailed)
	}

	return int32(wasm.ResultStatusCode_OK)
}

// TODO SetData if rid not exist, should be assigned by wasm?
func (ef *ExportFuncs) SetData(rid, addr, size int32) int32 {
	buf, err := ef.rt.Read(addr, size)
	if err != nil {
		ef.log.Error(err)
		return int32(wasm.ResultStatusCode_TransDataToVMFailed)
	}
	ef.res.Store(uint32(rid), buf)
	return int32(wasm.ResultStatusCode_OK)
}

func (ef *ExportFuncs) GetDB(kAddr, kSize int32, vmAddrPtr, vmSizePtr int32) int32 {
	key, err := ef.rt.Read(kAddr, kSize)
	if err != nil {
		ef.log.Error(err)
		return int32(wasm.ResultStatusCode_ResourceNotFound)
	}

	val, exist := ef.kvs.Get(string(key))
	if exist != nil || val == nil {
		return int32(wasm.ResultStatusCode_ResourceNotFound)
	}

	ef.log.WithValues("key", string(key), "val", string(val)).Info("host.GetDB")

	if err := ef.rt.Copy(val, vmAddrPtr, vmSizePtr); err != nil {
		ef.log.Error(err)
		return int32(wasm.ResultStatusCode_TransDataToVMFailed)
	}

	return int32(wasm.ResultStatusCode_OK)
}

func (ef *ExportFuncs) SetDB(kAddr, kSize, vAddr, vSize int32) int32 {
	key, err := ef.rt.Read(kAddr, kSize)
	if err != nil {
		ef.log.Error(err)
		return int32(wasm.ResultStatusCode_ResourceNotFound)
	}
	val, err := ef.rt.Read(vAddr, vSize)
	if err != nil {
		ef.log.Error(err)
		return int32(wasm.ResultStatusCode_ResourceNotFound)
	}
	ef.log.WithValues("key", string(key), "val", string(val)).Info("host.SetDB")

	err = ef.kvs.Set(string(key), val)
	if err != nil {
		ef.log.Error(err)
		return int32(wasm.ResultStatusCode_Failed)
	}
	return int32(wasm.ResultStatusCode_OK)
}

func (ef *ExportFuncs) SetSQLDB(addr, size int32) int32 {
	data, err := ef.rt.Read(addr, size)
	if err != nil {
		ef.log.Error(err)
		return int32(wasm.ResultStatusCode_ResourceNotFound)
	}

	prestate, params, err := sql_util.ParseQuery(data)
	if err != nil {
		ef.log.Error(err)
		return wasm.ResultStatusCode_Failed
	}

	_, err = ef.db.ExecContext(context.Background(), prestate, params...)
	if err != nil {
		ef.log.Error(err)
		return wasm.ResultStatusCode_Failed
	}

	return int32(wasm.ResultStatusCode_OK)
}

func (ef *ExportFuncs) GetSQLDB(addr, size int32, vmAddrPtr, vmSizePtr int32) int32 {
	data, err := ef.rt.Read(addr, size)
	if err != nil {
		ef.log.Error(err)
		return int32(wasm.ResultStatusCode_ResourceNotFound)
	}

	prestate, params, err := sql_util.ParseQuery(data)
	if err != nil {
		ef.log.Error(err)
		return wasm.ResultStatusCode_Failed
	}

	rows, err := ef.db.QueryContext(context.Background(), prestate, params...)
	if err != nil {
		ef.log.Error(err)
		return wasm.ResultStatusCode_Failed
	}

	ret, err := sql_util.JsonifyRows(rows)
	if err != nil {
		ef.log.Error(err)
		return wasm.ResultStatusCode_Failed
	}

	if err := ef.rt.Copy(ret, vmAddrPtr, vmSizePtr); err != nil {
		ef.log.Error(err)
		return int32(wasm.ResultStatusCode_TransDataToVMFailed)
	}

	return int32(wasm.ResultStatusCode_OK)
}

// TODO: add chainID in sendtx abi
// TODO: make sendTX async, and add callback if possible
func (ef *ExportFuncs) SendTX(offset, size int32) int32 {
	if ef.cl == nil {
		ef.log.Error(errors.New("eth client doesn't exist"))
		return wasm.ResultStatusCode_Failed
	}
	buf, err := ef.rt.Read(offset, size)
	if err != nil {
		ef.log.Error(err)
		return wasm.ResultStatusCode_Failed
	}
	ret := gjson.Parse(string(buf))
	// fmt.Println(ret)
	txHash, err := ef.cl.SendTX(ret.Get("to").String(), ret.Get("value").String(), ret.Get("data").String())
	if err != nil {
		ef.log.Error(err)
		return wasm.ResultStatusCode_Failed
	}
	ef.log.Info("tx hash: %s", txHash)
	return int32(wasm.ResultStatusCode_OK)
}

func (ef *ExportFuncs) CallContract(offset, size int32, vmAddrPtr, vmSizePtr int32) int32 {
	if ef.cl == nil {
		ef.log.Error(errors.New("eth client doesn't exist"))
		return wasm.ResultStatusCode_Failed
	}
	buf, err := ef.rt.Read(offset, size)
	if err != nil {
		ef.log.Error(err)
		return wasm.ResultStatusCode_Failed
	}
	ret := gjson.Parse(string(buf))
	// fmt.Println(ret)
	data, err := ef.cl.CallContract(ret.Get("to").String(), ret.Get("data").String())
	if err != nil {
		ef.log.Error(err)
		return wasm.ResultStatusCode_Failed
	}
	if err = ef.rt.Copy(data, vmAddrPtr, vmSizePtr); err != nil {
		ef.log.Error(err)
		return wasm.ResultStatusCode_Failed
	}
	return int32(wasm.ResultStatusCode_OK)
}

func (ef *ExportFuncs) GetEnv(kAddr, kSize int32, vmAddrPtr, vmSizePtr int32) int32 {
	key, err := ef.rt.Read(kAddr, kSize)
	if err != nil {
		return int32(wasm.ResultStatusCode_TransDataToVMFailed)
	}

	val, ok := ef.env.Get(string(key))
	if !ok {
		return int32(wasm.ResultStatusCode_EnvKeyNotFound)
	}

	if err = ef.rt.Copy([]byte(val), vmAddrPtr, vmSizePtr); err != nil {
		ef.log.Error(err)
		return int32(wasm.ResultStatusCode_TransDataToVMFailed)
	}
	return int32(wasm.ResultStatusCode_OK)
}
