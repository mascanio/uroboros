package syslog

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSyslogGeneratorRFC5424(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		app      string
		proc     string
		msgid    string
		sd       string
		msg      string
		fakeTime time.Time
		eol      string
	}{
		{
			name:     "basic",
			host:     "host1",
			app:      "app1",
			proc:     "1",
			msgid:    "MSGID",
			sd:       "-",
			msg:      "test message",
			fakeTime: time.Date(2023, 7, 4, 12, 34, 56, 0, time.UTC),
			eol:      "\n",
		},
		{
			name:     "weird fields",
			host:     "h",
			app:      "a-b_c.1",
			proc:     "007",
			msgid:    "X-Y",
			sd:       "[id@32456 foo=\"bar\"]",
			msg:      "msg2",
			fakeTime: time.Date(2023, 12, 9, 1, 2, 3, 0, time.UTC),
			eol:      "END",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gen := NewSyslogGenerator(RFC5424,
				WithHostname(tc.host),
				WithAppName(tc.app),
				WithProcID(tc.proc),
				WithMsgID(tc.msgid),
				WithStructuredData(tc.sd),
				WithEndOfLine([]byte(tc.eol)),
			).(*SyslogGeneratorRFC5424)
			buf := make([]byte, 512)
			n, err := gen.GenerateEvent(buf, tc.fakeTime, []byte(tc.msg))
			require.NoError(t, err)
			out := buf[:n]
			prefix := "<33>1 " + tc.fakeTime.UTC().Format("2006-01-02T15:04:05Z") + " " + tc.host + " " + tc.app + " " + tc.proc + " " + tc.msgid + " " + tc.sd + " "
			expected := prefix + tc.msg + tc.eol
			require.Equal(t, expected, string(out))
		})
	}
}
