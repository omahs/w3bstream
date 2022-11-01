package monitor

import (
	"context"

	"github.com/iotexproject/Bumblebee/kit/httptransport/httpx"

	"github.com/iotexproject/w3bstream/pkg/modules/blockchain"
	"github.com/iotexproject/w3bstream/pkg/types"
)

type CreateMonitor struct {
	httpx.MethodPost
	ProjectID                   types.SFID `in:"path" name:"projectID"`
	blockchain.CreateMonitorReq `in:"body"`
}

func (r *CreateMonitor) Path() string { return "/:projectID" }

func (r *CreateMonitor) Output(ctx context.Context) (interface{}, error) {
	// fmt.Println("DEBUG 1", r.ProjectID)
	// ca := middleware.CurrentAccountFromContext(ctx)
	// fmt.Println("DEBUG 2")
	// p, err := ca.ValidateProjectPerm(ctx, r.ProjectID)
	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Println("DEBUG 3")
	name := "project2"
	return blockchain.CreateMonitor(ctx, name, &r.CreateMonitorReq)
}
