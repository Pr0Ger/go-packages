package table

import (
	"bytes"
	"strings"
	"text/template"
)

// Alignment controls how text is aligned inside table columns.
type Alignment int

const (
	// AlignLeft aligns text to the left side of a column.
	AlignLeft Alignment = iota
	// AlignCenter centers text inside a column.
	AlignCenter
	// AlignRight aligns text to the right side of a column.
	AlignRight
)

type Table struct {
	columns    []string
	templates  []*template.Template
	alignments []Alignment
	data       []interface{}
}

func (t *Table) AddColumn(header, format string) {
	t.AddColumnAligned(header, format, AlignCenter)
}

// AddColumnAligned adds a column with the given text alignment.
func (t *Table) AddColumnAligned(header, format string, alignment Alignment) {
	t.columns = append(t.columns, header)

	tmpl, err := template.New("template for " + header).Parse(format)
	if err != nil {
		panic(err)
	}

	t.templates = append(t.templates, tmpl)
	t.alignments = append(t.alignments, alignment)
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
		line[i] = alignText(column, columnWidths[i], t.alignments[i])
	}
	result.WriteString(strings.Join(line, "|") + "\n")

	for i, width := range columnWidths {
		line[i] = strings.Repeat("-", width+2)
	}
	result.WriteString(strings.Join(line, "+") + "\n")

	for _, row := range lines {
		for i, column := range row {
			line[i] = alignText(column, columnWidths[i], t.alignments[i])
		}
		result.WriteString(strings.Join(line, "|") + "\n")
	}

	return result.String()
}

func alignText(text string, width int, alignment Alignment) string {
	padding := width - len(text)

	var leftPadding, rightPadding int
	switch alignment {
	case AlignLeft:
		rightPadding = padding
	case AlignCenter:
		leftPadding = padding / 2
		rightPadding = padding - leftPadding
	case AlignRight:
		leftPadding = padding
	}

	return strings.Repeat(" ", leftPadding+1) + text + strings.Repeat(" ", rightPadding+1)
}
