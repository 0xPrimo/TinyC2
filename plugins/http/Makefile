all:
	go build -buildmode=plugin -o ./build/http .
	make -C extension

clean:
	rm -rf bin
	make -C extension clean