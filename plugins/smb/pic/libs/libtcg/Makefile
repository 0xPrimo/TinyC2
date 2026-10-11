CC=i686-w64-mingw32-gcc
CC_64=x86_64-w64-mingw32-gcc
CFLAGS=-O1 -fno-jump-tables -shared -Wall -Wno-pointer-arith

.PHONY: all x86 x64 clean

all: x86 x64

bin:
	mkdir -p bin

x86: bin
	$(CC) -DWIN_X86 $(CFLAGS) -c src/loaddll.c     -o bin/loaddll.x86.o
	$(CC) -DWIN_X86 $(CFLAGS) -c src/resolve_eat.c -o bin/resolve_eat.x86.o
	$(CC) -DWIN_X86 $(CFLAGS) -c src/picorun.c     -o bin/picorun.x86.o
	$(CC) -DWIN_X86 $(CFLAGS) -c src/debug.c       -o bin/debug.x86.o
	$(CC) -DWIN_X86 $(CFLAGS) -c src/util.c        -o bin/util.x86.o
	zip -q -j libtcg.x86.zip bin/*.x86.o

x64: bin
	$(CC_64) -DWIN_X64 $(CFLAGS) -c src/loaddll.c     -o bin/loaddll.x64.o
	$(CC_64) -DWIN_X64 $(CFLAGS) -c src/resolve_eat.c -o bin/resolve_eat.x64.o
	$(CC_64) -DWIN_X64 $(CFLAGS) -c src/picorun.c     -o bin/picorun.x64.o
	$(CC_64) -DWIN_X64 $(CFLAGS) -c src/debug.c       -o bin/debug.x64.o
	$(CC_64) -DWIN_X64 $(CFLAGS) -c src/util.c        -o bin/util.x64.o
	zip -q -j libtcg.x64.zip bin/*.x64.o

clean:
	rm -rf bin
	rm -f libtcg.x86.zip
	rm -f libtcg.x64.zip
