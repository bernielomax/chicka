package exec

import (
	"context"
	"errors"
	"time"
)

var (
	errValidationIntervalLength = errors.New("the interval must be greater than 5")
)

type Controller struct {
	Ctx    context.Context
	Cancel context.CancelFunc
}

type Check struct {
	Command  string   `json:"check"`
	Args     []string `json:"args"`
	Interval int      `json:"interval"`
}

type Result struct {
	Output   string `json:"output"`
	ExitCode int    `json:"exit_code"`
}

type Checks []Check

func NewController() *Controller {

	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)

	return &Controller{
		Ctx:    ctx,
		Cancel: cancel,
	}
}

func (c *Controller) Reset() {
	ctx := context.Background()
	c.Ctx, c.Cancel = context.WithCancel(ctx)
}

func (c *Check) Validate() error {

	if c.Interval < 5 {
		return errValidationIntervalLength
	}

	return nil
}

func (c *Controller) Run(cfg *Config, l LoggerSvc, e ErrorSvc) {

	for true {

		total := len(cfg.Checks)

		done := make(chan bool, total)

		for _, check := range cfg.Checks {

			go func(check Check) {

				run := true

				for run {

					select {
					case <-time.After(time.Duration(check.Interval) * time.Second):

						err := check.Validate()
						if err != nil {
							e.Send(err)
						}

						r := Result{
							Output: check.Command,
						}

						l.Send(r)

					case <-c.Ctx.Done():
						run = false
						done <- true
					}
				}

			}(check)
		}

		for i := 0; i < total; i++ {
			<-done
		}

		c.Reset()
	}
}
