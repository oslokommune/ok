package main

import (
	"github.com/oslokommune/ok/cmd"
	"github.com/oslokommune/ok/pkg/version"
)

var (
	versionStr = "dev"
	date       = "unknown"
	commit     = "none"
)

func main() {

	//slog.SetLogLoggerLevel(slog.LevelDebug)

	version.Data = version.Version{
		Version: versionStr,
		Date:    date,
		Commit:  commit,
	}

	cmd.Execute()

}
