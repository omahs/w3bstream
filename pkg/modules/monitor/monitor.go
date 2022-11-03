package monitor

import (
	"context"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/blockchain"
	"github.com/machinefi/w3bstream/pkg/types"
)

type CreateMonitorReq struct {
	Contractlog *CreateContractlogReq `json:"contractLog,omitempty"`
	Chaintx     *CreateChaintxReq     `json:"chainTx,omitempty"`
	ChainHeight *CreateChainHeightReq `json:"chainHeight,omitempty"`
}

type (
	CreateContractlogReq = models.ContractlogInfo
	CreateChaintxReq     = models.ChaintxInfo
	CreateChainHeightReq = models.ChainHeightInfo
)

func CreateMonitor(ctx context.Context, project *models.Project, r *CreateMonitorReq) (interface{}, error) {
	d := types.MustDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)

	_, l = l.Start(ctx, "CreateMonitor")
	defer l.End()

	mt, err := getMonitorType(r)
	if err != nil {
		return nil, err
	}

	m := &models.Monitor{
		RelMonitor: models.RelMonitor{MonitorID: idg.MustGenSFID()},
		RelProject: models.RelProject{ProjectID: project.ProjectID},
		MonitorData: models.MonitorData{
			State: enums.MONITOR_STATE__SYNCING,
			Data: models.MonitorInfo{
				Type:        mt,
				ChainTx:     r.Chaintx,
				ContractLog: r.Contractlog,
				ChainHeight: r.ChainHeight,
			},
		},
	}
	if err := m.Create(d); err != nil {
		l.Error(err)
		return nil, status.CheckDatabaseError(err, "CreateMonitor")
	}
	if _, err := blockchain.CreateMonitor(ctx, project.Name, &blockchain.CreateMonitorReq{
		Contractlog: r.Contractlog,
		Chaintx:     r.Chaintx,
		ChainHeight: r.ChainHeight,
	}); err != nil {
		// TODO need a async task to keep server,monitor consistency
		l.Error(err)
		return m, nil
	}
	m.State = enums.MONITOR_STATE__SYNCED
	if err := m.UpdateByMonitorID(d); err != nil {
		l.Error(err)
		return nil, status.CheckDatabaseError(err, "UpdateByMonitorID")
	}
	return m, nil
}

func getMonitorType(r *CreateMonitorReq) (enums.MonitorType, error) {
	switch {
	case r.Contractlog != nil:
		return enums.MONITOR_TYPE__CONTRACT_LOG, nil
	case r.Chaintx != nil:
		return enums.MONITOR_TYPE__CHAIN_TX, nil
	case r.ChainHeight != nil:
		return enums.MONITOR_TYPE__CHAIN_HEIGHT, nil
	default:
		return enums.MONITOR_TYPE_UNKNOWN, status.BadRequest
	}
}

// TODO delete monitor
