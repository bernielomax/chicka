package cmd

import (
	"fmt"
	"github.com/bernielomax/chicka/exec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"io/ioutil"
	"os"
	"regexp"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get latest chicka plugins repo",
	RunE:  getCmdRun,
}

func init() {
	getCmd.Flags().StringP("ssh-key", "k", fmt.Sprintf("%v/.ssh/id_rsa", os.Getenv("HOME")), "The SSH key to use for Git authentication")
	viper.BindPFlag("ssh_key", getCmd.Flags().Lookup("ssh-key"))
}

func getCmdRun(cmd *cobra.Command, args []string) error {

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	cfg, err := exec.ReadConfig()
	if err != nil {
		return err
	}

	if _, err := os.Stat(cfg.Plugins.Path); os.IsNotExist(err) {

		os.MkdirAll(cfg.Plugins.Path, 0755)

		b, err := regexp.MatchString("^https?.*", cfg.Git.URL)
		if err != nil {
			return err
		}

		if b {

			r, err := git.PlainClone(cfg.Plugins.Path, false, &git.CloneOptions{
				URL:  cfg.Git.URL,
			})

			if err != nil {
				return err
			}

			w, err := r.Worktree()
			if err != nil {
				return err
			}

			err = w.Checkout(&git.CheckoutOptions{
				Hash: plumbing.NewHash(""),
			})
			if err != nil {
				return err
			}

		} else {

			key, err := ioutil.ReadFile(viper.GetString("ssh_key"))
			if err != nil {
				return err
			}

			signer, err := ssh.ParsePrivateKey(key)
			if err != nil {
				return err
			}

			auth := &gitssh.PublicKeys{User: "git", Signer: signer}

			r, err := git.PlainClone(cfg.Plugins.Path, false, &git.CloneOptions{
				URL:  cfg.Git.URL,
				Auth: auth,
			})

			if err != nil {
				return err
			}

			w, err := r.Worktree()
			if err != nil {
				return err
			}

			err = w.Checkout(&git.CheckoutOptions{
				Hash: plumbing.NewHash(""),
			})
			if err != nil {
				return err
			}

		}

	}

	return nil
}
