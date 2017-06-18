package cmd

import (
	"fmt"
	"os"
	osExec "os/exec"
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

	runtime.GOMAXPROCS(runtime.NumCPU())

}

// RootCmd is the base command config for the cli.
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
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	errs := make(chan error)

	cfg, err := exec.ReadConfig()
	if err != nil {
		errs <- err
	}

	b, err := exec.PathExists(cfg.Plugins.Path)
	if err != nil {
		panic(err)
	}

	if !b {
		cmd := osExec.Command("git", "clone", cfg.Git.URL, cfg.Plugins.Path)
		err := cmd.Run()
		if err != nil {
			panic(err)
		}
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

	fmt.Println("HTTP", cfg.HTTP)

	go http.StartAPIServer(cfg.HTTP.APIAddr, cache)

	go http.StartFrontEndServer(cfg.HTTP.FrontendAddr)

	c.Run(cfg, cache, l, e)

}
