package syslog

import (
	"strconv"
	"time"
)

type SyslogGeneratorRFC3164 struct {
	hostname                      string
	appName                       string
	procID                        string
	rfc3164Prefix                 []byte
	rfc3164TimestampOffset        int
	lastRFC3164TimestampUpdate    time.Time
	lastRFC3164FormattedTimestamp string
	endOfLine                     []byte
}

func (g *SyslogGeneratorRFC3164) GenerateMessage(p []byte, now time.Time, msg []byte) (int, error) {
	if g.lastRFC3164TimestampUpdate.IsZero() {
		g.lastRFC3164TimestampUpdate = now
		g.lastRFC3164FormattedTimestamp = now.Format("Jan _2 15:04:05")
		copy(g.rfc3164Prefix[g.rfc3164TimestampOffset:], g.lastRFC3164FormattedTimestamp)
	} else if now.Sub(g.lastRFC3164TimestampUpdate) >= 500*time.Millisecond {
		g.lastRFC3164TimestampUpdate = now
		g.lastRFC3164FormattedTimestamp = now.Format("Jan _2 15:04:05")
		copy(g.rfc3164Prefix[g.rfc3164TimestampOffset:], g.lastRFC3164FormattedTimestamp)
	}
	// Copy prefix (with cached timestamp)
	copy(p, g.rfc3164Prefix)
	idx := len(g.rfc3164Prefix)
	// MSG
	copy(p[idx:], msg)
	idx += len(msg)
	copy(p[idx:], g.endOfLine)
	idx += len(g.endOfLine)
	return idx, nil
}

// buildRFC3164Prefix returns the static prefix for RFC3164 (up to just before the timestamp)
func buildRFC3164Prefix(g *SyslogGeneratorRFC3164) ([]byte, int) {
	// Format: <PRI>TIMESTAMP HOSTNAME APP-NAME[PROCID]:
	pri := 33
	var buf [128]byte
	idx := 0
	buf[idx] = '<'
	idx++
	n := strconv.AppendInt(buf[idx:idx], int64(pri), 10)
	idx += len(n)
	buf[idx] = '>'
	idx++
	timestampOffset := idx
	// Reserve 15 spaces for timestamp
	for range 15 {
		buf[idx] = ' '
		idx++
	}
	buf[idx] = ' '
	idx++
	// Write hostname
	hostname := g.hostname
	if len(hostname) == 0 {
		return append([]byte{}, buf[:idx]...), timestampOffset
	}
	idx += copy(buf[idx:], hostname)
	if idx < len(buf) {
		buf[idx] = ' '
		idx++
	}
	// Write appName
	appName := g.appName
	if len(appName) == 0 {
		return append([]byte{}, buf[:idx]...), timestampOffset
	}
	idx += copy(buf[idx:], appName)
	if idx < len(buf) {
		buf[idx] = '['
		idx++
	}
	// Write procID
	procID := g.procID
	if procID == "" {
		procID = "1000"
	}
	procIDStr := procID
	if len(procIDStr) == 0 {
		return append([]byte{}, buf[:idx]...), timestampOffset
	}
	idx += copy(buf[idx:], procIDStr)
	if idx < len(buf) {
		buf[idx] = ']'
		idx++
	}
	if idx < len(buf) {
		buf[idx] = ':'
		idx++
	}
	if idx < len(buf) {
		buf[idx] = ' '
		idx++
	}
	return append([]byte{}, buf[:idx]...), timestampOffset
}
