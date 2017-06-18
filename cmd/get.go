package cmd

import (
	"github.com/bernielomax/chicka/exec"
	"github.com/spf13/cobra"
	"gopkg.in/src-d/go-git.v4"
	"os"
	"fmt"
	"github.com/spf13/viper"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get latest chicka plugins repo",
	RunE:  getCmdRun,
}

func getCmdRun(cmd *cobra.Command, args []string) error {

	fmt.Println(viper.AllKeys())

	cfg, err := exec.ReadConfig()
	if err != nil {
		exitOnError(err)
	}

	fmt.Println(cfg)

	_, err = git.PlainClone(cfg.Plugins.Path, true, &git.CloneOptions{
		URL:      cfg.Git.URL,
		Progress: os.Stdout,
	})

	if err != nil {
		return err
	}

	return nil
}
