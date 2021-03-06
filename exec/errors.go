package exec

import (
	"io"

	"github.com/go-kit/kit/log"
)

type errSvc struct {
	Queue  chan ErrorResult
	Logger log.Logger
}

// ErrorResult is the struct used to store error information.
type ErrorResult struct {
	TestCommand string `json:"test_command"`
	Error       error  `json:"error"`
}

// ErrorSvc is the interface to use for handling errors.
type ErrorSvc interface {
	Send(testCommand string, err error)
	Listen()
}

// NewErrorSvc returns the error service interface.
func NewErrorSvc(w io.Writer) ErrorSvc {

	l := log.NewLogfmtLogger(w)

	l = log.With(l, "time", log.DefaultTimestampUTC)

	return ErrorSvc(&errSvc{
		Queue:  make(chan ErrorResult),
		Logger: l,
	})
}

// Send adds errors the the error queue channel.
func (e *errSvc) Send(testCommand string, err error) {
	r := ErrorResult{
		TestCommand: testCommand,
		Error:       err,
	}
	e.Queue <- r
}

// Listen waits and displays any errors sent to the error queue channel.
func (e *errSvc) Listen() {
	go func() {
		for {
			r := <-e.Queue
			e.Logger.Log("test_command", r.TestCommand, "error", r.Error.Error())
		}
	}()
}
