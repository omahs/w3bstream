package monitor

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

const (
	syncInterval = 3 * time.Second
)

type sync struct {
	interval time.Duration
}

func (s *sync) run(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for range ticker.C {
		s.do(ctx)
	}
}

func (s *sync) do(ctx context.Context) {
	d := types.MustMonitorDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	m := &models.Monitor{}

	_, l = l.Start(ctx, "monitor.sync.run")
	defer l.End()

	// TODO how to list soft delete data
	ms, err := m.List(d, m.ColState().Eq(enums.MONITOR_STATE__SYNCING))
	if err != nil {
		l.Error(errors.Wrap(err, "list monitor db failed"))
		return
	}
	for _, m := range ms {
		p := &models.Project{RelProject: models.RelProject{ProjectID: m.ProjectID}}
		if err := p.FetchByProjectID(d); err != nil {
			l.Error(errors.Wrap(err, "fetch project db failed"))
			continue
		}

		if m.DeletedAt.IsZero() {
			s := enums.MONITOR_STATE__SYNCED
			if err := createBlockchain(ctx, p, m.MonitorID, &m.Data); err != nil {
				l.Error(errors.Wrap(err, "create blockchain failed"))
				s = enums.MONITOR_STATE__FAILED_UNKNOWN
				if sqlx.DBErr(err).IsConflict() {
					s = enums.MONITOR_STATE__FAILED_CONFLICT
				}
			}
			m.State = s
			if err := m.UpdateByMonitorID(d); err != nil {
				l.Error(errors.Wrap(err, "update monitor db failed"))
			}
		} else {
			if err := removeBlockchain(ctx, p, m.MonitorID, m.Data.Type); err != nil {
				l.Error(errors.Wrap(err, "remove blockchain failed"))
				if !sqlx.DBErr(err).IsNotFound() {
					continue
				}
			}
			if err := m.DeleteByMonitorID(d); err != nil {
				l.Error(errors.Wrap(err, "delete monitor db failed"))
			}
		}
	}
}

func Sync(ctx context.Context) {
	s := &sync{syncInterval}
	go s.run(ctx)
}
