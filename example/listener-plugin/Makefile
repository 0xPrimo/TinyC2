
all:
	go build -buildmode=plugin -o ./bin/protocol.plugin.so .
	make -C pic

clean:
	rm -rf ./bin
	make -C pic clean