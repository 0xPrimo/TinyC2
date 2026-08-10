all:
	go build -buildmode=plugin -o ./build/tcp .
	make -C extension

clean:
	rm -rf build
	make -C extension clean