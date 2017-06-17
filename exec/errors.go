package exec

import (
	"io"
	"time"

	"github.com/go-kit/kit/log"
)

type errSvc struct {
	Queue  chan error
	Logger log.Logger
}

// ErrorSvc is the interface to use for handling errors.
type ErrorSvc interface {
	Send(err error)
	Listen()
}

// NewErrorSvc returns the error service interface.
func NewErrorSvc(w io.Writer) ErrorSvc {

	l := log.NewLogfmtLogger(w)

	l = log.With(l, "time", time.Now().String())

	l = log.With(l, "severity", "error")

	return ErrorSvc(&errSvc{
		Queue:  make(chan error),
		Logger: l,
	})
}

// Send adds errors the the error queue channel.
func (e *errSvc) Send(err error) {
	e.Queue <- err
}

// Listen waits and displays any errors sent to the error queue channel.
func (e *errSvc) Listen() {
	go func() {
		for {
			e.Logger.Log("error", <-e.Queue)
		}
	}()
}
