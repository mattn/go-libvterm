VERSION=0.2

libvterm.a: libvterm
	cd libvterm && CFLAGS="-O3" DEBUG=0 make libvterm.la
	strip --strip-unneeded --strip-debug libvterm/.libs/libvterm.a

libvterm:
	curl -L -s https://launchpad.net/libvterm/trunk/v$(VERSION)/+download/libvterm-$(VERSION).tar.gz |\
		tar xz
	mv libvterm-$(VERSION) libvterm

clean:
	cd libvterm && make clean
