package schema_test

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/pkg/depends/schema"
	"github.com/machinefi/w3bstream/pkg/depends/x/ptrx"
)

func TestColumnType_DatabaseDatatype(t *testing.T) {
	cases := []*struct {
		Name       string
		Constrains *schema.ColumnType
		Expected   string
	}{
		{
			Name:       "Integer",
			Constrains: &schema.ColumnType{Datatype: schema.DATATYPE__UINT},
			Expected:   "integer",
		}, {
			Name:       "IntSerial",
			Constrains: &schema.ColumnType{AutoIncrement: true, Datatype: schema.DATATYPE__INT},
			Expected:   "serial",
		}, {
			Name:       "Bigint",
			Constrains: &schema.ColumnType{Datatype: schema.DATATYPE__INT64},
			Expected:   "bigint",
		}, {
			Name:       "BigSerial",
			Constrains: &schema.ColumnType{AutoIncrement: true, Datatype: schema.DATATYPE__UINT64},
			Expected:   "bigserial",
		}, {
			Name:       "Real",
			Constrains: &schema.ColumnType{Datatype: schema.DATATYPE__FLOAT32},
			Expected:   "real",
		}, {
			Name:       "DoublePrecision",
			Constrains: &schema.ColumnType{Datatype: schema.DATATYPE__FLOAT64},
			Expected:   "double precision",
		}, {
			Name:       "Boolean",
			Constrains: &schema.ColumnType{Datatype: schema.DATATYPE__BOOL},
			Expected:   "boolean",
		}, {
			Name:       "Text",
			Constrains: &schema.ColumnType{Datatype: schema.DATATYPE__TEXT, Length: 65536},
			Expected:   "text",
		}, {
			Name:       "Timestamp",
			Constrains: &schema.ColumnType{Datatype: schema.DATATYPE__TIMESTAMP},
			Expected:   "bigint",
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			NewWithT(t).Expect(c.Constrains.DatabaseDatatype()).To(Equal(c.Expected))
		})
	}
}

func TestColumnType_AutoCompleteLengthDatatype(t *testing.T) {
	cases := []*struct {
		Name       string
		Constrains *schema.ColumnType
		Expected   string
	}{
		{
			Name:       "DoublePrecisionWithDecimal",
			Constrains: &schema.ColumnType{Datatype: schema.DATATYPE__FLOAT64, Decimal: 3, Length: 100},
			Expected:   "double precision(100,3)",
		}, {
			Name:       "CharacterVaryingWithDefaultLength",
			Constrains: &schema.ColumnType{Datatype: schema.DATATYPE__TEXT},
			Expected:   "character varying(255)",
		}, {
			Name:       "CharacterVaryingWithLength",
			Constrains: &schema.ColumnType{Datatype: schema.DATATYPE__TEXT, Length: 1024},
			Expected:   "character varying(1024)",
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			NewWithT(t).Expect(c.Constrains.AutoCompleteLengthDatatype()).To(Equal(c.Expected))
		})
	}
}

func TestColumnType_TypeModify(t *testing.T) {
	cases := []*struct {
		Name       string
		Constrains *schema.ColumnType
		Expected   string
	}{
		{
			Name:       "TextWithDefaultValue",
			Constrains: &schema.ColumnType{Datatype: schema.DATATYPE__TEXT, Default: ptrx.Ptr("''"), Length: 65536},
			Expected:   " NOT NULL DEFAULT ''::text",
		}, {
			Name:       "IntegerWithDefaultValue",
			Constrains: &schema.ColumnType{Datatype: schema.DATATYPE__INT, Default: ptrx.Ptr("'0'")},
			Expected:   " NOT NULL DEFAULT '0'::integer",
		}, {
			Name:       "BigintWithDefaultValue",
			Constrains: &schema.ColumnType{Datatype: schema.DATATYPE__INT64, Default: ptrx.Ptr("'0'")},
			Expected:   " NOT NULL DEFAULT '0'::bigint",
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			NewWithT(t).Expect(c.Constrains.TypeModify()).To(Equal(c.Expected))
		})
	}
}
