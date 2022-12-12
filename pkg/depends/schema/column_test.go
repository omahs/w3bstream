package schema_test

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/schema"
)

func TestColumn_Ex(t *testing.T) {
	ctx := context.Background()
	col := schema.Col("f_column")

	ex := col.Ex(ctx)
	NewWithT(t).Expect(ex.Query()).To(Equal("f_column"))

	col.WithTable(schema.T("t_tab"))
	ex = col.Ex(builder.ContextWithToggleMultiTable(ctx, true))
	NewWithT(t).Expect(ex.Query()).To(Equal("t_tab.f_column"))

	ex = col.Ex(builder.ContextWithToggles(ctx, builder.Toggles{
		builder.ToggleNeedAutoAlias: true,
		builder.ToggleMultiTable:    true,
	}))
	NewWithT(t).Expect(ex.Query()).To(Equal("t_tab.f_column AS f_column"))
}
