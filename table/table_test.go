package table

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTable(t *testing.T) {
	tests := []struct {
		create func() *Table
		output string
	}{
		{
			func() *Table {
				tb := &Table{}
				tb.AddColumn("ID", "{{ .ID }}")
				tb.AddColumn("Name", "{{ .Name }}")

				tb.AddRow(struct{ ID, Name string }{"id", "name"})

				return tb
			},
			`
 ID | Name 
----+------
 id | name 
`,
		},
		{
			func() *Table {
				tb := &Table{}
				tb.AddColumn("ID", "{{ .ID }}")
				tb.AddColumn("Name", "{{ .Name }}")

				tb.AddRow(struct{ ID, Name string }{"really long id to test alignment", "name"})

				return tb
			},
			`
                ID                | Name 
----------------------------------+------
 really long id to test alignment | name 
`,
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			table := test.create()
			want := strings.TrimLeft(test.output, "\n")
			assert.Equal(t, want, table.String())
		})
	}
}
