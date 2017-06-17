package exec

import "fmt"

type loggerSvc struct {
	Queue chan Result
}

type LoggerSvc interface {
	Send(result Result)
	Listen()
}

func NewLoggerSvc() LoggerSvc {
	return LoggerSvc(&loggerSvc{
		Queue: make(chan Result),
	})
}

func (l *loggerSvc) Listen() {

	go func() {
		for true {
			result := <-l.Queue
			fmt.Println(result)
		}
	}()

}

func (l *loggerSvc) Send(result Result) {
	l.Queue <- result
}
