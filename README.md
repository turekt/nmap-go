# nmap-go - Go nmap lib using CGo bindings

This project is an nmap library for Go that calls nmap via CGo bindings rather than installing `nmap` on a machine and executing it via `os/exec`.

This approach might be more suitable in cases where:
- you are unable to install nmap on a machine
- you want to limit your program's external dependencies
- running a second process is not really optimal for your use case

## Building the lib

On Linux build nmap as a library:
```
make lib
```

This creates two files inside `lib/` directory:
- `libnmap.so`- dynamically linked library file (shared object)
- `libnmap.a` - statically linked library file (archive)

The lib uses `libnmap.a` by default.

## Usage

We try to build a simple Go program that runs a server and scans the port with nmap:
```go
package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/turekt/nmap-go"
)

func main() {
	server := &http.Server{Addr: "127.0.0.1:8080"}
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	str, err := nmap.Scan("-sC -sV -p 8080 127.0.0.1")
	if err != nil {
		panic(err)
	}
	fmt.Println("Test output")
	fmt.Printf("%s\n", str)
	server.Shutdown(context.Background())
}
```

### Building using libnmap.a - statically linked program

Static program do not require any external dependencies when being run. To build the program as a static binary, `ldflags` should point to nmap compiled dependencies (created after `make lib` is executed):
```
$ NMAP_SRCDIR=../nmap-go/src; go build -tags osusergo,netgo -ldflags "-extldflags \"-static -L${NMAP_SRCDIR}/nbase -L${NMAP_SRCDIR}/nsock/src/ -lstdc++ -lnsock -lnbase ${NMAP_SRCDIR}/libpcre/.libs/libpcre2-8.a ${NMAP_SRCDIR}/liblua/liblua.a -lm ${NMAP_SRCDIR}/libssh2/lib/libssh2.a -lssl -lcrypto ${NMAP_SRCDIR}/libnetutil/libnetutil.a ${NMAP_SRCDIR}/libdnet-stripped/src/.libs/libdnet.a ${NMAP_SRCDIR}/libpcap/libpcap.a ${NMAP_SRCDIR}/libz/libz.a ${NMAP_SRCDIR}/liblinear/liblinear.a\"" main.go
```

Make sure to set `NMAP_SRCDIR` to point to your nmap source. Build might output warnings but these can be ignored.

Compiled program works without any external dependencies:
```
$ file main
main: ELF 64-bit LSB executable, x86-64, version 1 (GNU/Linux), statically linked, BuildID[sha1]=bf14a3bc5cd02940da6011566dc33df5e4c1e917, for GNU/Linux 3.2.0, with debug_info, not stripped
$ ./main
Test output
Starting Nmap 7.94SVN ( https://nmap.org ) at 2023-07-09 13:27 CEST
Nmap scan report for localhost (127.0.0.1)
Host is up (0.00013s latency).

PORT     STATE SERVICE VERSION
8080/tcp open  http    Golang net/http server (Go-IPFS json-rpc or InfluxDB API)
|_http-open-proxy: Proxy might be redirecting requests
|_http-title: Site doesn't have a title (text/plain; charset=utf-8).

Service detection performed. Please report any incorrect results at https://nmap.org/submit/ .
Nmap done: 1 IP address (1 host up) scanned in 6.68 seconds
```

### Building using libnmap.so - dynamically linked program

In order to use `libnmap.so` on Linux, `nmap.go`'s `LDFLAGS` should be changed to use `libnmap.so` instead of `libnmap.a`:
```
#cgo LDFLAGS: -L${SRCDIR}/lib -l:libnmap.so
```

Build binary with go:
```
$ go build main.go
$ file main
main: ELF 64-bit LSB executable, x86-64, version 1 (SYSV), dynamically linked, interpreter /lib64/ld-linux-x86-64.so.2, BuildID[sha1]=749cdf2969900444ec5ea7857463f7ef48bf4236, for GNU/Linux 3.2.0, with debug_info, not stripped
```

In order for your program to work, `libnmap.so` has to be available on the device:
```
$ ./main
./main: error while loading shared libraries: libnmap.so: cannot open shared object file: No such file or directory
$ cp ../nmap-go/lib/libnmap.so /usr/local/lib
$ ldconfig
$ ./main 
Test output
Starting Nmap 7.94SVN ( https://nmap.org ) at 2023-07-09 13:20 CEST
Nmap scan report for localhost (127.0.0.1)
Host is up (0.00013s latency).

PORT     STATE SERVICE VERSION
8080/tcp open  http    Golang net/http server (Go-IPFS json-rpc or InfluxDB API)
|_http-title: Site doesn't have a title (text/plain; charset=utf-8).
|_http-open-proxy: Proxy might be redirecting requests

Service detection performed. Please report any incorrect results at https://nmap.org/submit/ .
Nmap done: 1 IP address (1 host up) scanned in 6.63 seconds
```

## nmap-go on Windows

The lib contains implementation for running and fetching execution output but has not been thoroughly tested. The approach should be the same: build the lib on windows (DLL and/or LIB) and compile your program against it.

Contributions are welcome.
