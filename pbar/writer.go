package pbar

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"golang.org/x/term"
)

type WriterWrapper struct {
	originalStdout *os.File
	pipeWriter     *os.File

	pbar *ProgressBar
}

func (ww *WriterWrapper) Width() uint {
	terminal := os.Stdout
	if ww.originalStdout != nil {
		terminal = ww.originalStdout
	}

	interactive := term.IsTerminal(int(terminal.Fd()))
	if !interactive {
		return 0
	}

	width, _, err := term.GetSize(int(terminal.Fd()))
	if err != nil {
		return 0
	}
	if width == 0 {
		width = 80
	}
	return uint(min(width, 80))
}

func (ww *WriterWrapper) Write(p []byte) (n int, err error) {
	var buf bytes.Buffer

	for _, b := range p {
		if b == '\n' {
			buf.WriteString(EraseLine)
		}
		buf.WriteByte(b)
	}

	buf.WriteString(ww.pbar.String() + "\r")

	n, err = ww.WriteRaw(buf.Bytes())
	if err != nil {
		return n, err
	}

	return len(p), nil
}

func (ww *WriterWrapper) WriteRaw(p []byte) (n int, err error) {
	return ww.originalStdout.Write(p)
}

func (ww *WriterWrapper) Start() error {
	ww.originalStdout = os.Stdout

	r, w, err := os.Pipe()
	if err != nil {
		return fmt.Errorf("could not create pipe: %w", err)
	}

	ww.pipeWriter = w
	os.Stdout = w

	go func() {
		defer r.Close()

		if _, err := io.Copy(ww, r); err != nil {
			log.Printf("Error copying from pipe to decorator: %v\n", err)
		}
	}()

	return nil
}

func (ww *WriterWrapper) Stop() error {
	if ww.originalStdout != nil {
		os.Stdout = ww.originalStdout
		ww.originalStdout = nil
	}
	if ww.pipeWriter != nil {
		ww.pipeWriter.Close()
		ww.pipeWriter = nil
	}
	return nil
}
