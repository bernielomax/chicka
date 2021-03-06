package cmd

import (
	"os"
	"runtime"
	"time"

	"github.com/bernielomax/chicka/exec"
	"github.com/bernielomax/chicka/http"
	"github.com/fsnotify/fsnotify"
	"github.com/go-kit/kit/log"
	"github.com/patrickmn/go-cache"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {

	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/chicka/")
	viper.AddConfigPath("$HOME/.chicka")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	RootCmd.AddCommand(getCmd)

	runtime.GOMAXPROCS(runtime.NumCPU())

}

// RootCmd is the base command config for the cli.
var RootCmd = &cobra.Command{
	Use:   "chicka",
	Short: "Chicka is pluggable monitoring system written in Go",

	RunE: runRootCmd,
}

func runRootCmd(cmd *cobra.Command, args []string) error {

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	errs := make(chan error)

	cfg, err := exec.ReadConfig()
	if err != nil {
		return err
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

	cache := cache.New(cfg.Cache.TTL*time.Second, cfg.Cache.TTL*time.Second)

	file, err := os.OpenFile(cfg.Logging.Path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		errs <- err
	}

	defer file.Close()

	l := exec.NewLoggerSvc(log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout)), log.NewJSONLogger(log.NewSyncWriter(file)))

	l.Listen()

	e := exec.NewErrorSvc(log.NewSyncWriter(os.Stderr))

	e.Listen()

	go http.StartAPIServer(cfg.HTTP.API, cache)

	go http.StartFrontEndServer(cfg.HTTP.WWW)

	c.Run(cfg, cache, l, e)

	return nil

}
