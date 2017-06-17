package exec

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"time"
)

var (
	errValidationIntervalLength = errors.New("the interval must be greater than 5")
)

// Controller is a struct for managing executions.
type Controller struct {
	Ctx    context.Context
	Cancel context.CancelFunc
}

// Check is a struct for defining check settings.
type Check struct {
	Command  string   `json:"check"`
	Args     []string `json:"args"`
	Interval int      `json:"interval"`
}

// Result is a struct for storing plugin exection results.
type Result struct {
	Status      bool   `json:"status"`
	Description string `json:"description"`
}

// Checks is a slice of Check.
type Checks []Check

// NewController sets up the controller for managing checks.
func NewController() *Controller {

	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)

	return &Controller{
		Ctx:    ctx,
		Cancel: cancel,
	}
}

// Reset sets a new context with cancel support into the controller.
func (c *Controller) Reset() {
	ctx := context.Background()
	c.Ctx, c.Cancel = context.WithCancel(ctx)
}

// Validate is used to validate the check.
func (c *Check) Validate() error {

	if c.Interval < 5 {
		return errValidationIntervalLength
	}

	return nil
}

// Run is used to execute all checks defined in the configuration.
func (c *Controller) Run(cfg *Config, l LoggerSvc, e ErrorSvc) {

	for {

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
							continue
						}

						r := Result{}

						cmd := exec.Command(fmt.Sprintf("%v%v", cfg.Plugins.Path, check.Command))

						reader, err := cmd.StdoutPipe()
						if err != nil {
							e.Send(err)
							continue
						}

						err = cmd.Start()
						if err != nil {
							e.Send(err)
							continue
						}

						err = json.NewDecoder(reader).Decode(&r)
						if err != nil {
							e.Send(err)
							continue
						}

						err = cmd.Wait()
						if err != nil {
							e.Send(err)
							continue
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
