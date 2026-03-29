all:
	go build -buildmode=plugin -o ./bin/tcp.plugin.so .
	make -C pic

clean:
	rm -rf bin
	make -C pic clean