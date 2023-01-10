package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/reactivex/rxgo/v2"

	"github.com/machinefi/w3bstream/cmd/stream-demo/global"
	"github.com/machinefi/w3bstream/cmd/stream-demo/tasks"
	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/depends/protocol/eventpb"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/mq"
	"github.com/machinefi/w3bstream/pkg/modules/vm"
	"github.com/machinefi/w3bstream/pkg/modules/vm/wasmtime"
	"github.com/machinefi/w3bstream/pkg/types"
)

var (
	ctx context.Context
	l   log.Logger
	d   sqlx.DBExecutor
	ins *wasmtime.Instance
	err error
)

func main() {

	global.Migrate()

	ch := make(chan rxgo.Item)

	ctx = global.WithContext(context.Background())

	l = types.MustLoggerFromContext(ctx)
	d = types.MustDBExecutorFromContext(ctx)

	channel := types.MustChannelNameFromContext(ctx)

	go func() {
		if err := initChannel(ctx, channel, func(ctx context.Context, channel string, data *eventpb.Event) (interface{}, error) {
			return onEventReceived(ctx, channel, ch, data)
		}); err != nil {
			l.Panic(err)
		}
	}()

	go kit.Run(tasks.Root, global.TaskServer())

	path := types.MustWasmPathFromContext(ctx)

	// new wasm runtime instance
	ins, err = newWasmRuntimeInstance(ctx, path)
	if err != nil {
		l.Panic(err)
	}
	id := vm.AddInstance(ctx, ins)

	//TODO mirco server, Does need mgr instance?
	err = vm.StartInstance(ctx, id)
	defer vm.StopInstance(ctx, id)

	observable := rxgo.FromChannel(ch).Filter(filterFunc).Map(mapFunc).
		GroupByDynamic(groupByKey, rxgo.WithBufferedChannel(10), rxgo.WithErrorStrategy(rxgo.ContinueOnError))

	c := observable.Observe()
	for item := range c {

		switch item.V.(type) {
		case rxgo.GroupedObservable: // group operator
			go func() {
				obs := item.V.(rxgo.GroupedObservable)
				for i := range obs.Observe() {
					sink(i)
				}
			}()
		case rxgo.ObservableImpl: // window operator
			obs := item.V.(rxgo.ObservableImpl)
			for i := range obs.Observe() {
				sink(i)
			}
		default:
			sink(item)
		}
	}
}

func filterFunc(i interface{}) bool {
	res := false

	// 1.Serialize
	b, err := json.Marshal(i.(models.SourceCustomer))
	if err != nil {
		l.Error(err)
	}

	// 2.Invoke wasm code
	code := ins.HandleEvent(ctx, "start", b).Code
	l.Info(fmt.Sprintf("start wasm code %d", code))

	// 3.Get & parse data
	if code < 0 {
		//TODO filter cant send error event
		l.Error(errors.New(fmt.Sprintf("%v filter error.", i.(models.Customer))))
		return res
		//errors.New(fmt.Sprintf("%v filter error.", i.(models.Customer)))
	}

	rb, ok := ins.GetResource(uint32(code))
	//TODO remove resource
	defer ins.RmvResource(ctx, uint32(code))
	if !ok {
		l.Error(errors.New("not found"))
		return res
	}

	result := strings.ToLower(string(rb))
	if result == "true" {
		res = true
	} else if result == "false" {
		res = false
	} else {
		l.Warn(errors.New("the value does not support"))
	}

	// 4.Return data
	return res
}

func mapFunc(c context.Context, i interface{}) (interface{}, error) {
	// 1.Serialize
	b, err := json.Marshal(i.(models.SourceCustomer))
	if err != nil {
		l.Error(err)
	}

	// 2.Invoke wasm code
	code := ins.HandleEvent(ctx, "mapTax", b).Code
	l.Info(fmt.Sprintf("mapTax wasm code %d", code))

	// 3.Get & parse data
	if code < 0 {
		l.Error(errors.New(fmt.Sprintf("%v %s error.", i.(models.Customer), "mapTax")))
		return nil, errors.New(fmt.Sprintf("%v %s error.", i.(models.Customer), "mapTax"))
	}

	rb, ok := ins.GetResource(uint32(code))
	defer ins.RmvResource(ctx, uint32(code))
	if !ok {
		l.Error(errors.New("mapTax result not found"))
		return nil, errors.New("mapTax result not found")
	}

	customer := models.Customer{}
	err = json.Unmarshal(rb, &customer)

	return customer, err
}

func groupByKey(item rxgo.Item) string {
	// 1.Serialize
	b, err := json.Marshal(item.V.(models.Customer))
	if err != nil {
		l.Error(err)
	}

	// 2.Invoke wasm code
	code := ins.HandleEvent(ctx, "groupByAge", b).Code
	l.Info(fmt.Sprintf("groupByAge wasm code %d", code))

	// 3.Get & parse data
	if code < 0 {
		l.Error(errors.New(fmt.Sprintf("%v %s error.", item.V.(models.Customer), "groupByAge")))
		//TODO handle exceptions
		return "error"
	}

	rb, ok := ins.GetResource(uint32(code))
	defer ins.RmvResource(ctx, uint32(code))
	if !ok {
		l.Error(errors.New("groupByAge result not found"))
		//TODO handle exceptions
		return "error"
	}

	groupKey := string(rb)
	return groupKey
}

func sink(item rxgo.Item) {
	customer := item.V.(models.Customer)
	l.Info(fmt.Sprintf("customer: %v", customer))
	if err := customer.Create(d); err != nil {
		l.Error(err)
	}
}

func initChannel(ctx context.Context, channel string, hdl mq.OnMessage) (err error) {
	err = mq.CreateChannel(ctx, channel, hdl)
	if err != nil {
		err = errors.Errorf("create channel: [channel:%s] [err:%v]", channel, err)
	}
	return err
}

func onEventReceived(ctx context.Context, projectName string, ch chan<- rxgo.Item, r *eventpb.Event) (interface{}, error) {
	sourceCustomer := models.SourceCustomer{}
	json.Unmarshal([]byte(r.Payload), &sourceCustomer)
	ch <- rxgo.Of(sourceCustomer)
	return nil, nil
}

func newWasmRuntimeInstance(ctx context.Context, path string) (*wasmtime.Instance, error) {
	idg := confid.MustSFIDGeneratorFromContext(ctx)
	insID := idg.MustGenSFID()

	code, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return wasmtime.NewInstanceByCode(ctx, insID, code)
}
