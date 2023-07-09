.PHONY: lib
lib:
	sed -i 's/^int nmap_main.*/#ifdef __cplusplus\nextern "C" int nmap_main(int argc, char *argv[]);\nextern "C" void set_program_name(const char *name);\n#else\nextern int nmap_main(int argc, char *argv[]);\nextern void set_program_name(const char *name);\n#endif/g' src/nmap.h
	cd src; CFLAGS="-fPIC" CXXFLAGS="-fPIC" ./configure --without-zenmap
	make -C src STATIC="-shared"
	mv src/nmap lib/libnmap.so
	cd src; ar rcs libnmap.a *.o
	mv src/libnmap.a lib/
	ranlib lib/libnmap.a

