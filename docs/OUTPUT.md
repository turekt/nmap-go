# Returning target information from `nmap_main`

A call to `nmap_main` runs nmap in a similar manner as from the command line: arguments are parsed, scan is initiated and target information is printed to the terminal.

Since `nmap_main` returns only the exit code, we have no other way to extract target information other than from terminal output which becomes tricky once you understand that Go's `os.Stdin`, `os.Stdout` and `os.Stderr` will not affect C's stdio `stdin`, `stdout` and `stderr`. Redirecting these standard files with `os.Pipe` and resetting Go's `os.Stdout` will not make any difference and extracting target information that way will tremendously fail. Instead, the solution is to go much deeper on the OS level (syscall level) to change stdio's `stdout` and forward terminal output back to Go.

On unix system we utilize `dup2` call to write stdout output into a write pipe created in Go. We read the stdout output through the read pipe.

On windows system there is an existing `SetStdHandle` function that is used to reset the stdout to a handle created in Go. This is similar to unix, just a different function call.
