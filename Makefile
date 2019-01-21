libvterm.a: libvterm
	cd libvterm && CFLAGS="-O3" DEBUG=0 make libvterm.la
	strip --strip-unneeded --strip-debug libvterm/.libs/libvterm.a

libvterm:
	curl -s http://bazaar.launchpad.net/~libvterm/libvterm/trunk/tarball/head: |\
		tar xz --transform "s!^~libvterm/libvterm/trunk!libvterm!"

clean:
	cd libvterm && make clean
