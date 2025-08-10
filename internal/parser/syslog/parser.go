package syslog

import "strconv"

// SyslogMessage represents a parsed syslog message (minimal fields)
type SyslogMessage struct {
	Priority   int
	Timestamp  []byte
	Hostname   []byte
	AppName    []byte
	ProcID     []byte
	MsgID      []byte // RFC5424
	Structured []byte // RFC5424
	Msg        []byte
}

// ParseSyslogMessageRFC3164 parses a syslog message in RFC3164 format.
func ParseSyslogMessageRFC3164(data []byte) SyslogMessage {
	msg := SyslogMessage{}
	// Try to parse priority
	if len(data) > 0 && data[0] == '<' {
		end := 1
		for ; end < len(data) && data[end] != '>'; end++ {
		}
		if end < len(data) && data[end] == '>' {
			if pri, err := strconv.Atoi(string(data[1:end])); err == nil {
				msg.Priority = pri
			}
			data = data[end+1:]
		}
	}
	start := 0
	next := 0
	// Extract timestamp
	for start < len(data) && data[start] == ' ' {
		start++
	}
	next = start
	for next < len(data) && data[next] != ' ' {
		next++
	}
	if next == start || next >= len(data) {
		msg.Msg = data
		return msg
	}
	msg.Timestamp = data[start:next]
	start = next
	// Extract hostname
	for start < len(data) && data[start] == ' ' {
		start++
	}
	next = start
	for next < len(data) && data[next] != ' ' {
		next++
	}
	if next == start || next >= len(data) {
		msg.Msg = data
		return msg
	}
	msg.Hostname = data[start:next]
	start = next
	// Extract appName
	for start < len(data) && data[start] == ' ' {
		start++
	}
	next = start
	for next < len(data) && data[next] != ' ' {
		next++
	}
	if next == start || next >= len(data) {
		msg.Msg = data
		return msg
	}
	msg.AppName = data[start:next]
	start = next
	// The rest is the message
	for start < len(data) && data[start] == ' ' {
		start++
	}
	msg.Msg = data[start:]
	return msg
}

// ParseSyslogMessageRFC5424 parses a syslog message in RFC5424 format.
func ParseSyslogMessageRFC5424(data []byte) SyslogMessage {
	msg := SyslogMessage{}
	// Try to parse priority
	if len(data) > 0 && data[0] == '<' {
		end := 1
		for ; end < len(data) && data[end] != '>'; end++ {
		}
		if end < len(data) && data[end] == '>' {
			if pri, err := strconv.Atoi(string(data[1:end])); err == nil {
				msg.Priority = pri
			}
			data = data[end+1:]
		}
	}
	start := 0
	next := 0
	// Extract version (skip)
	for start < len(data) && data[start] == ' ' {
		start++
	}
	next = start
	for next < len(data) && data[next] != ' ' {
		next++
	}
	if next == start || next >= len(data) {
		msg.Msg = data
		return msg
	}
	start = next
	// Extract timestamp
	for start < len(data) && data[start] == ' ' {
		start++
	}
	next = start
	for next < len(data) && data[next] != ' ' {
		next++
	}
	if next == start || next >= len(data) {
		msg.Msg = data
		return msg
	}
	msg.Timestamp = data[start:next]
	start = next
	// Extract hostname
	for start < len(data) && data[start] == ' ' {
		start++
	}
	next = start
	for next < len(data) && data[next] != ' ' {
		next++
	}
	if next == start || next >= len(data) {
		msg.Msg = data
		return msg
	}
	msg.Hostname = data[start:next]
	start = next
	// Extract appName
	for start < len(data) && data[start] == ' ' {
		start++
	}
	next = start
	for next < len(data) && data[next] != ' ' {
		next++
	}
	if next == start || next >= len(data) {
		msg.Msg = data
		return msg
	}
	msg.AppName = data[start:next]
	start = next
	// Extract procID
	for start < len(data) && data[start] == ' ' {
		start++
	}
	next = start
	for next < len(data) && data[next] != ' ' {
		next++
	}
	if next == start || next >= len(data) {
		msg.Msg = data
		return msg
	}
	msg.ProcID = data[start:next]
	start = next
	// Extract msgID
	for start < len(data) && data[start] == ' ' {
		start++
	}
	next = start
	for next < len(data) && data[next] != ' ' {
		next++
	}
	if next == start || next >= len(data) {
		msg.Msg = data
		return msg
	}
	msg.MsgID = data[start:next]
	start = next
	// Extract structured data
	for start < len(data) && data[start] == ' ' {
		start++
	}
	next = start
	for next < len(data) && data[next] != ' ' {
		next++
	}
	if next == start || next > len(data) {
		msg.Msg = data
		return msg
	}
	msg.Structured = data[start:next]
	start = next
	for start < len(data) && data[start] == ' ' {
		start++
	}
	msg.Msg = data[start:]
	return msg
}
