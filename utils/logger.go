package utils
/*
 * logger.go -  custom logger implementation
 *
 * Default log module is simply not enough
 *
 * History:
 *  1   Jul11   MR  The initial version
 *  2   May14   MR  Refactoring and simplification: the LogLevel type is out, Severity is now used instead. The second is 
 *                  introduction of concurency: log can now run as goroutine and messages are sent over a channel.
 */

import (
	"fmt"
	"os"
	"time"
)

/************************** Logger ***********************************/

// LogHandler is an interface defining methods for various log handlers
type LogHandler interface {
	Severity() Severity
	SetSeverity(Severity)
	Format() string
	SetFormat(fmt string)
	String() string
	Start() error
	Close()
	Send(Severity, string)
	Clear() error
}

/************************** logHandler ***********************************/

// a private struct that defines log handler data structures
type logHandler struct {
	// set severity for this handler
	sev Severity

	// a formatter for this handler
	format string

	// a handler's channel onto which log messages are sent
	msgch chan *logmsg

	// a channel to signal when to stop the handler goroutine
	stop chan int
}

// Severity returns the severity value.
func (l *logHandler) Severity() Severity { return l.sev }

// SetSeverity resets the severity value for the log handler.
func (l *logHandler) SetSeverity(s Severity) { l.sev = s }

// Format returns the log message format value.
func (l *logHandler) Format() string { return l.format }

// SetFormat resets the log message format.
func (l *logHandler) SetFormat(fmt string) { l.format = fmt }

// Create a new log handler instance.
func newLogHandler(fmt string, sev Severity) *logHandler { return &logHandler{sev, fmt, nil, nil} }

// Clear clears the log (empty implementation to satisfy the interface, only file logger needs this one...)
func (l *logHandler) Clear() error { return l.Clear() }

/************************** Log ***********************************/

// helper private struct that defines a log message: severity and message text
type logmsg struct {
	sev Severity
	msg string
}

// Log is a list (a slice) of different log handlers that can be added at will.
type Log struct {
	// Handlers is a list of log handlers
	Handlers []LogHandler
}

// String returns a human-readable representation of the Log instance.
func (l *Log) String() string {

	s := ""
	for _, h := range l.Handlers {
		if h != nil {
			s += fmt.Sprint(h.String())
		}
	}
	return s
}

// AddHandler appends a new handler to the list of handlers
func (l *Log) AddHandler(h LogHandler) []LogHandler { return append(l.Handlers, h) }

/*
// A dispatch log messages method.
// Calls all needed log handlers and logs the given message with given level.
// If an unknown log level is received, do nothing.
func (l *Log) dispatch(sev Severity, msg string) {
	for _, h := range l.Handlers {
		switch sev {
		case Emergency:
			h.Emergency(msg)
		case Alert:
			h.Alert(msg)
		case Critical:
			h.Critical(msg)
		case Error:
			h.Error(msg)
		case Warning:
			h.Warning(msg)
		case Notice:
			h.Notice(msg)
		case Informational:
			h.Info(msg)
		case Debug:
			h.Debug(msg)
		}
	}
}
*/

// Log is a generic log method: send a message with given severity.
func (l *Log) Log(sev Severity, msg string) {
	for _, h := range l.Handlers {
		h.Send(sev, msg)
	}
}

// LogS is a pure string version of the Log() method: send a message with given severity (here given as string).
func (l *Log) LogS(sev, msg string) {
	s := SeverityFromString(sev)
	for _, h := range l.Handlers {
		h.Send(s, msg)
	}
}

// Debug logs a debug message.
func (l *Log) Debug(msg string) {
	for _, h := range l.Handlers {
		h.Send(Debug, msg)
	}
}

// Info logs an informational message.
func (l *Log) Info(msg string) {
	for _, h := range l.Handlers {
		h.Send(Informational, msg)
	}
}

// Notice logs a notice message.
func (l *Log) Notice(msg string) {
	for _, h := range l.Handlers {
		h.Send(Notice, msg)
	}
}

// Warning logs a warning message.
func (l *Log) Warning(msg string) {
	for _, h := range l.Handlers {
		h.Send(Warning, msg)
	}
}

// Error logs an error message.
func (l *Log) Error(msg string) {
	for _, h := range l.Handlers {
		h.Send(Error, msg)
	}
}

// Critical logs a critical message.
func (l *Log) Critical(msg string) {
	for _, h := range l.Handlers {
		h.Send(Critical, msg)
	}
}

// Alert logs an alert message.
func (l *Log) Alert(msg string) {
	for _, h := range l.Handlers {
		h.Send(Alert, msg)
	}
}

// Emergency logs an emergency message.
func (l *Log) Emergency(msg string) {
	for _, h := range l.Handlers {
		h.Send(Emergency, msg)
	}
}

// Close closes the log.
func (l *Log) Close() {
	for _, h := range l.Handlers {
		h.Close()
	}
}

// Clear clears the contents of the log. (empty implementation to satisfy the interface, only FileHandler actually needs one...)
func (l *Log) Clear() {
	for _, h := range l.Handlers {
		h.Clear()
	}
}

// NewLog creates new logger, specify the number of log handlers and create needed channels: the one onto which the log messages 
// are sent and the other where signal when to stop is sent. Return the Log instance.
func NewLog() *Log {
	// create new Log instance
	return &Log{make([]LogHandler, 0, 2)}
}

// Start starts the log handlers.
func (l *Log) Start() error {
	var err error
	for _, h := range l.Handlers {
		if err = h.Start(); err != nil {
			return err
		}
	}
	return nil
}

/************************** Formatter  ***********************************/

// Formatter is an interface defining the generic formatter
type Formatter interface {
	Format(string)
}

/************************** FileHandler ***********************************/

// FileHandler is a handler that writes messages to local log file.
type FileHandler struct {
	// all handlers share common data structures
	*logHandler

	// file descriptor for the file log handler
	file *os.File

	// filename
	filename string
}

// Write a messages with given severity to a logfile.
func (f *FileHandler) write(sev Severity, msg string) {
	if f.Severity() >= sev {
		fmt.Fprintf(f.file, f.Format(), Now(), sev, msg)
	}
}

// Close closes the file handler.
func (f *FileHandler) Close() {

	// send a signal to quit goroutine
	if f.stop != nil {
		close(f.logHandler.msgch)
		f.stop <- 1
	}

	if f.file != nil {
		f.file.Close()
	}
}

// String returns a human-readable representation of the FileHandler instance.
func (f *FileHandler) String() string {
	return fmt.Sprintf("  FileHandler: fmt=%q, lvl=%-10s, fd=%d\n", f.Format(), f.Severity(), f.file.Fd())
}

// Send sends a log message onto an internal channel.
func (f *FileHandler) Send(sev Severity, msg string) {
	if f.logHandler.msgch != nil {
		f.logHandler.msgch <- &logmsg{sev, msg}
	}
}

// Clear clears the contents of the log file
func (f *FileHandler) Clear() error {

	var err error

	f.Close() // we must close the file

	if err = os.Remove(f.filename); err != nil {
		return err
	}
	if f.file, err = os.OpenFile(f.filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755); err != nil {
		return err
	}
	f.Start() // we must close the file
	return err
}

// Start runs handler as a goroutine.
func (f *FileHandler) Start() error {

	f.logHandler.msgch = make(chan *logmsg, 10) // message channel (buffered)
	f.logHandler.stop = make(chan int, 1)       // stop channel

	go func(f *FileHandler) {

		for {
			select {
			// when message is received over channel, write it
			case m, ok := <-f.logHandler.msgch:
				if ok {
					f.write(m.sev, m.msg)
				}
			// when data is received over stop channel, just exit the goroutine
			case <-f.logHandler.stop:
				return
	//		default: // do nothing
			}
		}
	}(f)

	return nil
}

// NewFileHandler creates a new file handler.
func NewFileHandler(filename string, fmt string, sev Severity) (*FileHandler, error) {
	// open log file for append data
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	return &FileHandler{newLogHandler(fmt, sev), f, filename}, err
}

/************************** StreamHandler ***********************************/

// StreamHandler is a handler that writes messages to STDOUT (console)
type StreamHandler FileHandler

// Write a message with given severity to STDOUT.
func (s *StreamHandler) write(sev Severity, msg string) {
	if s.Severity() >= sev {
		fmt.Printf(s.Format(), Now(), sev, msg)
	}
}

// String returns a human-readable representation of the StreamHandler instance.
func (s *StreamHandler) String() string {
	return fmt.Sprintf("StreamHandler: fmt=%q, lvl=%-10s\n", s.Format(), s.Severity())
}

// Close closes the stream handler.
func (s *StreamHandler) Close() {
	// send a signal to quit goroutine
	if s.stop != nil {
		close(s.logHandler.msgch)
		s.stop <- 1
	}
}

// Send sends a log message onto internal channel.
func (s *StreamHandler) Send(sev Severity, msg string) {
	if s.logHandler.msgch != nil {
		s.logHandler.msgch <- &logmsg{sev, msg}
	}
}

// Start runs handler as a goroutine.
func (s *StreamHandler) Start() error {

	s.logHandler.msgch = make(chan *logmsg, 10) // message channel (buffered)
	s.logHandler.stop = make(chan int, 1)       // stop channel

	go func(s *StreamHandler) {

		for {
			select {
			// when message is received over channel, write it
			case m, ok := <-s.logHandler.msgch:
				//fmt.Printf("DEBUG, logger: msg=%v\n", m) // DEBUG
				if ok {
					s.write(m.sev, m.msg)
				}
			// when data is received over stop channel, just exit the goroutine
			case <-s.logHandler.stop:
				return
			default: // do nothing
			}
		}
	}(s)

	return nil
}

// Clear clears the log (empty implementation to satisfy the interface, only file logger needs this one...).
func (s *StreamHandler) Clear() error { return nil }

// NewStreamHandler creates a new stream handler.
func NewStreamHandler(fmt string, sev Severity) *StreamHandler { return &StreamHandler{newLogHandler(fmt, sev), os.Stdout, ""} }

/************************** SyslogHandler ***********************************/

// SyslogHandler is a handler that sends the log messages to standard syslog port (UDP 514)
type SyslogHandler struct {
	// all handlers share common data structures
	*logHandler

	// IP address of the syslog server
	IP string

	// a syslog message built according to RFC
	*SyslogMsg
}

// Write a log message with given severity to wire.
func (s *SyslogHandler) write(level Severity, msg string) error {
	if s.Severity() >= level {
		s.Fac = FacLocal0
		s.Sev = level
		s.Msg = fmt.Sprintf("%s %s", level.String(), msg)
		t := time.Now()
		s.SetTimestamp(t)
		err := s.SyslogMsg.Send(s.IP)
		if err != nil {
			return err
		}
	}
	return nil
}

// String returns a human-readable representation of the SyslogHandler instance.
func (s *SyslogHandler) String() string {
	return fmt.Sprintf("SyslogHandler: fmt=%q, lvl=%-10s, Server=%q\n", s.Format(), s.Severity(), s.IP)
}

// Close closes the syslog handler.
func (s *SyslogHandler) Close() {
	// send a signal to quit goroutine
	if s.stop != nil {
		close(s.logHandler.msgch)
		s.logHandler.stop <- 1
	}
}

// Send sends a log message onto internal channel.
func (s *SyslogHandler) Send(sev Severity, msg string) {
	if s.logHandler.msgch != nil {
		s.logHandler.msgch <- &logmsg{sev, msg}
	}
}

// Start runs a handler as a goroutine.
func (s *SyslogHandler) Start() error {

	s.logHandler.msgch = make(chan *logmsg, 10) // message channel (buffered)
	s.logHandler.stop = make(chan int, 1)       // stop channel

	go func(s *SyslogHandler) {

		for {
			select {
			// when message is received over channel, write it
			case m, ok := <-s.logHandler.msgch:
				if ok {
					s.write(m.sev, m.msg)
				}
			// when data is received over stop channel, just exit the goroutine
			case <-s.logHandler.stop:
				return
			default: // do nothing
			}
		}
	}(s)

	return nil
}

// Clear clears the log (empty implementation to satisfy the interface, only file logger needs this one...).
func (s *SyslogHandler) Clear() error {
	return nil
}

// NewSyslogHandler creates a new sysloh handler.
func NewSyslogHandler(ip, fmt string, sev Severity) *SyslogHandler {
	return &SyslogHandler{newLogHandler(fmt, sev), ip, NewSyslogMsg()}
}
