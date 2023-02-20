package cmd

/*
Copyright © 2023 leig HERE <leigme@gmail.com>

*/

import (
	_ "embed"
	"fmt"
	"github.com/leigme/loki/app"
	loki "github.com/leigme/loki/cobra"
	"github.com/leigme/loki/shell"
	"github.com/spf13/cobra"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//go:embed template/sources.list
var sourcesData string

func init() {
	loki.Add(rootCmd, &Sources{
		sourcesListDir:  "/etc/apt",
		sourcesListName: "sources.list",
		sourcesListBak:  "sources.list.bak",
		cmd:             shell.New(),
	})
}

type Sources struct {
	sourcesListDir, sourcesListName, sourcesListBak string
	cmd                                             shell.Shell
}

func (s *Sources) Execute() loki.Exec {
	return func(cmd *cobra.Command, args []string) {
		sourcesFile := filepath.Join(s.sourcesListDir, s.sourcesListName)
		command := fmt.Sprintf("cat %s", sourcesFile)
		if len(args) > 0 {
			switch args[0] {
			case "backup":
				command = backupSources(sourcesFile, filepath.Join(app.WorkDir(), s.sourcesListBak))
			case "update":
				lsb := strings.Split(fmt.Sprint(s.cmd.Exe("lsb_release -c")), ":")
				if len(lsb) < 2 {
					log.Fatal("lsb_release -c error")
				}
				sourcesList := os.Expand(sourcesData, func(s string) string {
					return strings.TrimSpace(lsb[1])
				})

				temp, err := os.CreateTemp("", fmt.Sprintf("%s.tmp", s.sourcesListName))
				if err != nil {
					log.Fatal(err)
				}
				defer func(name string) {
					err = os.Remove(name)
					if err != nil {
						log.Println(err)
					}
				}(temp.Name())
				_, err = io.WriteString(temp, sourcesList)
				if err != nil {
					log.Fatal(err)
				}
				err = temp.Close()
				if err != nil {
					log.Println(err)
				}
				command = updateSources(temp.Name(), sourcesList)
			case "restore":
				_, err := os.Stat(filepath.Join(app.WorkDir(), s.sourcesListBak))
				if err != nil {
					log.Fatal(err)
				}
				command = updateSources(filepath.Join(app.WorkDir(), s.sourcesListBak), sourcesFile)
			}
		}
		fmt.Println(s.cmd.Exe(command))
	}
}

func backupSources(src, dest string) (backCmd string) {
	backCmd = fmt.Sprintf("cp %s %s", src, dest)
	return
}

func updateSources(src, dest string) (updateCmd string) {
	updateCmd = fmt.Sprintf("sudo cp %s %s && sudo chmod a+r %s", src, dest, dest)
	return
}
