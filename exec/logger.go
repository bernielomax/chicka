package exec

import (
	"io"

	"github.com/go-kit/kit/log"
)

type loggerSvc struct {
	Queue  chan Result
	Logger log.Logger
}

// LoggerSvc is the interface for managing logging functions.
type LoggerSvc interface {
	Send(result Result)
	Listen()
}

// NewLoggerSvc returns a LoggerSvc interface.
func NewLoggerSvc(w io.Writer) LoggerSvc {

	l := log.NewLogfmtLogger(w)

	l = log.With(l, "time", log.DefaultTimestampUTC)

	l = log.With(l, "severity", "info")

	return LoggerSvc(&loggerSvc{
		Queue:  make(chan Result),
		Logger: l,
	})
}

// Listen binds to the logger queue channel and displays any logs messages that are sent.
func (l *loggerSvc) Listen() {

	go func() {
		for {
			r := <-l.Queue
			l.Logger.Log("check_command", r.CheckCommand, "description", r.Description, "status", r.Status, "result", r.Result, "expect", r.Expect)
		}
	}()

}

// Send is used add log messages to the log queue channel.
func (l *loggerSvc) Send(result Result) {
	l.Queue <- result
}
