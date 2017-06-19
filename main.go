package main

import (
	"fmt"
	"github.com/bernielomax/chicka/cmd"
	"os"
)

func exitOnError(err error) {
	fmt.Println("ERROR:", err)
	os.Exit(1)
}

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		exitOnError(err)
	}
}
