package sender

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type tcpSender struct {
	io.Writer
	io.Closer
	buffer *bufio.Writer
}

type Dialer interface {
	DialContext(ctx context.Context) (net.Conn, error)
}

type retryDialer struct {
	dialer              *net.Dialer
	network, host, port string
	retryFn             RetryFn
}

type RetryFn func(ctx context.Context, do func() error) error

func RetryNWait(nRetries int, waitTime time.Duration) RetryFn {
	if nRetries <= 0 {
		panic(errors.New("nRetries should be > 0"))
	}
	return func(ctx context.Context, do func() error) error {
		var err error
		for range nRetries + 1 {
			err = do()
			if err == nil {
				return nil
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.Tick(waitTime):
			}
		}
		return fmt.Errorf("retries exhausted (total %v), last error: %v", nRetries, err)
	}
}

func NewDialerRetry(
	dialer *net.Dialer,
	network, host, port string,
	retryFn RetryFn,
) *retryDialer {
	return &retryDialer{
		dialer:  dialer,
		network: network,
		host:    host,
		port:    port,
		retryFn: retryFn,
	}
}

func (d *retryDialer) DialContext(ctx context.Context) (net.Conn, error) {
	var rv net.Conn
	err := d.retryFn(ctx, func() error {
		var errDial error
		rv, errDial = d.dialer.DialContext(ctx, d.network, net.JoinHostPort(d.host, d.port))
		return errDial
	})
	if err != nil {
		return nil, err
	}
	log.Print("connected")
	return rv, nil
}

func NewTCPSender(ctx context.Context, dialer Dialer) (*tcpSender, error) {
	const defaultBufSize = 4096
	return NewTCPSenderSize(ctx, dialer, defaultBufSize)
}

func NewTCPSenderSize(ctx context.Context, dialer Dialer, bufSize int) (*tcpSender, error) {
	var conn net.Conn
	var err error
	conn, err = dialer.DialContext(ctx)
	if err != nil {
		log.Print("error dialing: ", err)
		return nil, err
	}
	buffer := bufio.NewWriterSize(conn, bufSize)

	return &tcpSender{
		Writer: buffer,
		Closer: conn,
		buffer: buffer,
	}, nil
}

func (s *tcpSender) Close() error {
	flushErr := s.buffer.Flush()
	closeErr := s.Closer.Close()
	return errors.Join(flushErr, closeErr)
}

// func (s *tcpSender) Start(
// 	ctx context.Context,
// 	wg *sync.WaitGroup,
// 	next func(w io.Writer) (int, error),
// ) error {
// 	dialer := net.Dialer{}
// 	c, err := dialer.DialContext(
// 		ctx,
// 		"tcp",
// 		net.JoinHostPort(s.host, s.port),
// 	)
// 	if err != nil {
// 		return err
// 	}
// 	buffer := bufio.NewWriter(c)
// 	wg.Add(1)
// 	go func() {
// 		defer log.Print("sender done")
// 		defer wg.Done()
// 		defer func() {
// 			_ = buffer.Flush()
// 			_ = c.Close()
// 		}()
// 		ctx, cancel := context.WithCancelCause(ctx)
// 		defer cancel(nil)
// 		for {
// 			select {
// 			case <-ctx.Done():
// 				return
// 			default:
// 			}
// 			_, err := next(buffer)
// 			if err != nil {
// 				cancel(err)
// 				return
// 			}
// 		}
// 	}()
// 	return nil
// }
