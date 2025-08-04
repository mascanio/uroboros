package syslog

import (
	"fmt"
	"io"
	"math/rand"
	"time"
)

// Format specifies the syslog message format
// RFC3164 (BSD) or RFC5424 (IETF)
type Format int

const (
	RFC3164 Format = iota
	RFC5424
)

// SyslogGenerator generates valid syslog messages and implements io.Reader
// It generates a new random message on each Read, up to a specified count.
type SyslogGenerator struct {
	random            *rand.Rand // random source for message fields
	format            Format     // message format
	messagesTotal     int        // total number of messages to generate
	messagesGenerated int        // number of messages generated so far
	message           []byte
}

// NewSyslogGenerator returns a new SyslogGenerator with the given format and message count
func NewSyslogGenerator(format Format, numMessages int) *SyslogGenerator {
	return &SyslogGenerator{
		random:        rand.New(rand.NewSource(time.Now().UnixNano())),
		format:        format,
		messagesTotal: numMessages,
	}
}

// generateSyslogMessage creates a single syslog message in the selected format
func (g *SyslogGenerator) generateSyslogMessage(p []byte) (int, error) {
	if len(g.message) == 0 {
		switch g.format {
		case RFC5424:
			g.message = g.createSyslogMessageRFC5424()
		case RFC3164:
			g.message = g.createSyslogMessageRFC3164()
		}
	}
	if len(g.message) > len(p) {
		return 0, fmt.Errorf("buffer too small, %v %v", len(p), len(g.message))
	}
	copy(p, g.message)
	return len(g.message), nil
}

// generateSyslogMessageRFC3164 creates a single RFC 3164 syslog message
func (g *SyslogGenerator) createSyslogMessageRFC3164() []byte {
	pri := 33
	timestamp := time.Now().Format("Jan 2 15:04:05")
	hostname := "host"
	appName := "app"
	pid := 1000
	msg := fmt.Sprintf("This is a test syslog message %d", g.random.Intn(100000))
	s := fmt.Sprintf("<%d>%s %s %s[%d]: %s\n", pri, timestamp, hostname, appName, pid, msg)
	return []byte(s)
}

// generateSyslogMessageRFC5424 creates a single RFC 5424 syslog message
func (g *SyslogGenerator) createSyslogMessageRFC5424() []byte {
	pri := 33
	timestamp := time.Now().UTC().Format(time.RFC3339)
	version := 1
	hostname := "host.example.com"
	appName := "app"
	procID := "1"
	msgID := "ID1000"
	structuredData := "-" // could be extended
	msg := fmt.Sprintf("This is a test syslog message %d", g.random.Intn(100000))
	s := fmt.Sprintf("<%d>%d %s %s %s %s %s %s %s\n", pri, version, timestamp, hostname, appName, procID, msgID, structuredData, msg)
	return []byte(s)
}

// Read implements io.Reader, filling p with generated syslog messages
// Returns io.EOF after the specified number of messages have been generated.
func (g *SyslogGenerator) Read(p []byte) (int, error) {
	if g.messagesGenerated >= g.messagesTotal {
		return 0, io.EOF
	}
	var totalBytes int
	var nBytesRead int
	var err error
	for g.messagesGenerated < g.messagesTotal && err == nil {
		nBytesRead, err = g.generateSyslogMessage(p)
		if nBytesRead != 0 {
			g.messagesGenerated++
		}
		totalBytes += nBytesRead
		p = p[nBytesRead:]
	}
	return totalBytes, err
}
