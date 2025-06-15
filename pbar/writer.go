package pbar

import (
	"fmt"
	"io"
	"log"
	"os"
)

type WriterWrapper struct {
	originalStdout *os.File
	pipeWriter     *os.File
}

func (ww *WriterWrapper) Write(p []byte) (n int, err error) {
	return ww.originalStdout.Write(p)
}

func (ww *WriterWrapper) Start() error {
	ww.originalStdout = os.Stdout

	r, w, err := os.Pipe()
	if err != nil {
		return fmt.Errorf("could not create pipe: %v", err)
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
