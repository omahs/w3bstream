package main

import (
	"context"
	"encoding/json"
	"fmt"
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

func main() {

	global.Migrate()

	ch := make(chan rxgo.Item)

	ctx := global.WithContext(context.Background())
	//ctx = types.WithProject()
	//ctx = types.WithApplet()

	log := types.MustLoggerFromContext(ctx)
	d := types.MustDBExecutorFromContext(ctx)
	//d := types.MustInsDBExecutorFromContext(ctx)

	channel := "stream-dome"

	go func() {
		if err := initChannel(ctx, channel, func(ctx context.Context, channel string, data *eventpb.Event) (interface{}, error) {
			return onEventReceived(ctx, channel, ch, data)
		}); err != nil {
			log.Panic(err)
		}
	}()

	go kit.Run(tasks.Root, global.TaskServer())

	path := "./_examples/stream/stream.wasm"

	// new wasm runtime instance
	ins, err := newWasmRuntimeInstance(ctx, path)
	if err != nil {
		log.Panic(err)
	}
	id := vm.AddInstance(ctx, ins)

	//TODO mirco server, Does need mgr instance?
	err = vm.StartInstance(ctx, id)
	defer vm.StopInstance(ctx, id)

	observable := rxgo.FromChannel(ch).Filter(func(i interface{}) bool {
		res := false

		// 1.Serialize
		b, err := json.Marshal(i.(models.Customer))
		if err != nil {
			log.Error(err)
		}

		// 2.Invoke wasm code
		//res = i.(Customer).Age >= 18
		code := ins.HandleEvent(ctx, "start", b).Code
		log.Info(fmt.Sprintf("start wasm code %d", code))

		// 3.Get & parse data
		if code < 0 {
			//TODO filter cant send error event
			log.Error(errors.New(fmt.Sprintf("%v filter error.", i.(models.Customer))))
			return res
			//errors.New(fmt.Sprintf("%v filter error.", i.(models.Customer)))
		}

		rb, ok := ins.GetResource(uint32(code))
		//TODO remove resource
		defer ins.RmvResource(ctx, uint32(code))
		if !ok {
			log.Error(errors.New("not found"))
			return res
		}

		result := strings.ToLower(string(rb))
		if result == "true" {
			res = true
		} else if result == "false" {
			res = false
		} else {
			log.Warn(errors.New("the value does not support"))
		}

		// 4.Return data
		return res
	}).Map(func(c context.Context, i interface{}) (interface{}, error) {
		// 1.Serialize
		b, err := json.Marshal(i.(models.Customer))
		if err != nil {
			log.Error(err)
		}

		// 2.Invoke wasm code
		code := ins.HandleEvent(ctx, "mapTax", b).Code
		log.Info(fmt.Sprintf("mapTax wasm code %d", code))

		// 3.Get & parse data
		if code < 0 {
			log.Error(errors.New(fmt.Sprintf("%v %s error.", i.(models.Customer), "mapTax")))
			return nil, errors.New(fmt.Sprintf("%v %s error.", i.(models.Customer), "mapTax"))
		}

		rb, ok := ins.GetResource(uint32(code))
		defer ins.RmvResource(ctx, uint32(code))
		if !ok {
			log.Error(errors.New("mapTax result not found"))
			return nil, errors.New("mapTax result not found")
		}

		customer := models.Customer{}
		err = json.Unmarshal(rb, &customer)

		return customer, err
	}).GroupByDynamic(func(item rxgo.Item) string {
		// 1.Serialize
		b, err := json.Marshal(item.V.(models.Customer))
		if err != nil {
			log.Error(err)
		}

		// 2.Invoke wasm code
		code := ins.HandleEvent(ctx, "groupByAge", b).Code
		log.Info(fmt.Sprintf("groupByAge wasm code %d", code))

		// 3.Get & parse data
		if code < 0 {
			log.Error(errors.New(fmt.Sprintf("%v %s error.", item.V.(models.Customer), "groupByAge")))
			//TODO handle exceptions
			return "error"
		}

		rb, ok := ins.GetResource(uint32(code))
		defer ins.RmvResource(ctx, uint32(code))
		if !ok {
			log.Error(errors.New("groupByAge result not found"))
			//TODO handle exceptions
			return "error"
		}

		groupKey := string(rb)
		return groupKey
	}, rxgo.WithBufferedChannel(10), rxgo.WithErrorStrategy(rxgo.ContinueOnError))

	c := observable.Observe()
	for item := range c {
		//fmt.Println(item.V)

		switch item.V.(type) {
		case rxgo.GroupedObservable: // group operator
			go func() {
				obs := item.V.(rxgo.GroupedObservable)
				log.Info(fmt.Sprintf("New observable: %s", obs.Key))
				for i := range obs.Observe() {
					log.Info(fmt.Sprintf("item: %v", i.V))
					customer := i.V.(models.Customer)
					log.Info(fmt.Sprintf("customer: %v", customer))
					if err := customer.Create(d); err != nil {
						log.Error(err)
					}
				}
			}()
		case rxgo.ObservableImpl: // window operator
			obs := item.V.(rxgo.ObservableImpl)
			for i := range obs.Observe() {
				//for i := range obs.Count().Observe() {
				log.Info(fmt.Sprintf("item: %v", i.V))
			}
		default:
			log.Info(fmt.Sprintf("item: %v", item.V))
			customer := item.V.(models.Customer)
			log.Info(fmt.Sprintf("customer: %v", customer))
			if err := customer.Create(d); err != nil {
				log.Error(err)
			}
		}
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
	customer := models.Customer{}
	//fmt.Println("r.Payload " + r.Payload)
	json.Unmarshal([]byte(r.Payload), &customer)
	//fmt.Println(customer)
	ch <- rxgo.Of(customer)
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
