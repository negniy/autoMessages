package logger

import (
	"io"
	"sync"
)

type CustomWriter struct {
	outputs []io.Writer
	mu      sync.Mutex
}

func (cw *CustomWriter) Write(p []byte) (n int, err error) {
	cw.mu.Lock()
	defer cw.mu.Unlock()
	for _, w := range cw.outputs {
		n, err = w.Write(p)
		if err != nil {
			return n, err
		}
	}
	return len(p), nil
}

func (cw *CustomWriter) AddOutput(w io.Writer) {
	cw.mu.Lock()
	defer cw.mu.Unlock()
	cw.outputs = append(cw.outputs, w)
}
