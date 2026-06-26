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

// Option configures a table created with New.
type Option func(*Table)

type rowSeparatorConfig struct {
	minRows int
	maxRows int
}

type Table struct {
	columns       []string
	templates     []*template.Template
	alignments    []Alignment
	data          []interface{}
	rowSeparators *rowSeparatorConfig
}

// New creates a table with the provided options.
func New(options ...Option) *Table {
	t := &Table{}

	for _, option := range options {
		option(t)
	}

	return t
}

// WithRowSeparators adds separator lines between row groups.
//
// Groups are sized between minRows and maxRows when possible. No separators are
// added if splitting would leave a group smaller than minRows.
func WithRowSeparators(minRows, maxRows int) Option {
	if minRows <= 0 {
		panic("table: minRows must be positive")
	}
	if maxRows < minRows {
		panic("table: maxRows must be greater than or equal to minRows")
	}

	return func(t *Table) {
		t.rowSeparators = &rowSeparatorConfig{
			minRows: minRows,
			maxRows: maxRows,
		}
	}
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

	separator := makeSeparatorLine(columnWidths)
	result.WriteString(separator)

	rowSeparatorBreaks := t.rowSeparatorBreaks(len(lines))
	for rowNumber, row := range lines {
		for i, column := range row {
			line[i] = alignText(column, columnWidths[i], t.alignments[i])
		}
		result.WriteString(strings.Join(line, "|") + "\n")

		if _, ok := rowSeparatorBreaks[rowNumber+1]; ok {
			result.WriteString(separator)
		}
	}

	return result.String()
}

func makeSeparatorLine(columnWidths []int) string {
	line := make([]string, len(columnWidths))
	for i, width := range columnWidths {
		line[i] = strings.Repeat("-", width+2)
	}

	return strings.Join(line, "+") + "\n"
}

func (t *Table) rowSeparatorBreaks(rows int) map[int]struct{} {
	if t.rowSeparators == nil {
		return nil
	}

	config := t.rowSeparators
	groupCount := (rows + config.maxRows - 1) / config.maxRows
	if groupCount < 2 || rows < groupCount*config.minRows {
		return nil
	}

	baseGroupSize := rows / groupCount
	extraRows := rows % groupCount
	breaks := make(map[int]struct{}, groupCount-1)

	rowNumber := 0
	for group := 0; group < groupCount-1; group++ {
		groupSize := baseGroupSize
		if group < extraRows {
			groupSize++
		}

		rowNumber += groupSize
		breaks[rowNumber] = struct{}{}
	}

	return breaks
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
