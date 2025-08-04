package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os/signal"
	"slices"
	"sync"
	"syscall"
	"time"

	"github.com/mascanio/uroboros/internal/generator/syslog"
	"github.com/mascanio/uroboros/internal/receiver"
	"github.com/mascanio/uroboros/internal/sender"
)

func main() {
	// Start pprof server
	go func() {
		log.Println("pprof listening on :6060")
		http.ListenAndServe(":6060", nil)
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGKILL, syscall.SIGABRT)
	defer cancel()

	wg := sync.WaitGroup{}

	receiver := receiver.NewTCPReceiver("localhost", "4444")
	receiverCtx, cancelReceiver := context.WithCancel(ctx)
	defer cancelReceiver()
	err := receiver.Start(receiverCtx, &wg, func(conn net.Conn) {
		wg.Add(1)
		go consumer(&wg, conn)
	})
	if err != nil {
		log.Fatal(err)
	}

	wg.Add(1)
	go func() {
		defer log.Print("sender done")
		defer wg.Done()
		dialer := net.Dialer{}
		sender, err := sender.NewTCPSender(ctx, sender.NewDialerRetry(
			&dialer,
			"tcp",
			"localhost",
			"4444",
			sender.RetryNWait(33, 100*time.Millisecond),
		))
		if err != nil {
			log.Fatal(err)
		}
		defer sender.Close()

		// gen := io.LimitReader(&r{ctx}, 1<<20)
		gen := syslog.NewSyslogGenerator(syslog.RFC5424, 100000000)
		r := bufio.NewReaderSize(gen, 1<<14)
		w := bufio.NewWriterSize(sender.Writer, 1<<14)
		_, err = w.ReadFrom(r)
		if err != nil {
			log.Fatal(err)
		}
	}()

	wg.Wait()
}

func consumer(wg *sync.WaitGroup, conn net.Conn) {
	nRead := 0
	lastNRead := 0
	lastTime := time.Now()

	log.Print("consumer started")
	defer log.Print("consumer done")
	defer wg.Done()
	defer conn.Close()
	buf := make([]byte, 1<<12)
	buffReader := bufio.NewReader(conn)
	for {
		n, readErr := buffReader.Read(buf)
		nRead += n
		if errors.Is(readErr, io.EOF) {
			log.Print(ByteCountIEC(int64(nRead) - int64(lastNRead)))
			log.Print("total read: ", ByteCountIEC(int64(nRead)))
			return
		}
		if time.Since(lastTime) > time.Second {
			lastTime = time.Now()
			log.Print(ByteCountIEC(int64(nRead) - int64(lastNRead)))
			lastNRead = nRead
		}
	}
}

type r struct{ ctx context.Context }

var repeated = slices.Repeat([]byte("a"), 1<<30)

func (r *r) Read(p []byte) (n int, err error) {
	select {
	case <-r.ctx.Done():
		return 0, context.Cause(r.ctx)
	default:
	}
	return copy(p, repeated[:len(p)]), nil
}

func ByteCountIEC(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}
