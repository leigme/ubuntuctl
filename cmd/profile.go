package cmd

import (
	"fmt"
	loki "github.com/leigme/loki/cobra"
	"github.com/spf13/cobra"
)

/*
Copyright © 2023 leig HERE <leigme@gmail.com>

*/

func init() {
	loki.Add(rootCmd, &profile{})
}

type profile struct{}

type profileHandle func()

func (p *profile) Execute() loki.Exec {
	return func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			p.handlers()[args[0]]()
		}
		fmt.Println("profile backup and restore")
	}
}

func (p *profile) handlers() map[string]profileHandle {
	return map[string]profileHandle{
		"show":    func() { fmt.Println("Show profile") },
		"init":    func() { fmt.Println("Initialize profile") },
		"backup":  func() { fmt.Println("Backup profile") },
		"restore": func() { fmt.Println("Restore profile") },
	}
}
