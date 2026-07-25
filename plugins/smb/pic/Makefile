CYRSTAL_LINKER	= crystal-link
CC_64			= x86_64-w64-mingw32-gcc
CFLAGS			= -O1 -fno-jump-tables -shared -Wall -Wno-pointer-arith

.PHONY: all x64 clean

all: x64

bin:
	mkdir -p bin

x64: bin
	$(CC_64) -DWIN_X64 $(CFLAGS) -Iinclude -c src/services.c -o bin/services.x64.o
	$(CC_64) -DWIN_X64 $(CFLAGS) -Iinclude -c src/channel.c -o bin/channel.x64.o
	$(CYRSTAL_LINKER) $(CURDIR)/channel.spec $(CURDIR)/bin/channel.x64.o $(CURDIR)/../bin/channel.x64.pic

clean:
	rm -rf bin
