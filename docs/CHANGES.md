# Changes needed to interface with `nmap_main`

Since parts of nmap are written in C++, interfacing with `nmap_main` was not so simple as building the project with `-fPIC` and `-shared` flags. What C++ does after building the library with appropriate flags is **name mangling**. Meaning that the function and variable names are encoded into unique function name in order to facilitate function overloading and visibility.

This makes interfacing harder and there is a need to label `nmap_main` as `extern`, which prevents name mangling happening on that function.

We can do this with `sed` by replacing the function signature in `nmap.h` and adding `extern`. Since nmap utilizes both C and C++, the change should include adding `extern` for both languages with `#ifdef` directive.

This is what we are changing in `nmap.h`:
```
index 276bf2367..a58bd6ccf 100644
--- a/nmap.h
+++ b/nmap.h
@@ -259,7 +259,13 @@
 /***********************PROTOTYPES**********************************/
 
 /* Renamed main so that interactive mode could preprocess when necessary */
-int nmap_main(int argc, char *argv[]);
+#ifdef __cplusplus
+extern "C" int nmap_main(int argc, char *argv[]);
+extern "C" void set_program_name(const char *name);
+#else
+extern int nmap_main(int argc, char *argv[]);
+extern void set_program_name(const char *name);
+#endif
 
 int nmap_fetchfile(char *filename_returned, int bufferlen, const char *file);
 int gather_logfile_resumption_state(char *fname, int *myargc, char ***myargv);
```

The extra function `set_program_name` is unfortunately needed because nmap will assert that `program_name` variable is set and `nmap_main` expects it to be set before the function call is made: https://github.com/nmap/nmap/blob/7ae5c4d9278b981e10bbd8b3877587bfccff163d/nmap.cc#L2676.
