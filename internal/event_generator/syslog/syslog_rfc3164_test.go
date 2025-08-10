package syslog

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSyslogGeneratorRFC3164(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		app      string
		proc     string
		msg      string
		fakeTime time.Time
		eol      string
	}{
		{
			name:     "one-digit day",
			host:     "host1",
			app:      "app1",
			proc:     "1",
			msg:      "test message",
			fakeTime: time.Date(2023, 7, 4, 2, 3, 4, 0, time.UTC),
			eol:      "\n",
		},
		{
			name:     "two-digit day",
			host:     "host2",
			app:      "app2",
			proc:     "22",
			msg:      "test message2",
			fakeTime: time.Date(2023, 7, 14, 12, 34, 56, 0, time.UTC),
			eol:      "\r\n",
		},
		{
			name:     "weird app/proc",
			host:     "h",
			app:      "a-b_c.1",
			proc:     "007",
			msg:      "msg3",
			fakeTime: time.Date(2023, 12, 9, 23, 59, 59, 0, time.UTC),
			eol:      "END",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gen := NewSyslogGenerator(RFC3164,
				WithHostname(tc.host),
				WithAppName(tc.app),
				WithProcID(tc.proc),
				WithEndOfLine([]byte(tc.eol)),
			).(*SyslogGeneratorRFC3164)
			buf := make([]byte, 256)
			n, err := gen.GenerateMessage(buf, tc.fakeTime, []byte(tc.msg))
			require.NoError(t, err)
			out := buf[:n]
			// Build expected prefix
			prefix := "<33>" + tc.fakeTime.Format("Jan _2 15:04:05") + " " + tc.host + " " + tc.app + "[" + tc.proc + "]: "
			expected := prefix + tc.msg + tc.eol
			require.Equal(t, expected, string(out))
		})
	}
}
