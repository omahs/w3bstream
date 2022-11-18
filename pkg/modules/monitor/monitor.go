package monitor

import (
	"context"
	"time"

	"github.com/pkg/errors"

	bt "github.com/machinefi/w3bstream/pkg/depends/base/types"
	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/blockchain"
	"github.com/machinefi/w3bstream/pkg/types"
)

type (
	CreateMonitorReq = models.MonitorInfo
)

func Create(ctx context.Context, project *models.Project, r *CreateMonitorReq) (*models.Monitor, error) {
	d := types.MustDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)

	_, l = l.Start(ctx, "MonitorCreate")
	defer l.End()

	id := idg.MustGenSFID()

	l = l.WithValues("monitor_id", id, "project_name", project.Name)

	m := &models.Monitor{
		RelMonitor: models.RelMonitor{MonitorID: id},
		RelProject: models.RelProject{ProjectID: project.ProjectID},
		MonitorData: models.MonitorData{
			State: enums.MONITOR_STATE__SYNCING,
			Data:  *r,
		},
	}
	if err := m.Create(d); err != nil {
		l.Error(err)
		return nil, status.CheckDatabaseError(err, "CreateMonitor")
	}

	err := createBlockchain(ctx, project, id, r)
	m.State = enums.MONITOR_STATE__SYNCED
	if err != nil {
		l.Error(err)
		m.State = enums.MONITOR_STATE__FAILED_UNKNOWN
		if err == status.Conflict {
			m.State = enums.MONITOR_STATE__FAILED_CONFLICT
		}
	}

	if err := m.UpdateByMonitorID(d); err != nil {
		l.Error(err)
		return nil, status.CheckDatabaseError(err, "UpdateByMonitorID")
	}
	return m, nil
}

func Remove(ctx context.Context, project *models.Project, id types.SFID, t enums.MonitorType) error {
	d := types.MustDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "MonitorRemove")
	defer l.End()

	l = l.WithValues("monitor_id", id, "project_name", project.Name)

	m := &models.Monitor{RelMonitor: models.RelMonitor{MonitorID: id}}
	if err := m.FetchByMonitorID(d); err != nil {
		return status.CheckDatabaseError(err, "FetchByMonitorID")
	}
	if project.ProjectID != m.ProjectID {
		l.Error(errors.New("monitor project mismatch"))
		return status.BadRequest.StatusErr().WithDesc("monitor project mismatch")
	}
	m.DeletedAt = bt.AsTimestamp(time.Now())
	m.State = enums.MONITOR_STATE__SYNCING
	if err := m.UpdateByMonitorID(d); err != nil {
		l.Error(err)
		return status.CheckDatabaseError(err, "UpdateByMonitorID")
	}

	if err := removeBlockchain(ctx, project, id, t); err != nil && !sqlx.DBErr(err).IsNotFound() {
		l.Error(err)
		return nil
	}

	if err := m.DeleteByMonitorID(d); err != nil {
		l.Error(err)
		return status.CheckDatabaseError(err, "DeleteByMonitorID")
	}
	return nil
}

func createBlockchain(ctx context.Context, project *models.Project, id types.SFID, r *CreateMonitorReq) error {
	var err error
	switch r.Type {
	case enums.MONITOR_TYPE__CONTRACT_LOG:
		_, err = blockchain.CreateContractLog(ctx, project.Name, id, r.ContractLog)
	case enums.MONITOR_TYPE__CHAIN_TX:
		_, err = blockchain.CreateChainTx(ctx, project.Name, id, r.ChainTx)
	case enums.MONITOR_TYPE__CHAIN_HEIGHT:
		_, err = blockchain.CreateChainHeight(ctx, project.Name, id, r.ChainHeight)
	}
	return err
}

func removeBlockchain(ctx context.Context, project *models.Project, id types.SFID, t enums.MonitorType) error {
	var err error
	switch t {
	case enums.MONITOR_TYPE__CONTRACT_LOG:
		err = blockchain.RemoveContractLog(ctx, project.Name, id)
	case enums.MONITOR_TYPE__CHAIN_TX:
		err = blockchain.RemoveChainTx(ctx, project.Name, id)
	case enums.MONITOR_TYPE__CHAIN_HEIGHT:
		err = blockchain.RemoveChainHeight(ctx, project.Name, id)
	}
	return err
}
