package scriptrunner

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
)

//go:embed ok.sh port-forward.sh
var scripts embed.FS

func RunScript(scriptName string, args []string) {
	scriptContent, err := scripts.ReadFile(scriptName)
	if err != nil {
		fmt.Println("Error reading script:", err)
		return
	}

	// Using fmt.Sprintf for string formatting
	tmpFileName := fmt.Sprintf("*-%s", scriptName)
	tmpFile, err := os.CreateTemp("", tmpFileName)
	if err != nil {
		fmt.Println("Error creating temp file:", err)
		return
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(scriptContent); err != nil {
		fmt.Println("Error writing to temp file:", err)
		return
	}
	if err := tmpFile.Close(); err != nil {
		fmt.Println("Error closing temp file:", err)
		return
	}

	cmdArgs := append([]string{"/bin/bash", tmpFile.Name()}, args...)
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("Error running script:", err)
	}
}
