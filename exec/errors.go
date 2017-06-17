package exec

type errSvc struct {
	Queue chan error
}

type ErrorSvc interface {
	Send(err error)
	Listen()
}

func NewErrorSvc() ErrorSvc {
	return ErrorSvc(&errSvc{
		Queue: make(chan error),
	})
}

func (e *errSvc) Send(err error) {
	e.Queue <- err
}

func (e *errSvc) Listen() {
	go func() {
		panic(<-e.Queue)
	}()
}
