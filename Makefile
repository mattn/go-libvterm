libvterm.a: libvterm
	cd libvterm && \
	CFLAGS="-O3 -fstack-protector-strong -march=x86-64 -mtune=generic" DEBUG=0 make && \
	cd .. && \
	cp ./libvterm/.libs/libvterm.a .
	strip --strip-unneeded --strip-debug ./libvterm.a

clean:
	cd libvterm && \
	make clean && \
	cd .. && \
	rm -f libvterm.a
