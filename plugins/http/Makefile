all:
	go build -buildmode=plugin -o ./bin/http.plugin.so .
	make -C pic

clean:
	rm -rf bin
	make -C pic clean