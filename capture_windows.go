//go:build windows

package windows

import (
	"bytes"
	"io"
	"os"

	"golang.org/x/sys/windows"
)

func CaptureOut(exec func()) (string, error) {
	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}

	if err := windows.SetStdHandle(windows.STD_OUTPUT_HANDLE, windows.Handle(w.Fd())); err != nil {
		return "", err
	}

	exec()
	w.Close()

	if err := windows.SetStdHandle(windows.STD_OUTPUT_HANDLE, windows.Stdout); err != nil {
		return "", err
	}

	var buf bytes.Buffer
	io.Copy(&buf, r)
	r.Close()
	return buf.String(), nil
}
