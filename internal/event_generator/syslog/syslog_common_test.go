package syslog

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewSyslogGenerator_BuilderTypes(t *testing.T) {
	g3164 := NewSyslogGenerator(RFC3164)
	_, ok := g3164.(*SyslogGeneratorRFC3164)
	require.True(t, ok)

	g5424 := NewSyslogGenerator(RFC5424)
	_, ok = g5424.(*SyslogGeneratorRFC5424)
	require.True(t, ok)
}

func TestOptions_All(t *testing.T) {
	g := NewSyslogGenerator(RFC5424,
		WithHostname("h"),
		WithAppName("a"),
		WithProcID("p"),
		WithMsgID("m"),
		WithStructuredData("[sd@1 x=\"y\"]"),
		WithEndOfLine([]byte("EOL")),
	).(*SyslogGeneratorRFC5424)
	require.Equal(t, "h", g.hostname)
	require.Equal(t, "a", g.appName)
	require.Equal(t, "p", g.procID)
	require.Equal(t, "m", g.msgID)
	require.Equal(t, "[sd@1 x=\"y\"]", g.structuredData)
	require.Equal(t, []byte("EOL"), g.endOfLine)
}
