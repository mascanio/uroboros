package main

import (
	"bufio"
	"context"
	"fmt"
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

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGKILL,
		syscall.SIGABRT,
	)
	defer cancel()

	wg := sync.WaitGroup{}

	receiver := receiver.NewTCPReceiver("localhost", "4444")
	receiverCtx, cancelReceiver := context.WithCancel(ctx)
	defer cancelReceiver()
	err := receiver.Start(receiverCtx, &wg, func(conn net.Conn) {
		wg.Add(1)
		go consumerRFC(&wg, conn, syslog.RFC5424) // Change to RFC3164 to test BSD
	})
	if err != nil {
		log.Fatal(err)
	}

	wg.Add(1)
	go func() {
		defer log.Print("sender done")
		defer wg.Done()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
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
		go func() {
			defer sender.Close()
			<-ctx.Done()
		}()

		// gen := io.LimitReader(&r{ctx}, 1<<20)
		gen := syslog.NewSyslogGenerator(
			syslog.RFC5424,
			syslog.WithEndOfLine([]byte("\n")),
		)
		w := bufio.NewWriterSize(sender.Writer, 1<<14)
		defer w.Flush()
		buf := make([]byte, 1<<10)
		msg := []byte("This is a test syslog message ")
		for range 10000000 {
			_, err := gen.GenerateMessage(buf, time.Now(), msg)
			if err != nil {
				return
			}
			w.Write(buf)
		}
	}()

	wg.Wait()
}

// consumerRFC is a generic consumer for both RFCs
func consumerRFC(wg *sync.WaitGroup, conn net.Conn, rfc syslog.Format) {
	nRead := 0
	lastNRead := 0
	lastTime := time.Now()

	log.Print("consumer started")
	defer log.Print("consumer done")
	defer wg.Done()
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	scanner.Split(syslog.NewSyslogSplitFuncDelimiter([]byte("\n")))
	for scanner.Scan() {
		nRead += len(scanner.Bytes())
		_ = syslog.ParseSyslogMessageRFC5424(scanner.Bytes())
		// log.Printf("Received syslog message: %+v", parsed)
		if time.Since(lastTime) > time.Second {
			lastTime = time.Now()
			log.Print(ByteCountIEC(int64(nRead) - int64(lastNRead)))
			lastNRead = nRead
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Scanner error: %v", err)
	}
	log.Print("total read: ", ByteCountIEC(int64(nRead)))
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
