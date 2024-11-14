package main

import (
	"github.com/oslokommune/ok/cmd"
)

var (
	version = "dev"
	date    = "unknown"
	commit  = "none"
)

func main() {

	//slog.SetLogLoggerLevel(slog.LevelDebug)

	cmd.VersionData = cmd.Version{
		Version: version,
		Date:    date,
		Commit:  commit,
	}

	cmd.Execute()

}
