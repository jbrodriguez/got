package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kong"

	"got/cli/cmd"
)

var cli struct {
	DataDir string `flag:"" short:"f" default:"${datadir}" help:"location of data files"`

	On     cmd.On     `cmd:"" help:"start the day"`
	Add    cmd.Add    `cmd:"" help:"add an activity"`
	Report cmd.Report `cmd:"" help:"report activities"`
}

func main() {
	home := os.Getenv("HOME")
	dataDir := filepath.Join(home, ".local", "share", "got")

	if !strings.HasSuffix(dataDir, "/") {
		dataDir += "/"
	}

	ctx := kong.Parse(&cli, kong.Vars{"datadir": dataDir})
	err := ctx.Run(&cmd.Context{DataDir: cli.DataDir})
	ctx.FatalIfErrorf(err)
}
