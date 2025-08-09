package syslog

import (
	"bytes"
	"strconv"
)

// SyslogSplitConfig allows configuring the split function for syslog.
type SyslogSplitConfig struct {
	UseOctetCounting bool // for RFC5424
	Delimiter        byte // for newline, e.g. '\n'
}

// NewSyslogSplitFuncOctetCounting returns a bufio.SplitFunc for octet-counted syslog (RFC5424).
func NewSyslogSplitFuncOctetCounting() func(data []byte, atEOF bool) (advance int, token []byte, err error) {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if len(data) == 0 {
			return 0, nil, nil
		}
		spaceIdx := bytes.IndexByte(data, ' ')
		if spaceIdx > 0 && isAllDigits(data[:spaceIdx]) {
			msgLen, convErr := strconv.Atoi(string(data[:spaceIdx]))
			if convErr == nil && len(data) >= spaceIdx+1+msgLen {
				return spaceIdx + 1 + msgLen, data[spaceIdx+1 : spaceIdx+1+msgLen], nil
			}
			return 0, nil, nil
		}
		return 0, nil, nil
	}
}

// NewSyslogSplitFuncDelimiter returns a bufio.SplitFunc for delimiter-based syslog (e.g., RFC3164, newline-delimited).
func NewSyslogSplitFuncDelimiter(delim []byte) func(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if len(delim) == 0 {
		delim = []byte{'\n'}
	}
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if len(data) == 0 {
			return 0, nil, nil
		}
		if i := bytes.Index(data, delim); i >= 0 {
			return i + len(delim), data[:i], nil
		}
		if atEOF {
			return len(data), data, nil
		}
		return 0, nil, nil
	}
}

func isAllDigits(b []byte) bool {
	for _, c := range b {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(b) > 0
}
