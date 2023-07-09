package nmap

/*
#cgo CFLAGS: -I./src
#cgo LDFLAGS: -L${SRCDIR}/lib -l:libnmap.so
#include <stdlib.h>
#include <nmap.h>

static void nmap_main_wrapper(int argc, char** argv) {
	set_program_name(argv[0]);
	nmap_main(argc, argv);
}
*/
import "C"

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"unsafe"
)

//go:embed src/nmap-mac-prefixes
//go:embed src/nmap-os-db
//go:embed src/nmap-protocols
//go:embed src/nmap-rpc
//go:embed src/nmap-service-probes
//go:embed src/nmap-services
//go:embed src/*.lua
//go:embed src/scripts
//go:embed src/nselib
var content embed.FS

const FlagDataDir = "--datadir"

type NmapOptions struct {
	ProgramName string
	Args        string
}

func Scan(args string) (string, error) {
	return ScanWithOptions(&NmapOptions{"nmap", args})
}

func ScanWithOptions(opt *NmapOptions) (string, error) {
	opts := strings.Split(opt.Args, " ")
	dataDirSet := false
	for _, s := range opts {
		if s == FlagDataDir {
			dataDirSet = true
			break
		}
	}

	args := []string{opt.ProgramName}
	if !dataDirSet {
		dir, err := os.MkdirTemp("", "cfgs")
		if err != nil {
			return "", err
		}
		defer os.RemoveAll(dir)

		if err := fs.WalkDir(content, ".", func(path string, f fs.DirEntry, err error) error {
			nonsrcPath := strings.TrimPrefix(path, "src/")
			if f.IsDir() {
				os.Mkdir(filepath.Join(dir, nonsrcPath), 0755)
				return nil
			}

			data, err := content.ReadFile(path)
			if err != nil {
				return err
			}

			newPath := filepath.Join(dir, nonsrcPath)
			if err := os.WriteFile(newPath, data, 0644); err != nil {
				return err
			}

			return nil
		}); err != nil {
			return "", err
		}

		args = append(args, FlagDataDir)
		args = append(args, dir)
	}

	args = append(args, opts...)
	argc := C.int(len(args))
	argv := make([]*C.char, len(args))

	for i, s := range args {
		argv[i] = C.CString(s)
		defer C.free(unsafe.Pointer(argv[i]))
	}

	return CaptureOut(func() {
		C.nmap_main_wrapper(argc, &argv[0])
	})
}
