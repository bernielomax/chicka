package exec

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

const (
	validSleepIntervalLength = 5
)

var (
	errValidationIntervalLength = fmt.Errorf("the interval must be greater than %d", validSleepIntervalLength)
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
	CheckCommand string      `json:"check_command"`
	Expect       bool        `json:"expect"`
	Status       bool        `json:"status"`
	Result       interface{} `json:"result"`
	Description  string      `json:"description"`
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

					if check.Interval < validSleepIntervalLength {
						check.Interval = 60
					}

					select {
					case <-time.After(time.Duration(check.Interval) * time.Second):

						err := check.Validate()
						if err != nil {
							e.Send(check.Command, err)
							continue
						}

						r := Result{
							CheckCommand: check.Command,
						}

						args := strings.Split(check.Command, " ")

						command := args[0]

						args = append(args[:0], args[0+1:]...)

						cmd := exec.Command(fmt.Sprintf("%v%v", cfg.Plugins.Path, command), args...)

						reader, err := cmd.StdoutPipe()
						if err != nil {
							e.Send(check.Command, err)
							continue
						}

						err = cmd.Start()
						if err != nil {
							e.Send(check.Command, err)
							continue
						}

						buf := new(bytes.Buffer)
						buf.ReadFrom(reader)

						err = cmd.Wait()

						output := buf.String()

						if err != nil {
							e.Send(check.Command, fmt.Errorf("error: %v, output: %v", err, output))
							continue
						}

						err = json.Unmarshal([]byte(output), &r)
						if err != nil {
							e.Send(check.Command, fmt.Errorf("error: %v, output: %v", err, output))
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
