package blockchain

import (
	"context"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

const chainUniqFlag = 0

type (
	CreateContractLogReq = models.ContractLogInfo
	CreateChainTxReq     = models.ChainTxInfo
	CreateChainHeightReq = models.ChainHeightInfo
)

func CreateContractLog(ctx context.Context, projectName string, id types.SFID, r *CreateContractLogReq) (*models.ContractLog, error) {
	d := types.MustMonitorDBExecutorFromContext(ctx)

	if err := checkChainID(d, r.ChainID); err != nil {
		return nil, err
	}

	n := *r
	n.BlockCurrent = n.BlockStart
	n.EventType = getEventType(n.EventType)
	m := &models.ContractLog{
		RelContractLog: models.RelContractLog{ContractLogID: id},
		ContractLogData: models.ContractLogData{
			ProjectName:     projectName,
			Uniq:            chainUniqFlag,
			ContractLogInfo: n,
		},
	}
	if err := m.Create(d); err != nil {
		return nil, status.CheckDatabaseError(err, "CreateContractLog")
	}
	return m, nil
}

func CreateChainTx(ctx context.Context, projectName string, id types.SFID, r *CreateChainTxReq) (*models.ChainTx, error) {
	d := types.MustMonitorDBExecutorFromContext(ctx)

	if err := checkChainID(d, r.ChainID); err != nil {
		return nil, err
	}

	n := *r
	n.EventType = getEventType(n.EventType)
	m := &models.ChainTx{
		RelChainTx: models.RelChainTx{ChainTxID: id},
		ChainTxData: models.ChainTxData{
			ProjectName: projectName,
			Uniq:        chainUniqFlag,
			Finished:    datatypes.FALSE,
			ChainTxInfo: n,
		},
	}
	if err := m.Create(d); err != nil {
		return nil, status.CheckDatabaseError(err, "CreateChainTx")
	}
	return m, nil
}

func CreateChainHeight(ctx context.Context, projectName string, id types.SFID, r *CreateChainHeightReq) (*models.ChainHeight, error) {
	d := types.MustMonitorDBExecutorFromContext(ctx)

	if err := checkChainID(d, r.ChainID); err != nil {
		return nil, err
	}

	n := *r
	n.EventType = getEventType(n.EventType)
	m := &models.ChainHeight{
		RelChainHeight: models.RelChainHeight{ChainHeightID: id},
		ChainHeightData: models.ChainHeightData{
			ProjectName:     projectName,
			Uniq:            chainUniqFlag,
			Finished:        datatypes.FALSE,
			ChainHeightInfo: n,
		},
	}
	if err := m.Create(d); err != nil {
		return nil, status.CheckDatabaseError(err, "CreateChainHeight")
	}
	return m, nil
}

func checkChainID(d sqlx.DBExecutor, id uint64) error {
	b := &models.Blockchain{RelBlockchain: models.RelBlockchain{ChainID: id}}
	if err := b.FetchByChainID(d); err != nil {
		return status.CheckDatabaseError(err, "GetBlockchainByChainID")
	}
	return nil
}

func RemoveContractLog(ctx context.Context, projectName string, id types.SFID) error {
	d := types.MustMonitorDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "RemoveContractLog")
	defer l.End()

	l = l.WithValues("project", projectName)

	m := &models.ContractLog{RelContractLog: models.RelContractLog{ContractLogID: id}}
	if err := m.FetchByContractLogID(d); err != nil {
		return status.CheckDatabaseError(err, "FetchByContractLogID")
	}
	if err := checkProjectName(m.ProjectName, projectName, l); err != nil {
		return err
	}
	if err := m.DeleteByContractLogID(d); err != nil {
		return status.CheckDatabaseError(err, "DeleteByContractLogID")
	}
	return nil
}

func RemoveChainTx(ctx context.Context, projectName string, id types.SFID) error {
	d := types.MustMonitorDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "RemoveChainTx")
	defer l.End()

	l = l.WithValues("project", projectName)

	m := &models.ChainTx{RelChainTx: models.RelChainTx{ChainTxID: id}}
	if err := m.FetchByChainTxID(d); err != nil {
		return status.CheckDatabaseError(err, "FetchByChainTxID")
	}
	if err := checkProjectName(m.ProjectName, projectName, l); err != nil {
		return err
	}
	if err := m.DeleteByChainTxID(d); err != nil {
		return status.CheckDatabaseError(err, "DeleteByChainTxID")
	}
	return nil
}

func RemoveChainHeight(ctx context.Context, projectName string, id types.SFID) error {
	d := types.MustMonitorDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "RemoveChainHeight")
	defer l.End()

	l = l.WithValues("project", projectName)

	m := &models.ChainHeight{RelChainHeight: models.RelChainHeight{ChainHeightID: id}}
	if err := m.FetchByChainHeightID(d); err != nil {
		return status.CheckDatabaseError(err, "FetchByChainHeightID")
	}
	if err := checkProjectName(m.ProjectName, projectName, l); err != nil {
		return err
	}
	if err := m.DeleteByChainHeightID(d); err != nil {
		return status.CheckDatabaseError(err, "DeleteByChainHeightID")
	}
	return nil
}

func checkProjectName(want, curr string, l log.Logger) error {
	if want != curr {
		l.Error(errors.New("monitor project mismatch"))
		return status.BadRequest.StatusErr().WithDesc("monitor project mismatch")
	}
	return nil
}

func getEventType(eventType string) string {
	if eventType == "" {
		return enums.MONITOR_EVENTTYPEDEFAULT
	}
	return eventType
}
