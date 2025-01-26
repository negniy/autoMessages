package logger

import (
	"io"
	"strings"
	"sync"

	"fyne.io/fyne/widget"
)

type CustomWriter struct {
	outputs   []io.Writer
	logWidget *widget.Entry
	mu        sync.Mutex
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

	if cw.logWidget != nil {
		cw.logWidget.SetText(cw.logWidget.Text + string(p))
		cw.logWidget.CursorRow = len(strings.Split(cw.logWidget.Text, "\n")) - 1
		cw.logWidget.Refresh()
	}

	return len(p), nil
}

func (cw *CustomWriter) AddOutput(w io.Writer) {
	cw.mu.Lock()
	defer cw.mu.Unlock()
	cw.outputs = append(cw.outputs, w)
}

func (cw *CustomWriter) SetLogWidget(widget *widget.Entry) {
	cw.mu.Lock()
	defer cw.mu.Unlock()
	cw.logWidget = widget
}
