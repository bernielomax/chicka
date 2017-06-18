package exec

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
	"os"
	"strconv"
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

// Test is a struct for defining test settings.
type Test struct {
	Command  string   `json:"command"`
	Args     []string `json:"args"`
	Interval int      `json:"interval"`
}

// Tests is a slice of test.
type Tests []Test

// Result is a struct for storing plugin exection results.
type Result struct {
	TestCommand string      `json:"test_command"`
	Expect      bool        `json:"expect"`
	Result      bool        `json:"result"`
	Data        interface{} `json:"data"`
	Description string      `json:"description"`
}

// Results is a slice of result.
type Results map[string]Result

// NewController sets up the controller for managing tests.
func NewController() *Controller {

	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)

	return &Controller{
		Ctx:    ctx,
		Cancel: cancel,
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// Reset sets a new context with cancel support into the controller.
func (ctrl *Controller) Reset() {
	ctx := context.Background()
	ctrl.Ctx, ctrl.Cancel = context.WithCancel(ctx)
}

// Validate is used to validate the test.
func (t *Test) Validate() error {

	if t.Interval < 5 {
		return errValidationIntervalLength
	}

	return nil
}

// Run is used to execute all tests defined in the configuration.
func (ctrl *Controller) Run(cfg *Config, c *cache.Cache, l LoggerSvc, e ErrorSvc) {

	for {

		total := len(cfg.Tests)

		done := make(chan bool, total)

		for _, test := range cfg.Tests {

			go func(test Test) {

				run := true

				for run {

					if test.Interval < validSleepIntervalLength {
						test.Interval = 60
					}

					select {
					case <-time.After(time.Duration(test.Interval) * time.Second):

						err := test.Validate()
						if err != nil {
							e.Send(test.Command, err)
							continue
						}

						r := Result{
							TestCommand: test.Command,
						}

						args := strings.Split(test.Command, " ")

						command := args[0]

						args = append(args[:0], args[0+1:]...)

						cmd := exec.Command(fmt.Sprintf("%v%v", cfg.Plugins.Path, command), args...)

						reader, err := cmd.StdoutPipe()
						if err != nil {
							e.Send(test.Command, err)
							continue
						}

						err = cmd.Start()
						if err != nil {
							e.Send(test.Command, err)
							continue
						}

						buf := new(bytes.Buffer)
						buf.ReadFrom(reader)

						err = cmd.Wait()

						output := buf.String()

						if err != nil {
							e.Send(test.Command, fmt.Errorf("error: %v, output: %v", err, output))
							continue
						}

						err = json.Unmarshal([]byte(output), &r)
						if err != nil {
							e.Send(test.Command, fmt.Errorf("error: %v, output: %v", err, output))
							continue
						}

						l.Send(r)

						c.Set(strconv.Itoa(int(time.Now().UnixNano())), r, cache.DefaultExpiration)

					case <-ctrl.Ctx.Done():
						run = false
						done <- true
					}
				}

			}(test)
		}

		for i := 0; i < total; i++ {
			<-done
		}

		ctrl.Reset()
	}
}
