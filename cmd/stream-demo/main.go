package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/reactivex/rxgo/v2"
	"syreclabs.com/go/faker"

	"github.com/machinefi/w3bstream/cmd/stream-demo/global"
	"github.com/machinefi/w3bstream/cmd/stream-demo/models"
	"github.com/machinefi/w3bstream/cmd/stream-demo/tasks"
	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/modules/vm"
	"github.com/machinefi/w3bstream/pkg/modules/vm/wasmtime"
	"github.com/machinefi/w3bstream/pkg/types"
)

func main() {

	global.Migrate()

	ch := make(chan rxgo.Item)
	go producer(ch)

	ctx := global.WithContext(context.Background())
	//ctx = types.WithProject()
	//ctx = types.WithApplet()

	log := types.MustLoggerFromContext(ctx)
	d := types.MustInsDBExecutorFromContext(ctx)

	go kit.Run(tasks.Root, global.TaskServer())

	path := "/Users/hunhun/project_go/iotex/w3bstream/_examples/stream/stream.wasm"

	// new wasm runtime instance
	ins, err := newWasmRuntimeInstance(ctx, path)
	if err != nil {
		log.Panic(err)
	}
	id := vm.AddInstance(ctx, ins)

	//TODO mirco server 不需要mgr instance
	err = vm.StartInstance(ctx, id)
	defer vm.StopInstance(ctx, id)

	observable := rxgo.FromChannel(ch).Filter(func(i interface{}) bool {
		res := false

		// 1.序列化
		b, err := json.Marshal(i.(models.Customer))
		if err != nil {
			log.Error(err)
		}

		// 2.调用wasm code
		//res = i.(Customer).Age >= 18
		code := ins.HandleEvent(ctx, "start", b).Code
		log.Info(fmt.Sprintf("start wasm code %d", code))

		// 3.get/parse data
		if code < 0 {
			//TODO filter cant send error event
			log.Error(errors.New(fmt.Sprintf("%v filter error.", i.(models.Customer))))
			return res
			//errors.New(fmt.Sprintf("%v filter error.", i.(models.Customer)))
		}

		rb, ok := ins.GetResource(uint32(code))
		if !ok {
			log.Error(errors.New("not found"))
		}

		result := strings.ToLower(string(rb))
		if result == "true" {
			res = true
		} else if result == "false" {
			res = false
		} else {
			log.Warn(errors.New("the value does not support"))
		}

		// 4.return data
		return res
	}).Map(func(ctx context.Context, i interface{}) (interface{}, error) {
		// 1.序列化
		b, err := json.Marshal(i.(models.Customer))
		if err != nil {
			log.Error(err)
		}

		// 2.调用wasm code
		code := ins.HandleEvent(ctx, "mapTax", b).Code
		log.Info(fmt.Sprintf("mapTax wasm code %d", code))

		// 3.get/parse data
		if code < 0 {
			log.Error(errors.New(fmt.Sprintf("%v %s error.", i.(models.Customer), "mapTax")))
			return nil, errors.New(fmt.Sprintf("%v %s error.", i.(models.Customer), "mapTax"))
		}

		rb, ok := ins.GetResource(uint32(code))
		if !ok {
			log.Error(errors.New("not found"))
		}

		customer := &models.Customer{}
		err = json.Unmarshal(rb, customer)

		return customer, err
	}).GroupByDynamic(func(item rxgo.Item) string {
		// 1.序列化
		b, err := json.Marshal(item.V.(models.Customer))
		if err != nil {
			log.Error(err)
		}

		// 2.调用wasm code
		code := ins.HandleEvent(ctx, "mapTax", b).Code
		log.Info(fmt.Sprintf("mapTax wasm code %d", code))

		// 3.get/parse data
		if code < 0 {
			log.Error(errors.New(fmt.Sprintf("%v %s error.", item.V.(models.Customer), "mapTax")))
			//TODO error 怎么处理
			return "error"
		}

		rb, ok := ins.GetResource(uint32(code))
		if !ok {
			log.Error(errors.New("not found"))
		}

		groupKey := string(rb)
		return groupKey
	}, rxgo.WithBufferedChannel(10), rxgo.WithErrorStrategy(rxgo.ContinueOnError))

	c := observable.Observe()
	for item := range c {
		fmt.Println(item.V)

		switch item.V.(type) {
		case rxgo.GroupedObservable: // group operator
			go func() {
				obs := item.V.(rxgo.GroupedObservable)
				fmt.Printf("New observable: %s\n", obs.Key)
				for i := range obs.Observe() {
					fmt.Printf("item: %v\n", i.V)
					customer := i.V.(models.Customer)
					fmt.Println(fmt.Sprintf("customer: %v", customer))
					customer.Create(d)
				}
			}()
		case rxgo.ObservableImpl: // window operator
			obs := item.V.(rxgo.ObservableImpl)
			for i := range obs.Observe() {
				//for i := range obs.Count().Observe() {
				fmt.Printf("item: %v\n", i.V)
			}
		default:
			fmt.Printf("item: %v\n", item.V)
		}
	}
}

func producer(ch chan<- rxgo.Item) {
	for _ = range time.Tick(time.Second) {
		for i := 0; i < 10; i++ {
			ch <- rxgo.Of(models.Customer{
				ID:        faker.Code().Isbn10(),
				FirstName: faker.Name().FirstName(),
				LastName:  faker.Name().LastName(),
				Age:       faker.RandomInt(15, 53),
				City:      faker.Address().City(),
			})
		}
	}
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
