package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	_ "net/http/pprof"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/mascanio/uroboros/internal/event_generator/syslog"
	syslogparser "github.com/mascanio/uroboros/internal/parser/syslog"
	payloadgenerator "github.com/mascanio/uroboros/internal/payload_generator"
	"github.com/mascanio/uroboros/internal/receiver"
	"github.com/mascanio/uroboros/internal/sender"
	sequencegenerator "github.com/mascanio/uroboros/internal/sequence_generator"
	"golang.org/x/time/rate"
)

func main() {
	// Command-line flags
	senderFlag := flag.Bool("sender", false, "Enable sender")
	receiverFlag := flag.Bool("receiver", false, "Enable receiver")
	flag.Parse()

	if !*senderFlag && !*receiverFlag {
		fmt.Println("Usage: uroboros [-sender] [-receiver]")
		fmt.Println("  -sender     Enable sender")
		fmt.Println("  -receiver   Enable receiver")
		return
	}

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGKILL,
		syscall.SIGABRT,
	)
	defer cancel()

	wg := sync.WaitGroup{}

	if *receiverFlag {
		setupReceivers(ctx, &wg, []receiverConfig{{
			host:   "localhost",
			port:   "4444",
			format: syslog.RFC5424,
		}})
	}

	if *senderFlag {
		setupSenders(ctx, &wg)
	}

	wg.Wait()
}

func setupSenders(ctx context.Context, wg *sync.WaitGroup) {
	sequenceGenerator := sequencegenerator.NewSequenceGenerator(0, 100_000_000)
	for range 1 {
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

			gen := syslog.NewSyslogGenerator(
				syslog.RFC5424,
				syslog.WithEndOfLine([]byte("\n")),
			)
			// payloadGenerator := payloadgenerator.NewRandomMessageGenerator(
			// 	payloadgenerator.WithSequenceGenerator(sequenceGenerator),
			// 	payloadgenerator.WithFixedLength(300),
			// 	// payloadgenerator.WithMinMaxLength(500, 900),
			// 	payloadgenerator.WithIncludeSeqInMsg(true),
			// )
			payloadGenerator := payloadgenerator.NewIDMessageGenerator(
				"seq: 0000166354, thread: 0000, runid: 1755194303, stamp: 2025-08-14T19:58:25 PADDPADDPADDPADDPADDPADDPADDPADDPADDPADDPADDPADDPADDPADDPADDPADDPADDPADDPADDPADDPADDPADDPADDPADDPADDPADDPADDPADDPADDPADDPADDPAD",
				sequenceGenerator,
			)
			limit := rate.NewLimiter(600, 1)
			w := bufio.NewWriterSize(sender.Writer, 1<<10)
			defer w.Flush()
			for {
				msg := payloadGenerator.GenerateMessage()
				if msg == nil {
					return
				}
				if limit.Tokens() < 1.0 {
					w.Flush()
				}
				err = limit.Wait(ctx)
				if err != nil {
					return
				}
				_, err := gen.GenerateEvent(w, time.Now(), msg)
				if err != nil {
					return
				}
			}
		}()
	}
}

type receiverConfig struct {
	host, port string
	format     syslog.Format
}

func setupReceivers(
	ctx context.Context,
	wg *sync.WaitGroup,
	cfg []receiverConfig,
) (context.CancelCauseFunc, error) {
	allReceiverCtx, cancelReceivers := context.WithCancelCause(ctx)
	for _, config := range cfg {
		receiver := receiver.NewTCPReceiver(config.host, config.port)
		receiverCtx := context.WithValue(allReceiverCtx, "host", config.host)
		receiverCtx = context.WithValue(receiverCtx, "port", config.port)
		err := receiver.Start(
			receiverCtx,
			wg,
			func(ctx context.Context, conn net.Conn) {
				consumerRFC(ctx, conn, config.format) // Change to RFC3164 to test BSD
			})
		if err != nil {
			cancelReceivers(err)
			return nil, err
		}
	}
	return cancelReceivers, nil
}

// consumerRFC is a generic consumer for both RFCs
func consumerRFC(ctx context.Context, conn net.Conn, rfc syslog.Format) {
	nRead := 0
	lastNRead := 0
	nEvents := 0
	lastNEvents := 0
	lastTime := time.Now()

	log.Print("consumer started")
	defer log.Print("consumer done")
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	scanner.Split(syslog.NewSyslogSplitFuncDelimiter([]byte("\n")))
	for scanner.Scan() {
		nRead += len(scanner.Bytes())
		nEvents++
		_ = syslogparser.ParseSyslogMessageRFC5424(scanner.Bytes())
		// parsed := syslogparser.ParseSyslogMessageRFC5424(scanner.Bytes())
		// log.Printf("Received syslog message: %s, %s", parsed.Msg, string(scanner.Bytes()))
		if time.Since(lastTime) > time.Second {
			lastTime = time.Now()
			bytesPerSec := nRead - lastNRead
			eventsPerSec := nEvents - lastNEvents
			log.Printf("%s/s, %d events/s", ByteCountIEC(int64(bytesPerSec)), eventsPerSec)
			lastNRead = nRead
			lastNEvents = nEvents
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Scanner error: %v", err)
	}
	log.Printf("total read: %s, total events: %d", ByteCountIEC(int64(nRead)), nEvents)
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
