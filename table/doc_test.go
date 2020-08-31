package table

import (
	"fmt"
)

func ExampleTable_String() {
	table := Table{}

	table.AddColumn("ID", "{{ .ID }}")
	table.AddColumn("Name", "{{ .Name }}")

	type row struct {
		ID   string
		Name string
	}

	table.AddRow(row{
		ID:   "0",
		Name: "Test",
	})

	fmt.Println(table.String())
}
