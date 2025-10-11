package pbar

import (
	"strings"
)

type ProgressBar struct {
	Bars []Bar

	stdoutWriter WriterWrapper
}

func (b *ProgressBar) AsString(width uint) string {
	formattedBars := make([]string, len(b.Bars))
	for i, bar := range b.Bars {
		formattedBars[i] = bar.asString(width)
	}

	return strings.Join(formattedBars, "\n")
}

func (b *ProgressBar) String() string {
	return b.AsString(b.stdoutWriter.Width())
}

func (b *ProgressBar) SetBarValue(id int, value uint) {
	b.Bars[id].Current = value
	_, _ = b.stdoutWriter.WriteRaw([]byte(b.String() + "\r"))
}

func (b *ProgressBar) Start() error {
	err := b.stdoutWriter.Start()
	if err != nil {
		return err
	}
	_, _ = b.stdoutWriter.WriteRaw([]byte(b.String() + "\r"))
	b.stdoutWriter.pbar = b
	return nil
}
