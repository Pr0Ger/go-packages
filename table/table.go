package table

import (
	"bytes"
	"strings"
	"text/template"
)

type Table struct {
	columns   []string
	templates []*template.Template
	data      []interface{}
}

func (t *Table) AddColumn(header, format string) {
	t.columns = append(t.columns, header)

	tmpl, err := template.New("template for " + header).Parse(format)
	if err != nil {
		panic(err)
	}

	t.templates = append(t.templates, tmpl)
}

func (t *Table) AddRow(data interface{}) {
	t.data = append(t.data, data)
}

func (t *Table) String() string {
	columnWidths := make([]int, len(t.columns))
	for i, column := range t.columns {
		columnWidths[i] = len(column)
	}

	lines := make([][]string, 0, len(t.data))
	buf := bytes.NewBuffer(nil)

	for _, data := range t.data {
		row := make([]string, len(t.templates))
		for i, tmpl := range t.templates {
			err := tmpl.Execute(buf, data)
			if err != nil {
				panic(err)
			}

			row[i] = buf.String()
			buf.Reset()

			if len(row[i]) > columnWidths[i] {
				columnWidths[i] = len(row[i])
			}
		}
		lines = append(lines, row)
	}

	result := strings.Builder{}

	line := make([]string, len(t.columns))
	for i, column := range t.columns {
		leftPadding := (columnWidths[i] - len(column)) / 2
		rightPadding := columnWidths[i] - len(column) - leftPadding
		line[i] = strings.Repeat(" ", leftPadding+1) + column + strings.Repeat(" ", rightPadding+1)
	}
	result.WriteString(strings.Join(line, "|") + "\n")

	for i, width := range columnWidths {
		line[i] = strings.Repeat("-", width+2)
	}
	result.WriteString(strings.Join(line, "+") + "\n")

	for _, row := range lines {
		for i, column := range row {
			leftPadding := (columnWidths[i] - len(column)) / 2
			rightPadding := columnWidths[i] - len(column) - leftPadding
			line[i] = strings.Repeat(" ", leftPadding+1) + column + strings.Repeat(" ", rightPadding+1)
		}
		result.WriteString(strings.Join(line, "|") + "\n")
	}

	return result.String()
}
