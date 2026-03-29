package listener

import (
	"encoding/binary"
	"io"
	"net"
	"sync"
)

func TcpHandler(listener *TcpListener, conn net.Conn) {
	var tasks sync.WaitGroup

	defer listener.wg.Done()
	defer conn.Close()

	queue := make(chan []byte, 32)

	listener.wg.Add(1)

	// write frames
	go func() {
		defer listener.wg.Done()
		defer conn.Close()

		head := make([]byte, 4)
		for data := range queue {
			binary.BigEndian.PutUint32(head, uint32(len(data)))

			if _, err := conn.Write(head); err != nil {
				return
			}
			if _, err := conn.Write(data); err != nil {
				return
			}
		}
	}()

	// read frames
	head := make([]byte, 4)
	for {
		if _, err := io.ReadFull(conn, head); err != nil {
			break
		}

		size := binary.BigEndian.Uint32(head)
		if size > 4096*5 {
			break
		}

		body := make([]byte, size)
		if _, err := io.ReadFull(conn, body); err != nil {
			break
		}

		// process request
		tasks.Add(1)
		go func(data []byte) {
			defer tasks.Done()

			reply, err := listener.Engine.ImplantProcess(listener.Name, data)
			if err != nil {
				queue <- []byte("error")
				return
			}

			queue <- reply
		}(body)
	}

	tasks.Wait()
	close(queue)
}
