package pbar

import (
	"fmt"
	"strings"
)

type Bar struct {
	Current uint
	Total   uint
}

func (b Bar) asString(width uint) string {
	total := width - 7
	filled := b.Current * total / b.Total

	res := strings.Builder{}
	res.WriteString("[")
	res.WriteString(strings.Repeat("=", int(filled)))
	if filled > 0 && filled < total {
		res.WriteString(">")
		res.WriteString(strings.Repeat(" ", int(total-filled-1)))
	} else {
		res.WriteString(strings.Repeat(" ", int(total-filled)))
	}

	res.WriteString("] ")
	res.WriteString(fmt.Sprintf("%d%%", b.Current*100/b.Total))

	return res.String()
}
