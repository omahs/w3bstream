package monitor

import (
	"context"
	"time"

	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/pkg/errors"
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
		s := enums.MONITOR_STATE__SYNCED
		if m.DeletedAt.IsZero() {
			if err := createBlockchain(ctx, p, m.MonitorID, &m.Data); err != nil {
				l.Error(errors.Wrap(err, "create blockchain failed"))
				s = enums.MONITOR_STATE__FAILED_UNKNOWN
				if err == status.Conflict {
					s = enums.MONITOR_STATE__FAILED_CONFLICT
				}
			}
		} else {
			if err := removeBlockchain(ctx, p, m.MonitorID, m.Data.Type); err != nil {
				l.Error(errors.Wrap(err, "remove blockchain failed"))
				if err != status.NotFound {
					continue
				}
			}
		}
		m.State = s
		if err := m.UpdateByMonitorID(d); err != nil {
			l.Error(errors.Wrap(err, "update monitor db failed"))
		}
	}
}

func Sync(ctx context.Context) {
	s := &sync{syncInterval}
	go s.run(ctx)
}
