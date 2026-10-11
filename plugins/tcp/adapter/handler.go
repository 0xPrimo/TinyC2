package adapter

import (
	"context"
	"encoding/binary"
	"io"
	"net"
	"sync"
)

func (l *Listener) handler(conn net.Conn) {
	defer func() {
		l.connsMu.Lock()
		if l.conns != nil {
			delete(l.conns, conn)
		}
		l.connsMu.Unlock()
		conn.Close()
		l.wg.Done()
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	queue := make(chan []byte, 32)
	var tasks sync.WaitGroup

	// writer goroutine
	l.wg.Add(1)
	go l.writeLoop(ctx, conn, queue, cancel)

	// reader loop
	l.readLoop(ctx, conn, queue, &tasks)

	// cleanup
	cancel()
	tasks.Wait()
	close(queue)
}

func (l *Listener) writeLoop(ctx context.Context, conn net.Conn, queue <-chan []byte, cancel context.CancelFunc) {
	defer l.wg.Done()
	defer cancel()

	head := make([]byte, 4)
	for {
		select {
		case <-ctx.Done():
			return
		case data, ok := <-queue:
			if !ok {
				return
			}
			binary.BigEndian.PutUint32(head, uint32(len(data)))
			if _, err := conn.Write(head); err != nil {
				return
			}
			if _, err := conn.Write(data); err != nil {
				return
			}
		}
	}
}

func (l *Listener) readLoop(ctx context.Context, conn net.Conn, queue chan<- []byte, tasks *sync.WaitGroup) {
	head := make([]byte, 4)
	for {
		// check for context cancellation
		select {
		case <-ctx.Done():
			return
		default:
		}

		if _, err := io.ReadFull(conn, head); err != nil {
			return
		}

		// read the frame size
		size := binary.BigEndian.Uint32(head)
		if size > 4096*5 {
			l.LogError("tcp", "Frame size %d exceeds limit of %d bytes", size, 4096*5)
			return
		}

		body := make([]byte, size)
		if _, err := io.ReadFull(conn, body); err != nil {
			return
		}

		// process the frame
		tasks.Add(1)
		go func(data []byte) {
			defer tasks.Done()

			reply, err := l.ImplantProcess(l.name, data)
			if err != nil {
				l.LogError("tcp", "Failed to process implant request: %v", err)
				select {
				case queue <- []byte("error"):
				case <-ctx.Done():
				}
				return
			}

			select {
			case queue <- reply:
			case <-ctx.Done():
			}
		}(body)
	}
}
