package syslog

import (
	"time"
)

// Format specifies the syslog message format
// RFC3164 (BSD) or RFC5424 (IETF)
type Format int

const (
	RFC3164 Format = iota
	RFC5424
)

type SyslogMessageGenerator interface {
	GenerateMessage(p []byte, now time.Time, msg []byte) (int, error)
}

type Option interface {
	apply3164(*SyslogGeneratorRFC3164)
	apply5424(*SyslogGeneratorRFC5424)
}

type endOfLineOption []byte

func (e endOfLineOption) apply3164(g *SyslogGeneratorRFC3164) { g.endOfLine = []byte(e) }
func (e endOfLineOption) apply5424(g *SyslogGeneratorRFC5424) { g.endOfLine = []byte(e) }
func WithEndOfLine(e []byte) Option                           { return endOfLineOption(e) }

type hostnameOption string

func (h hostnameOption) apply3164(g *SyslogGeneratorRFC3164) { g.hostname = string(h) }
func (h hostnameOption) apply5424(g *SyslogGeneratorRFC5424) { g.hostname = string(h) }
func WithHostname(h string) Option                           { return hostnameOption(h) }

type appNameOption string

func (a appNameOption) apply3164(g *SyslogGeneratorRFC3164) { g.appName = string(a) }
func (a appNameOption) apply5424(g *SyslogGeneratorRFC5424) { g.appName = string(a) }
func WithAppName(a string) Option                           { return appNameOption(a) }

type procIDOption string

func (p procIDOption) apply3164(g *SyslogGeneratorRFC3164) { g.procID = string(p) }
func (p procIDOption) apply5424(g *SyslogGeneratorRFC5424) { g.procID = string(p) }
func WithProcID(p string) Option                           { return procIDOption(p) }

type msgIDOption string

func (m msgIDOption) apply3164(g *SyslogGeneratorRFC3164) {}
func (m msgIDOption) apply5424(g *SyslogGeneratorRFC5424) { g.msgID = string(m) }
func WithMsgID(m string) Option                           { return msgIDOption(m) }

type structuredDataOption string

func (s structuredDataOption) apply3164(g *SyslogGeneratorRFC3164) {}
func (s structuredDataOption) apply5424(g *SyslogGeneratorRFC5424) { g.structuredData = string(s) }
func WithStructuredData(s string) Option                           { return structuredDataOption(s) }

func NewSyslogGenerator(rfc Format, opts ...Option) SyslogMessageGenerator {
	defaultEOL := []byte("\n")
	switch rfc {
	case RFC3164:
		g := &SyslogGeneratorRFC3164{
			hostname:  "host.example.com",
			appName:   "app",
			procID:    "1",
			endOfLine: defaultEOL,
		}
		for _, opt := range opts {
			opt.apply3164(g)
		}
		g.rfc3164Prefix, g.rfc3164TimestampOffset = buildRFC3164Prefix(g)
		return g
	case RFC5424:
		g := &SyslogGeneratorRFC5424{
			hostname:       "host.example.com",
			appName:        "app",
			procID:         "1",
			msgID:          "ID1000",
			structuredData: "-",
			endOfLine:      defaultEOL,
		}
		for _, opt := range opts {
			opt.apply5424(g)
		}
		g.rfc5424Prefix, g.rfc5424TimestampOffset = buildRFC5424Prefix(g)
		return g
	default:
		return nil
	}
}
