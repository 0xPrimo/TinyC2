
all:
	make -C plugins/tcp
	make -C plugins/http
	make -C plugins/smb
	go build -C ./server -o ../tinyc2

clean:
	make -C plugins/tcp clean
	make -C plugins/http clean
	make -C plugins/smb clean

	rm -rf tinyc2
