// This is a generated source file. DO NOT EDIT
// Source: models/customer__generated.go

package models

import (
	"fmt"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
)

var CustomerTable *builder.Table

func init() {
	CustomerTable = DB.Register(&Customer{})
}

type CustomerIterator struct {
}

func (*CustomerIterator) New() interface{} {
	return &Customer{}
}

func (*CustomerIterator) Resolve(v interface{}) *Customer {
	return v.(*Customer)
}

func (*Customer) TableName() string {
	return "t_customer"
}

func (*Customer) TableDesc() []string {
	return []string{
		"Customer stream demo",
	}
}

func (*Customer) Comments() map[string]string {
	return map[string]string{}
}

func (*Customer) ColDesc() map[string][]string {
	return map[string][]string{}
}

func (*Customer) ColRel() map[string][]string {
	return map[string][]string{}
}

func (*Customer) PrimaryKey() []string {
	return []string{
		"ID",
	}
}

func (m *Customer) IndexFieldNames() []string {
	return []string{
		"ID",
	}
}

func (m *Customer) ColID() *builder.Column {
	return CustomerTable.ColByFieldName(m.FieldID())
}

func (*Customer) FieldID() string {
	return "ID"
}

func (m *Customer) ColFirstName() *builder.Column {
	return CustomerTable.ColByFieldName(m.FieldFirstName())
}

func (*Customer) FieldFirstName() string {
	return "FirstName"
}

func (m *Customer) ColLastName() *builder.Column {
	return CustomerTable.ColByFieldName(m.FieldLastName())
}

func (*Customer) FieldLastName() string {
	return "LastName"
}

func (m *Customer) ColAge() *builder.Column {
	return CustomerTable.ColByFieldName(m.FieldAge())
}

func (*Customer) FieldAge() string {
	return "Age"
}

func (m *Customer) ColTaxNumber() *builder.Column {
	return CustomerTable.ColByFieldName(m.FieldTaxNumber())
}

func (*Customer) FieldTaxNumber() string {
	return "TaxNumber"
}

func (m *Customer) ColCity() *builder.Column {
	return CustomerTable.ColByFieldName(m.FieldCity())
}

func (*Customer) FieldCity() string {
	return "City"
}

func (m *Customer) ColCreatedAt() *builder.Column {
	return CustomerTable.ColByFieldName(m.FieldCreatedAt())
}

func (*Customer) FieldCreatedAt() string {
	return "CreatedAt"
}

func (m *Customer) ColUpdatedAt() *builder.Column {
	return CustomerTable.ColByFieldName(m.FieldUpdatedAt())
}

func (*Customer) FieldUpdatedAt() string {
	return "UpdatedAt"
}

func (m *Customer) ColDeletedAt() *builder.Column {
	return CustomerTable.ColByFieldName(m.FieldDeletedAt())
}

func (*Customer) FieldDeletedAt() string {
	return "DeletedAt"
}

func (m *Customer) CondByValue(db sqlx.DBExecutor) builder.SqlCondition {
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

func (m *Customer) Create(db sqlx.DBExecutor) error {

	if m.CreatedAt.IsZero() {
		m.CreatedAt.Set(time.Now())
	}

	if m.UpdatedAt.IsZero() {
		m.UpdatedAt.Set(time.Now())
	}

	_, err := db.Exec(sqlx.InsertToDB(db, m, nil))
	return err
}

func (m *Customer) List(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) ([]Customer, error) {
	var (
		tbl = db.T(m)
		lst = make([]Customer, 0)
	)
	cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("Customer.List")}, adds...)
	err := db.QueryAndScan(builder.Select(nil).From(tbl, adds...), &lst)
	return lst, err
}

func (m *Customer) Count(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) (cnt int64, err error) {
	tbl := db.T(m)
	cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("Customer.List")}, adds...)
	err = db.QueryAndScan(builder.Select(builder.Count()).From(tbl, adds...), &cnt)
	return
}

func (m *Customer) FetchByID(db sqlx.DBExecutor) error {
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
				builder.Comment("Customer.FetchByID"),
			),
		m,
	)
	return err
}

func (m *Customer) UpdateByIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

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
				builder.Comment("Customer.UpdateByIDWithFVs"),
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

func (m *Customer) UpdateByID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByIDWithFVs(db, fvs)
}

func (m *Customer) Delete(db sqlx.DBExecutor) error {
	_, err := db.Exec(
		builder.Delete().
			From(
				db.T(m),
				builder.Where(m.CondByValue(db)),
				builder.Comment("Customer.Delete"),
			),
	)
	return err
}

func (m *Customer) DeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("Customer.DeleteByID"),
			),
	)
	return err
}

func (m *Customer) SoftDeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("Customer.SoftDeleteByID"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}
