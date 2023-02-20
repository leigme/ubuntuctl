package cmd

import (
	"fmt"
	loki "github.com/leigme/loki/cobra"
	"github.com/leigme/loki/shell"
	"github.com/leigme/progressing"
	"github.com/spf13/cobra"
)

/*
Copyright © 2023 leig HERE <leigme@gmail.com>

*/

func init() {
	loki.Add(rootCmd, &install{
		s: shell.New(shell.WithProcess(progressing.New())),
		command: "sudo apt update" +
			" && sudo apt upgrad" +
			" && sudo apt install git -y" +
			" && sudo apt install curl -y" +
			" && sudo apt install wget -y" +
			" && sudo apt install tree -y" +
			" && sudo apt install zsh -y",
	})
}

type install struct {
	s       shell.Shell
	command string
}

func (i *install) Execute() loki.Exec {
	return func(cmd *cobra.Command, args []string) {
		fmt.Println(i.s.Exe(i.command))
	}
}
