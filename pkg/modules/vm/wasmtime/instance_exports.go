package wasmtime

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/bytecodealliance/wasmtime-go"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	conflog "github.com/iotexproject/Bumblebee/conf/log"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"

	"github.com/iotexproject/w3bstream/pkg/models"
	"github.com/iotexproject/w3bstream/pkg/modules/httpclient"
	"github.com/iotexproject/w3bstream/pkg/types/wasm"
)

const (
	logTraceLevel uint32 = iota + 1
	logDebugLevel
	logInfoLevel
	logWarnLevel
	logErrorLevel
)

type (
	ExportFuncs struct {
		instance *Instance
		logger   conflog.Logger
		cl       *ChainClient
	}

	ChainClient struct {
		pvk   *ecdsa.PrivateKey
		chain *ethclient.Client
	}
)

func (ef *ExportFuncs) Log(c *wasmtime.Caller, logLevel, ptr, size int32) int32 {
	membuf := c.GetExport("memory").Memory().UnsafeData(ef.instance.vmStore)
	buf, err := read(membuf, ptr, size)
	if err != nil {
		ef.logger.Error(err)
		return wasm.ResultStatusCode_Failed
	}
	switch uint32(logLevel) {
	case logTraceLevel:
		ef.logger.Trace(string(buf))
	case logDebugLevel:
		ef.logger.Debug(string(buf))
	case logInfoLevel:
		ef.logger.Info(string(buf))
	case logWarnLevel:
		ef.logger.Warn(errors.New(string(buf)))
	case logErrorLevel:
		ef.logger.Error(errors.New(string(buf)))
	default:
		return wasm.ResultStatusCode_Failed
	}
	return int32(wasm.ResultStatusCode_OK)
}

func (ef *ExportFuncs) GetData(c *wasmtime.Caller, rid, vmAddrPtr, vmSizePtr int32) int32 {
	data, ok := ef.instance.res.Load(uint32(rid))
	if !ok {
		return int32(wasm.ResultStatusCode_ResourceNotFound)
	}

	if err := ef.copyDataIntoWasm(c, data, vmAddrPtr, vmSizePtr); err != nil {
		ef.logger.Error(err)
		return int32(wasm.ResultStatusCode_TransDataToVMFailed)
	}

	return int32(wasm.ResultStatusCode_OK)
}

func (ef *ExportFuncs) copyDataIntoWasm(c *wasmtime.Caller, data []byte, vmAddrPtr, vmSizePtr int32) error {
	allocFn := c.GetExport("alloc")
	if allocFn == nil {
		return errors.New("alloc is nil")
	}
	size := len(data)
	result, err := allocFn.Func().Call(ef.instance.vmStore, int32(size))
	if err != nil {
		return err
	}

	addr := result.(int32)

	memBuf := c.GetExport("memory").Memory().UnsafeData(ef.instance.vmStore)
	if siz := copy(memBuf[addr:], data); siz != size {
		return errors.New("fail to copy data")
	}

	// fmt.Printf("host >> addr=%d\n", addr)
	// fmt.Printf("host >> size=%d\n", size)
	// fmt.Printf("host >> vmAddrPtr=%d\n", vmAddrPtr)
	// fmt.Printf("host >> vmSizePtr=%d\n", vmSizePtr)

	if err := putUint32Le(memBuf, vmAddrPtr, uint32(addr)); err != nil {
		return err
	}
	if err := putUint32Le(memBuf, vmSizePtr, uint32(size)); err != nil {
		return err
	}

	return nil
}

// TODO SetData if rid not exist, should be assigned by wasm?
func (ef *ExportFuncs) SetData(c *wasmtime.Caller, rid, addr, size int32) int32 {
	memBuf := c.GetExport("memory").Memory().UnsafeData(ef.instance.vmStore)
	if addr > int32(len(memBuf)) || addr+size > int32(len(memBuf)) {
		return int32(wasm.ResultStatusCode_TransDataToVMFailed)
	}
	buf, err := read(memBuf, addr, size)
	if err != nil {
		ef.logger.Error(err)
		return int32(wasm.ResultStatusCode_TransDataToVMFailed)
	}
	ef.instance.res.Store(uint32(rid), buf)
	return int32(wasm.ResultStatusCode_OK)
}

func (ef *ExportFuncs) SetDB(c *wasmtime.Caller, kAddr, kSize, vAddr, vSize int32) int32 {
	memBuf := c.GetExport("memory").Memory().UnsafeData(ef.instance.vmStore)
	key, err := read(memBuf, kAddr, kSize)
	if err != nil {
		ef.logger.Error(err)
		return int32(wasm.ResultStatusCode_ResourceNotFound)
	}
	value, err := read(memBuf, vAddr, vSize)
	if err != nil {
		ef.logger.Error(err)
		return int32(wasm.ResultStatusCode_ResourceNotFound)
	}

	ef.logger.WithValues(
		"key", string(key),
		"val", string(value),
	).Info("host.SetDB")

	ef.instance.db[string(key)] = value
	return int32(wasm.ResultStatusCode_OK)
}

func (ef *ExportFuncs) GetDB(c *wasmtime.Caller,
	kAddr, kSize int32, vmAddrPtr, vmSizePtr int32) int32 {
	memBuf := c.GetExport("memory").Memory().UnsafeData(ef.instance.vmStore)
	key, err := read(memBuf, kAddr, kSize)
	if err != nil {
		ef.logger.Error(err)
		return int32(wasm.ResultStatusCode_ResourceNotFound)
	}

	val, exist := ef.instance.db[string(key)]
	if !exist || val == nil {
		return int32(wasm.ResultStatusCode_ResourceNotFound)
	}

	ef.logger.WithValues(
		"key", string(key),
		"val", string(val),
	).Info("host.GetDB")

	if err := ef.copyDataIntoWasm(c, val, vmAddrPtr, vmSizePtr); err != nil {
		ef.logger.Error(err)
		return int32(wasm.ResultStatusCode_TransDataToVMFailed)
	}

	return int32(wasm.ResultStatusCode_OK)
}

// TODO: add chainID in sendtx abi
// TODO: make sendTX async, and add callback if possible
func (ef *ExportFuncs) SendTX(c *wasmtime.Caller,
	chainID, payloadOffset, payloadSize int32,
	onSuccessEventTypeAddr, onSuccessEventTypeSize int32,
	onFailureEventTypeAddr, onFailureEventTypeSize int32,
	hashAddrPtr, hashSizePtr int32) int32 {
	if ef.cl == nil {
		ef.logger.Error(errors.New("eth client doesn't exist"))
		return wasm.ResultStatusCode_Failed
	}
	if ef.cl.pvk == nil {
		ef.logger.Error(errors.New("private key is empty"))
		return wasm.ResultStatusCode_Failed
	}
	memBuf := c.GetExport("memory").Memory().UnsafeData(ef.instance.vmStore)
	buf, err := read(memBuf, payloadOffset, payloadSize)
	if err != nil {
		ef.logger.Error(err)
		return wasm.ResultStatusCode_Failed
	}
	ret := gjson.Parse(string(buf))
	// fmt.Println(ret)
	txHash, err := sendETHTx(ef.cl, ret.Get("to").String(), ret.Get("value").String(), ret.Get("data").String())
	if err != nil {
		ef.logger.Error(err)
		return wasm.ResultStatusCode_Failed
	}
	if err := ef.copyDataIntoWasm(c, []byte(txHash), hashAddrPtr, hashSizePtr); err != nil {
		return int32(wasm.ResultStatusCode_TransDataToVMFailed)
	}
	ef.logger.Info("tx hash: %s", txHash)
	if onSuccessEventTypeSize > 0 || onFailureEventTypeSize > 0 {
		successEventType, err := read(memBuf, onSuccessEventTypeAddr, onSuccessEventTypeSize)
		if err != nil {
			ef.logger.Error(err)
			return wasm.ResultStatusCode_Failed
		}
		failureEventType, err := read(memBuf, onFailureEventTypeAddr, onFailureEventTypeSize)
		if err != nil {
			ef.logger.Error(err)
			return wasm.ResultStatusCode_Failed
		}
		if err := ef.monitorTx(uint64(chainID), txHash, string(successEventType), string(failureEventType)); err != nil {
			ef.logger.Error(err)
			return wasm.ResultStatusCode_Failed
		}
	}
	return int32(wasm.ResultStatusCode_OK)
}

func sendETHTx(cl *ChainClient, toStr string, valueStr string, dataStr string) (string, error) {
	var (
		sender = crypto.PubkeyToAddress(cl.pvk.PublicKey)
		to     = common.HexToAddress(toStr)
	)
	value, ok := new(big.Int).SetString(valueStr, 10)
	if !ok {
		return "", errors.New("fail to read tx value")
	}
	data, err := hex.DecodeString(dataStr)
	if err != nil {
		return "", err

	}
	nonce, err := cl.chain.PendingNonceAt(context.Background(), sender)
	if err != nil {
		return "", err
	}

	gasPrice, err := cl.chain.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}

	msg := ethereum.CallMsg{
		From:     sender,
		To:       &to,
		GasPrice: gasPrice,
		Value:    value,
		Data:     data,
	}
	gasLimit, err := cl.chain.EstimateGas(context.Background(), msg)
	if err != nil {
		return "", err
	}

	// Create a new transaction
	tx := types.NewTx(
		&types.LegacyTx{
			Nonce:    nonce,
			GasPrice: gasPrice,
			Gas:      gasLimit,
			To:       &to,
			Value:    value,
			Data:     data,
		})

	chainid, err := cl.chain.ChainID(context.Background())
	if err != nil {
		return "", err
	}
	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainid), cl.pvk)
	if err != nil {
		return "", err
	}
	err = cl.chain.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}
	return signedTx.Hash().Hex(), nil
}

func (ef *ExportFuncs) monitorTx(chainID uint64, txHash string,
	successEventType string, failureEventType string) error {
	monitorReq := models.CreateMonitorReq2{
		Chaintx: &models.ChaintxInfo{
			ChainID:   chainID,
			EventType: successEventType,
			TxAddress: txHash,
		},
	}
	body, err := json.Marshal(monitorReq)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("http://localhost:8888/srv-applet-mgr/v0/monitor/%s", ef.instance.projectID) // TODO move to config
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	cli := httpclient.NewClient()
	_, err = cli.Send(httpReq)
	if err != nil {
		return err
	}
	return nil
}

func (ef *ExportFuncs) CallContract(c *wasmtime.Caller,
	chainID, offset, size int32, vmAddrPtr, vmSizePtr int32) int32 {
	if ef.cl == nil {
		ef.logger.Error(errors.New("eth client doesn't exist"))
		return wasm.ResultStatusCode_Failed
	}
	if ef.cl.pvk == nil {
		ef.logger.Error(errors.New("private key is empty"))
		return wasm.ResultStatusCode_Failed
	}
	memBuf := c.GetExport("memory").Memory().UnsafeData(ef.instance.vmStore)
	buf, err := read(memBuf, offset, size)
	if err != nil {
		ef.logger.Error(err)
		return wasm.ResultStatusCode_Failed
	}
	ret := gjson.Parse(string(buf))
	// fmt.Println(ret)
	data, err := callContract(ef.cl.chain, ret.Get("to").String(), ret.Get("data").String())
	if err != nil {
		ef.logger.Error(err)
		return wasm.ResultStatusCode_Failed
	}
	if err := ef.copyDataIntoWasm(c, data, vmAddrPtr, vmSizePtr); err != nil {
		ef.logger.Error(err)
		return wasm.ResultStatusCode_Failed
	}
	return int32(wasm.ResultStatusCode_OK)
}

func callContract(cl *ethclient.Client, toStr string, dataStr string) ([]byte, error) {
	var (
		to      = common.HexToAddress(toStr)
		data, _ = hex.DecodeString(dataStr)
	)

	msg := ethereum.CallMsg{
		To:   &to,
		Data: data,
	}

	return cl.CallContract(context.Background(), msg, nil)
}

func putUint32Le(buf []byte, addr int32, num uint32) error {
	if int32(len(buf)) < addr+4 {
		return errors.New("overflow")
	}
	binary.LittleEndian.PutUint32(buf[addr:], num)
	return nil
}

func read(memBuf []byte, addr int32, size int32) ([]byte, error) {
	if size == 0 {
		return []byte{}, nil
	}
	if addr > int32(len(memBuf)) || addr+size > int32(len(memBuf)) {
		return nil, errors.New("overflow")
	}
	buf := make([]byte, size)
	if siz := copy(buf, memBuf[addr:addr+size]); int32(siz) != size {
		return nil, errors.New("overflow")
	}
	return buf, nil
}
