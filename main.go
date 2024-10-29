package main

import (
	"github.com/oslokommune/ok/cmd"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

var (
	version = "dev"
	date    = "unknown"
	commit  = "none"
)

func main() {

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	cmd.VersionData = cmd.Version{
		Version: version,
		Date:    date,
		Commit:  commit,
	}

	cmd.Execute()

}
