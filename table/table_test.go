package table

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testRow struct {
	ID string
}

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
		{
			func() *Table {
				tb := &Table{}
				tb.AddColumnAligned("ID", "{{ .ID }}", AlignLeft)
				tb.AddColumnAligned("Status", "{{ .Status }}", AlignCenter)
				tb.AddColumnAligned("Count", "{{ .Count }}", AlignRight)

				tb.AddRow(struct{ ID, Status, Count string }{"1", "ok", "7"})
				tb.AddRow(struct{ ID, Status, Count string }{"longer", "warn", "120"})

				return tb
			},
			`
 ID     | Status | Count 
--------+--------+-------
 1      |   ok   |     7 
 longer |  warn  |   120 
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

func TestTableWithRowSeparators(t *testing.T) {
	table := New(WithRowSeparators(2, 3))
	table.AddColumn("ID", "{{ .ID }}")

	for i := 1; i <= 4; i++ {
		table.AddRow(testRow{ID: strconv.Itoa(i)})
	}

	want := strings.TrimLeft(`
 ID 
----
 1  
 2  
----
 3  
 4  
`, "\n")
	assert.Equal(t, want, table.String())
}

func TestTableWithRowSeparatorsKeepsSmallRemainderTogether(t *testing.T) {
	table := New(WithRowSeparators(4, 5))
	table.AddColumn("ID", "{{ .ID }}")

	for i := 1; i <= 6; i++ {
		table.AddRow(testRow{ID: strconv.Itoa(i)})
	}

	want := strings.TrimLeft(`
 ID 
----
 1  
 2  
 3  
 4  
 5  
 6  
`, "\n")
	assert.Equal(t, want, table.String())
}
