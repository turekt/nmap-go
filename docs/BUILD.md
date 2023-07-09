# Building nmap as library

Currently nmap does not offer a Makefile target to build it as a library. The approach taken for nmap-go was to offer an interface towards C code that:
- is simple
- does not require a lot of changes in the original nmap C code (none would be optimal)
- is futureproof for any further changes made in nmap C code

Because of this, lib simply interfaces with `nmap_main` function that takes an `argc` and `argv` parameter (https://github.com/nmap/nmap/blob/7ae5c4d9278b981e10bbd8b3877587bfccff163d/nmap.cc#L1817). This keeps the Go lib:
- simple: usual nmap command line switches work
- does not require a lot of changes to the original nmap code: only one function is labeled with `extern`
- is futureproof: I would not expect much changes to happen with regards to `nmap_main` as opposed to interfacing with other internal functions

There are two other approaches that come to mind and one might find them "more programmatic":
1. interface with scan functions directly (e.g. `script_scan` or `ultra_scan`) but this would be worst in terms of the wanted approach
  - more functions should be labeled as `extern`
  - implementing the lib and using it would not be as simple
  - seems much less futureproof
2. change the `nmap_main` function signature to return `std::vector<Target *> Targets` that contain target details but this requires much more changes to the nmap C code

The only issue with calling `nmap_main` instead of `main` is that `nmap_main` will not check for `--resume` in scans so this has to be done manually. For now this seems a good tradeoff to avoid interfacing directly with `int main(argc, argv)`.

This opens up several other problems down the line:
- Which changes should be made to the original C code? (extern)
- How do we return target information if `nmap_main` returns `int`?

Nevertheless, to simply build nmap as a shared object (nmap as libnmap.so) we need to do:
```
CFLAGS="-fPIC" CXXFLAGS="-fPIC" ./configure --without-zenmap
make STATIC="-shared"
```

For static library, we have to create an archive from all of the generated object files:
```
ar rcs libnmap.a *.o
ranlib libnmap.a
```
