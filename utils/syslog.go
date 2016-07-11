package utils

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

// Severity defines the syslog severity
type Severity int
const (
	Emergency Severity = iota
	Alert
	Critical
	Error
	Warning
	Notice
	Informational
	Debug
	UnknownSeverity
)

// String returns a human-readable representation of the Severity value.
func (s Severity) String() string {
	switch s {
	case Emergency:
		return "EMERGENCY"
	case Alert:
		return "ALERT"
	case Critical:
		return "CRITICAL"
	case Error:
		return "ERROR"
	case Warning:
		return "WARNING"
	case Notice:
		return "NOTICE"
	case Informational:
		return "INFO"
	case Debug:
		return "DEBUG"
	}
	panic(errors.New("syslog: Invalid Severity values"))
}

// SeverityFromString converts log level given as string into proper Severity value.
// If invalid string is given, function returns 'UnknownSeverity' value.
func SeverityFromString(lvl string) Severity {
	loglvl := UnknownSeverity
	switch strings.ToUpper(lvl) {
	case "EMERGENCY":
		loglvl = Emergency
	case "ALERT":
		loglvl = Alert
	case "CRITICAL":
		loglvl = Critical
	case "ERROR":
		loglvl = Error
	case "WARNING":
		loglvl = Warning
	case "NOTICE":
		loglvl = Notice
	case "INFO":
		loglvl = Informational
	case "DEBUG":
		loglvl = Debug
	}
	return loglvl
}

// Facility defines the syslog facility
type Facility int
const (
	FacKernel Facility = iota
	FacUser
	FacMail
	FacSystem
	FacSecurity4
	FacSyslogd
	FacLine
	FacNetwork
	FacUUCP
	FacClock9
	FacSecurity10
	FacFTP
	FacNTP
	FacLogAudit
	FacLogAlert
	FacClock15
	FacLocal0
	FacLocal1
	FacLocal2
	FacLocal3
	FacLocal4
	FacLocal5
	FacLocal6
	FacLocal7
)

const (
	// TimestampFmt defines a standard syslog message timestamp format
	TimestampFmt = "Jan _2 15:04:05"
	// SyslogPort defines the standard UDP port for syslog (514)
	SyslogPort = 514
)

// SyslogMsg defines a syslog message type.
type SyslogMsg struct {
	Sev                 Severity
	Fac                 Facility
	timestamp, Hostname string
	Msg                 string
}

// Priority returns a value of syslog priority.
func (s *SyslogMsg) Priority() string {
	pri := int(s.Sev) + (8 * int(s.Fac))
	return fmt.Sprintf("<%d>", pri)
}

// TimeStamp returns a syslog message timestamp.
func (s *SyslogMsg) TimeStamp() string { return s.timestamp }

// SetTimestamp sets a new timestamp for the syslog message.
func (s *SyslogMsg) SetTimestamp(stamp time.Time) { s.timestamp = stamp.Format(TimestampFmt) }

// SSetTimestamp sets a new timestamp for the syslog message (timestamp is given as a string value).
func (s *SyslogMsg) SSetTimestamp(stamp string) error {

	t, err := time.Parse(TimestampFmt, stamp)
	if err != nil {
		return err
	}
	s.SetTimestamp(t)
	return nil
}

// Get returns the properly formatted syslog message.
func (s *SyslogMsg) Get() string { return fmt.Sprintf("%s%s %s %s", s.Priority(), s.timestamp, s.Hostname, s.Msg) }

// Send sends the syslog message to given IP address.
func (s *SyslogMsg) Send(ip string) error {

	//var addr net.IP
	// local IP address overrides the Hostname field
	if ip != "" {
		s.Hostname = ip
	}
	addr := net.ParseIP(s.Hostname)

	// let's make an UDP connection and send the message
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{addr, SyslogPort, ""})
	if err != nil {
		return err
	}
	defer conn.Close()
	fmt.Fprintf(conn, s.Get())
	return nil
}

// NewSyslogMsg creates new syslog message with default fields.
func NewSyslogMsg() *SyslogMsg { return &SyslogMsg{Informational, FacLocal0, "", "", ""} }
