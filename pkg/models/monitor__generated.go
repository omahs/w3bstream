// This is a generated source file. DO NOT EDIT
// Source: models/monitor__generated.go

package models

import (
	"fmt"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
)

var MonitorTable *builder.Table

func init() {
	MonitorTable = DB.Register(&Monitor{})
}

type MonitorIterator struct {
}

func (*MonitorIterator) New() interface{} {
	return &Monitor{}
}

func (*MonitorIterator) Resolve(v interface{}) *Monitor {
	return v.(*Monitor)
}

func (*Monitor) TableName() string {
	return "t_monitor"
}

func (*Monitor) TableDesc() []string {
	return []string{
		"Monitor project monitor info",
	}
}

func (*Monitor) Comments() map[string]string {
	return map[string]string{}
}

func (*Monitor) ColDesc() map[string][]string {
	return map[string][]string{}
}

func (*Monitor) ColRel() map[string][]string {
	return map[string][]string{}
}

func (*Monitor) PrimaryKey() []string {
	return []string{
		"ID",
	}
}

func (m *Monitor) IndexFieldNames() []string {
	return []string{
		"ID",
		"MonitorID",
	}
}

func (*Monitor) UniqueIndexes() builder.Indexes {
	return builder.Indexes{
		"ui_monitor_id": []string{
			"MonitorID",
			"DeletedAt",
		},
	}
}

func (*Monitor) UniqueIndexUIMonitorID() string {
	return "ui_monitor_id"
}

func (m *Monitor) ColID() *builder.Column {
	return MonitorTable.ColByFieldName(m.FieldID())
}

func (*Monitor) FieldID() string {
	return "ID"
}

func (m *Monitor) ColMonitorID() *builder.Column {
	return MonitorTable.ColByFieldName(m.FieldMonitorID())
}

func (*Monitor) FieldMonitorID() string {
	return "MonitorID"
}

func (m *Monitor) ColProjectID() *builder.Column {
	return MonitorTable.ColByFieldName(m.FieldProjectID())
}

func (*Monitor) FieldProjectID() string {
	return "ProjectID"
}

func (m *Monitor) ColState() *builder.Column {
	return MonitorTable.ColByFieldName(m.FieldState())
}

func (*Monitor) FieldState() string {
	return "State"
}

func (m *Monitor) ColData() *builder.Column {
	return MonitorTable.ColByFieldName(m.FieldData())
}

func (*Monitor) FieldData() string {
	return "Data"
}

func (m *Monitor) ColCreatedAt() *builder.Column {
	return MonitorTable.ColByFieldName(m.FieldCreatedAt())
}

func (*Monitor) FieldCreatedAt() string {
	return "CreatedAt"
}

func (m *Monitor) ColUpdatedAt() *builder.Column {
	return MonitorTable.ColByFieldName(m.FieldUpdatedAt())
}

func (*Monitor) FieldUpdatedAt() string {
	return "UpdatedAt"
}

func (m *Monitor) ColDeletedAt() *builder.Column {
	return MonitorTable.ColByFieldName(m.FieldDeletedAt())
}

func (*Monitor) FieldDeletedAt() string {
	return "DeletedAt"
}

func (m *Monitor) CondByValue(db sqlx.DBExecutor) builder.SqlCondition {
	var (
		tbl  = db.T(m)
		fvs  = builder.FieldValueFromStructByNoneZero(m)
		cond = []builder.SqlCondition{tbl.ColByFieldName("DeletedAt").Eq(0)}
	)

	for _, fn := range m.IndexFieldNames() {
		if v, ok := fvs[fn]; ok {
			cond = append(cond, tbl.ColByFieldName(fn).Eq(v))
			delete(fvs, fn)
		}
	}
	if len(cond) == 0 {
		panic(fmt.Errorf("no field for indexes has value"))
	}
	for fn, v := range fvs {
		cond = append(cond, tbl.ColByFieldName(fn).Eq(v))
	}
	return builder.And(cond...)
}

func (m *Monitor) Create(db sqlx.DBExecutor) error {

	if m.CreatedAt.IsZero() {
		m.CreatedAt.Set(time.Now())
	}

	if m.UpdatedAt.IsZero() {
		m.UpdatedAt.Set(time.Now())
	}

	_, err := db.Exec(sqlx.InsertToDB(db, m, nil))
	return err
}

func (m *Monitor) List(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) ([]Monitor, error) {
	var (
		tbl = db.T(m)
		lst = make([]Monitor, 0)
	)
	cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("Monitor.List")}, adds...)
	err := db.QueryAndScan(builder.Select(nil).From(tbl, adds...), &lst)
	return lst, err
}

func (m *Monitor) Count(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) (cnt int64, err error) {
	tbl := db.T(m)
	cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("Monitor.List")}, adds...)
	err = db.QueryAndScan(builder.Select(builder.Count()).From(tbl, adds...), &cnt)
	return
}

func (m *Monitor) FetchByID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ID").Eq(m.ID),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("Monitor.FetchByID"),
			),
		m,
	)
	return err
}

func (m *Monitor) FetchByMonitorID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("MonitorID").Eq(m.MonitorID),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("Monitor.FetchByMonitorID"),
			),
		m,
	)
	return err
}

func (m *Monitor) UpdateByIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("ID").Eq(m.ID),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("Monitor.UpdateByIDWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByID(db)
	}
	return nil
}

func (m *Monitor) UpdateByID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByIDWithFVs(db, fvs)
}

func (m *Monitor) UpdateByMonitorIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("MonitorID").Eq(m.MonitorID),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("Monitor.UpdateByMonitorIDWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByMonitorID(db)
	}
	return nil
}

func (m *Monitor) UpdateByMonitorID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByMonitorIDWithFVs(db, fvs)
}

func (m *Monitor) Delete(db sqlx.DBExecutor) error {
	_, err := db.Exec(
		builder.Delete().
			From(
				db.T(m),
				builder.Where(m.CondByValue(db)),
				builder.Comment("Monitor.Delete"),
			),
	)
	return err
}

func (m *Monitor) DeleteByID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ID").Eq(m.ID),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("Monitor.DeleteByID"),
			),
	)
	return err
}

func (m *Monitor) SoftDeleteByID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	fvs := builder.FieldValues{}

	if _, ok := fvs["DeletedAt"]; !ok {
		fvs["DeletedAt"] = types.Timestamp{Time: time.Now()}
	}

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	_, err := db.Exec(
		builder.Update(db.T(m)).
			Where(
				builder.And(
					tbl.ColByFieldName("ID").Eq(m.ID),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("Monitor.SoftDeleteByID"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}

func (m *Monitor) DeleteByMonitorID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("MonitorID").Eq(m.MonitorID),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("Monitor.DeleteByMonitorID"),
			),
	)
	return err
}

func (m *Monitor) SoftDeleteByMonitorID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	fvs := builder.FieldValues{}

	if _, ok := fvs["DeletedAt"]; !ok {
		fvs["DeletedAt"] = types.Timestamp{Time: time.Now()}
	}

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	_, err := db.Exec(
		builder.Update(db.T(m)).
			Where(
				builder.And(
					tbl.ColByFieldName("MonitorID").Eq(m.MonitorID),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("Monitor.SoftDeleteByMonitorID"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}
