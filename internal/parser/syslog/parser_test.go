package syslog

import (
	"reflect"
	"testing"
)

func equalSyslogMessage(a, b SyslogMessage) bool {
	// nil and empty slices are considered equal for all fields
	eq := func(x, y []byte) bool {
		return len(x) == 0 && len(y) == 0 || reflect.DeepEqual(x, y)
	}
	return a.Priority == b.Priority &&
		eq(a.Timestamp, b.Timestamp) &&
		eq(a.Hostname, b.Hostname) &&
		eq(a.AppName, b.AppName) &&
		eq(a.ProcID, b.ProcID) &&
		eq(a.MsgID, b.MsgID) &&
		eq(a.Structured, b.Structured) &&
		eq(a.Msg, b.Msg)
}

func TestParseSyslogMessageRFC3164(t *testing.T) {
	tests := []struct {
		name   string
		input  []byte
		expect SyslogMessage
	}{
		{
			name:  "Valid RFC3164",
			input: []byte("<34>Oct 11 22:14:15 mymachine su: 'su root' failed for lonvick on /dev/pts/8"),
			expect: SyslogMessage{
				Priority:  34,
				Timestamp: []byte("Oct"),
				Hostname:  []byte("11"),
				AppName:   []byte("22:14:15"),
				Msg:       []byte("mymachine su: 'su root' failed for lonvick on /dev/pts/8"),
			},
		},
		{
			name:  "No Priority",
			input: []byte("Oct 11 22:14:15 mymachine su: test"),
			expect: SyslogMessage{
				Priority:  0,
				Timestamp: []byte("Oct"),
				Hostname:  []byte("11"),
				AppName:   []byte("22:14:15"),
				Msg:       []byte("mymachine su: test"),
			},
		},
		{
			name:   "Empty Input",
			input:  []byte(""),
			expect: SyslogMessage{},
		},
		{
			name:  "Malformed Priority",
			input: []byte("<xx>Oct 11 22:14:15 host app msg"),
			expect: SyslogMessage{
				Priority:  0,
				Timestamp: []byte("Oct"),
				Hostname:  []byte("11"),
				AppName:   []byte("22:14:15"),
				Msg:       []byte("host app msg"),
			},
		},
		{
			name:  "Short Input",
			input: []byte("<34>foo"),
			expect: SyslogMessage{
				Priority: 34,
				Msg:      []byte("foo"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseSyslogMessageRFC3164(tt.input)
			if !equalSyslogMessage(got, tt.expect) {
				t.Errorf("got %+v, want %+v", got, tt.expect)
			}
		})
	}
}

func TestParseSyslogMessageRFC5424(t *testing.T) {
	tests := []struct {
		name   string
		input  []byte
		expect SyslogMessage
	}{
		{
			name:  "Valid RFC5424",
			input: []byte("<165>1 2003-10-11T22:14:15Z mymachine.example.com evntslog 1234 ID47 [exampleSDID@32473 iut=\"3\" eventSource=\"Application\" eventID=\"1011\"] BOMAn application event log entry..."),
			expect: SyslogMessage{
				Priority:   165,
				Timestamp:  []byte("2003-10-11T22:14:15Z"),
				Hostname:   []byte("mymachine.example.com"),
				AppName:    []byte("evntslog"),
				ProcID:     []byte("1234"),
				MsgID:      []byte("ID47"),
				Structured: []byte("[exampleSDID@32473"),
				Msg:        []byte("iut=\"3\" eventSource=\"Application\" eventID=\"1011\"] BOMAn application event log entry..."),
			},
		},
		{
			name:  "No Priority",
			input: []byte("1 2003-10-11T22:14:15Z host app 1234 ID47 - msg"),
			expect: SyslogMessage{
				Priority:   0,
				Timestamp:  []byte("2003-10-11T22:14:15Z"),
				Hostname:   []byte("host"),
				AppName:    []byte("app"),
				ProcID:     []byte("1234"),
				MsgID:      []byte("ID47"),
				Structured: []byte("-"),
				Msg:        []byte("msg"),
			},
		},
		{
			name:   "Empty Input",
			input:  []byte(""),
			expect: SyslogMessage{},
		},
		{
			name:  "Malformed Priority",
			input: []byte("<xx>1 2003-10-11T22:14:15Z host app 1234 ID47 - msg"),
			expect: SyslogMessage{
				Priority:   0,
				Timestamp:  []byte("2003-10-11T22:14:15Z"),
				Hostname:   []byte("host"),
				AppName:    []byte("app"),
				ProcID:     []byte("1234"),
				MsgID:      []byte("ID47"),
				Structured: []byte("-"),
				Msg:        []byte("msg"),
			},
		},
		{
			name:  "Short Input",
			input: []byte("<165>foo"),
			expect: SyslogMessage{
				Priority: 165,
				Msg:      []byte("foo"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseSyslogMessageRFC5424(tt.input)
			if !equalSyslogMessage(got, tt.expect) {
				t.Errorf("got %+v, want %+v", got, tt.expect)
			}
		})
	}
}

func BenchmarkParseSyslogMessageRFC3164(b *testing.B) {
	msg := []byte("<34>Oct 11 22:14:15 mymachine su: 'su root' failed for lonvick on /dev/pts/8")
	for b.Loop() {
		_ = ParseSyslogMessageRFC3164(msg)
	}
}

func BenchmarkParseSyslogMessageRFC3164_Short(b *testing.B) {
	msg := []byte("<13>foo")
	for b.Loop() {
		_ = ParseSyslogMessageRFC3164(msg)
	}
}

func BenchmarkParseSyslogMessageRFC3164_Long(b *testing.B) {
	msg := []byte("<34>Oct 11 22:14:15 mymachine su: " + string(make([]byte, 1000)))
	for b.Loop() {
		_ = ParseSyslogMessageRFC3164(msg)
	}
}

func BenchmarkParseSyslogMessageRFC3164_Malformed(b *testing.B) {
	msg := []byte("<xx>Oct 11 22:14:15 host app msg")
	for b.Loop() {
		_ = ParseSyslogMessageRFC3164(msg)
	}
}

func BenchmarkParseSyslogMessageRFC5424(b *testing.B) {
	msg := []byte("<165>1 2003-10-11T22:14:15Z mymachine.example.com evntslog 1234 ID47 [exampleSDID@32473 iut=\"3\" eventSource=\"Application\" eventID=\"1011\"] BOMAn application event log entry...")
	for b.Loop() {
		_ = ParseSyslogMessageRFC5424(msg)
	}
}

func BenchmarkParseSyslogMessageRFC5424_Short(b *testing.B) {
	msg := []byte("<165>foo")
	for b.Loop() {
		_ = ParseSyslogMessageRFC5424(msg)
	}
}

func BenchmarkParseSyslogMessageRFC5424_Long(b *testing.B) {
	msg := []byte("<165>1 2003-10-11T22:14:15Z host app 1234 ID47 - " + string(make([]byte, 2000)))
	for b.Loop() {
		_ = ParseSyslogMessageRFC5424(msg)
	}
}

func BenchmarkParseSyslogMessageRFC5424_Malformed(b *testing.B) {
	msg := []byte("<xx>1 2003-10-11T22:14:15Z host app 1234 ID47 - msg")
	for b.Loop() {
		_ = ParseSyslogMessageRFC5424(msg)
	}
}
