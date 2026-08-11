all:
	go build -buildmode=plugin -o ./build/smb .
	make -C extension

clean:
	rm -rf build
	make -C extension clean
