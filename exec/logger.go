package exec

import "github.com/go-kit/kit/log"

type loggerSvc struct {
	Queue   chan Result
	Loggers []log.Logger
}

// LoggerSvc is the interface for managing logging functions.
type LoggerSvc interface {
	Send(result Result)
	Listen()
}

// NewLoggerSvc returns a LoggerSvc interface.
func NewLoggerSvc(loggers ...log.Logger) LoggerSvc {

	for k, l := range loggers {
		l = log.With(l, "time", log.DefaultTimestampUTC)
		loggers[k] = l
	}

	return LoggerSvc(&loggerSvc{
		Queue:   make(chan Result),
		Loggers: loggers,
	})
}

// Listen binds to the logger queue channel and displays any logs messages that are sent.
func (l *loggerSvc) Listen() {

	go func() {
		for {
			r := <-l.Queue
			for _, logger := range l.Loggers {
				logger.Log("command", r.TestCommand, "description", r.Description, "data", r.Data, "result", r.Result, "expect", r.Expect)
			}
		}
	}()

}

// Send is used add log messages to the log queue channel.
func (l *loggerSvc) Send(result Result) {
	l.Queue <- result
}
