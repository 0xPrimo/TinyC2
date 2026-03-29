
all:
	make -C plugins/tcp
	make -C plugins/http
	go build -C ./server -o ../tinyc2

clean:
	make -C plugins/tcp clean
	make -C plugins/http clean
	rm -rf tinyc2
