package consumer

import (
	"bufio"
	"context"
	"log"
	"net"
	"sync"
)

type Consumer interface {
	Start(ctx context.Context)
}

type tcpConsumer struct {
	conn net.Conn
}

func NewTcpReceiver(conn net.Conn) *tcpConsumer {
	rv := &tcpConsumer{conn}
	return rv
}

func (c *tcpConsumer) Start(
	ctx context.Context,
	wg *sync.WaitGroup,
	next func(b []byte) (int, error),
) error {
	wg.Add(1)
	go func() {
		defer log.Print("consumer done")
		defer wg.Done()
		defer func() {
			_ = c.conn.Close()
		}()
		buffer := bufio.NewReader(c.conn)
		readBuf := make([]byte, 1<<10)
		ctx, cancel := context.WithCancelCause(ctx)
		defer cancel(nil)

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			_, err2 := buffer.Read(readBuf)
			if err2 != nil {
				cancel(err2)
			}
			_, err := next(readBuf)
			if err != nil {
				cancel(err)
				return
			}
		}
	}()
	return nil
}
