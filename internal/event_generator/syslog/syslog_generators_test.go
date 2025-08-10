package syslog

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestSyslogGenerator_Table(t *testing.T) {
	tests := []struct {
		name     string
		format   Format
		host     string
		app      string
		proc     string
		msgid    string
		sd       string
		msg      string
		fakeTime time.Time
		eol      string
		expected string
	}{
		{
			name:     "RFC3164 basic",
			format:   RFC3164,
			host:     "testhost",
			app:      "testapp",
			proc:     "123",
			msg:      "test message",
			fakeTime: time.Date(2023, 7, 4, 12, 34, 56, 0, time.UTC),
			eol:      "\r\n",
			expected: "<33>Jul  4 12:34:56 testhost testapp[123]: test message\r\n",
		},
		{
			name:     "RFC5424 basic",
			format:   RFC5424,
			host:     "testhost",
			app:      "testapp",
			proc:     "123",
			msgid:    "MSGID",
			sd:       "-",
			msg:      "test message",
			fakeTime: time.Date(2023, 7, 4, 12, 34, 56, 0, time.UTC),
			eol:      "\n",
			expected: "<33>1 2023-07-04T12:34:56Z testhost testapp 123 MSGID - test message\n",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var gen SyslogMessageGenerator
			if tc.format == RFC3164 {
				gen = NewSyslogGenerator(RFC3164,
					WithHostname(tc.host),
					WithAppName(tc.app),
					WithProcID(tc.proc),
					WithEndOfLine([]byte(tc.eol)),
				)
			} else {
				gen = NewSyslogGenerator(RFC5424,
					WithHostname(tc.host),
					WithAppName(tc.app),
					WithProcID(tc.proc),
					WithMsgID(tc.msgid),
					WithStructuredData(tc.sd),
					WithEndOfLine([]byte(tc.eol)),
				)
			}
			buf := make([]byte, 512)
			n, err := gen.GenerateMessage(buf, tc.fakeTime, []byte(tc.msg))
			require.NoError(t, err)
			out := buf[:n]
			require.Equal(t, tc.expected, string(out))
		})
	}
}
