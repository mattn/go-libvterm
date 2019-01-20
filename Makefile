libvterm.a : libvterm
	cd libvterm && CFLAGS="-O3" DEBUG=0 make libvterm.la
	cp libvterm/.libs/libvterm.a .
	strip --strip-unneeded --strip-debug ./libvterm.a

clean:
	cd libvterm && make clean
	rm -f libvterm.a
