package receiver

import (
	"context"
	"log"
	"net"
	"sync"
)

type Receiver interface {
	Start(ctx context.Context)
}

type tcpReceiver struct {
	host, port string
}

func NewTCPReceiver(host, port string) *tcpReceiver {
	rv := &tcpReceiver{host: host, port: port}
	return rv
}

func (c *tcpReceiver) Start(
	ctx context.Context,
	wg *sync.WaitGroup,
	next func(conn net.Conn),
) error {
	lc := net.ListenConfig{}
	l, err := lc.Listen(ctx, "tcp", net.JoinHostPort(c.host, c.port))
	if err != nil {
		return err
	}
	wg.Add(1)
	go func() {
		defer log.Print("listen done")
		defer wg.Done()
		defer l.Close()
		<-ctx.Done()
	}()
	ctx, cancel := context.WithCancelCause(ctx)
	wg.Add(1)
	go func() {
		defer log.Print("accept done")
		defer wg.Done()
		defer cancel(nil)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			log.Print("accept...")
			conn, err := l.Accept()
			if err != nil {
				cancel(err)
				return
			}
			next(conn)
		}
	}()
	return nil
}
