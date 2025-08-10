package syslog

import (
	"strconv"
	"time"
)

type SyslogGeneratorRFC5424 struct {
	hostname                      string
	appName                       string
	procID                        string
	msgID                         string
	structuredData                string
	rfc5424Prefix                 []byte
	rfc5424TimestampOffset        int
	lastRFC5424TimestampUpdate    time.Time
	lastRFC5424FormattedTimestamp string
	endOfLine                     []byte
}

func (g *SyslogGeneratorRFC5424) GenerateEvent(p []byte, now time.Time, msg []byte) (int, error) {
	if g.lastRFC5424TimestampUpdate.IsZero() {
		g.lastRFC5424TimestampUpdate = now
		g.lastRFC5424FormattedTimestamp = now.UTC().Format(time.RFC3339)
		copy(g.rfc5424Prefix[g.rfc5424TimestampOffset:], g.lastRFC5424FormattedTimestamp)
	} else if now.Sub(g.lastRFC5424TimestampUpdate) >= 500*time.Millisecond {
		g.lastRFC5424TimestampUpdate = now
		g.lastRFC5424FormattedTimestamp = now.UTC().Format(time.RFC3339)
		copy(g.rfc5424Prefix[g.rfc5424TimestampOffset:], g.lastRFC5424FormattedTimestamp)
	}
	copy(p, g.rfc5424Prefix)
	idx := len(g.rfc5424Prefix)
	copy(p[idx:], msg)
	idx += len(msg)
	copy(p[idx:], g.endOfLine)
	idx += len(g.endOfLine)
	return idx, nil
}

// buildRFC5424Prefix returns the static prefix for RFC5424 (up to just before the timestamp)
func buildRFC5424Prefix(g *SyslogGeneratorRFC5424) ([]byte, int) {
	// Format: <PRI>VERSION TIMESTAMP HOSTNAME APP-NAME PROCID MSGID STRUCTUREDDATA
	pri := 33
	version := 1
	var buf [256]byte
	idx := 0
	buf[idx] = '<'
	idx++
	n := strconv.AppendInt(buf[idx:idx], int64(pri), 10)
	idx += len(n)
	buf[idx] = '>'
	idx++
	n = strconv.AppendInt(buf[idx:idx], int64(version), 10)
	idx += len(n)
	buf[idx] = ' '
	idx++
	timestampOffset := idx
	// Reserve 20 spaces for timestamp
	for range 20 {
		buf[idx] = ' '
		idx++
	}
	buf[idx] = ' '
	idx++
	copy(buf[idx:], g.hostname)
	idx += len(g.hostname)
	buf[idx] = ' '
	idx++
	copy(buf[idx:], g.appName)
	idx += len(g.appName)
	buf[idx] = ' '
	idx++
	copy(buf[idx:], g.procID)
	idx += len(g.procID)
	buf[idx] = ' '
	idx++
	copy(buf[idx:], g.msgID)
	idx += len(g.msgID)
	buf[idx] = ' '
	idx++
	copy(buf[idx:], g.structuredData)
	idx += len(g.structuredData)
	buf[idx] = ' '
	idx++
	return append([]byte{}, buf[:idx]...), timestampOffset
}
