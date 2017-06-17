package cmd

import (
	"fmt"
	"github.com/bernielomax/chicka/exec"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"runtime"
)

func init() {

	runtime.GOMAXPROCS(runtime.NumCPU())

}

var RootCmd = &cobra.Command{
	Use:   "chicka",
	Short: "Chicka is pluggable monitoring system written in Go",

	Run: runRootCmd,
}

func runRootCmd(cmd *cobra.Command, args []string) {
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/ckicka/")
	viper.AddConfigPath("$HOME/.chicka")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	errs := make(chan error)

	cfg, err := exec.ReadConfig()
	if err != nil {
		errs <- err
	}

	c := exec.NewController()

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		err := cfg.Refresh()
		if err != nil {
			errs <- err
		}
		c.Cancel()
	})

	l := exec.NewLoggerSvc()

	l.Listen()

	e := exec.NewErrorSvc()

	e.Listen()

	c.Run(cfg, l, e)

	exitOnError(<-errs)
}

func exitOnError(err error) {
	fmt.Println("Error:", err)
	os.Exit(1)
}
