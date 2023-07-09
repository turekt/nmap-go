//go:build !windows

package nmap

import (
	"bytes"
	"io"
	"os"

	"golang.org/x/sys/unix"
)

func CaptureOut(exec func()) (string, error) {
	stdout, err := unix.Dup(unix.Stdout)
	if err != nil {
		return "", err
	}

	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}

	if err := unix.Dup2(int(w.Fd()), unix.Stdout); err != nil {
		return "", err
	}
	exec()
	w.Close()

	if err := unix.Dup2(stdout, unix.Stdout); err != nil {
		return "", err
	}
	if err := unix.Close(stdout); err != nil {
		return "", err
	}

	var buf bytes.Buffer
	io.Copy(&buf, r)
	r.Close()

	// terminal loses echo so we restore it here
	f, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		// no point in failing here?
		return buf.String(), nil
	}

	fd := int(f.Fd())
	term, err := unix.IoctlGetTermios(fd, unix.TCGETS)
	if err != nil {
		return buf.String(), nil
	}

	term.Lflag |= unix.ECHO | unix.ICANON
	unix.IoctlSetTermios(fd, unix.TCSETS, term)
	return buf.String(), nil
}
