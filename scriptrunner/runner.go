package scriptrunner

import (
	"embed"
	"fmt"
	"os"
	"os/exec"

	"github.com/magefile/mage/sh"
)

//go:embed ok.sh port-forward.sh
var scripts embed.FS

// Loads a script from the embedded filesystem, writes it to a temp file, and returns the path to the temp file.
func createTempScriptFile(scriptName string) (string, error) {
	scriptContent, err := scripts.ReadFile(scriptName)
	if err != nil {
		return "", fmt.Errorf("error reading script: %w", err)
	}

	tmpFileName := fmt.Sprintf("*-%s", scriptName)
	tmpFile, err := os.CreateTemp("", tmpFileName)
	if err != nil {
		return "", fmt.Errorf("error creating temp file: %w", err)
	}

	if _, err := tmpFile.Write(scriptContent); err != nil {
		return "", fmt.Errorf("error writing to temp file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return "", fmt.Errorf("error closing temp file: %w", err)
	}

	return tmpFile.Name(), nil
}

// Executes a script with the given arguments, and returns the output.
func executeScript(scriptFile string, args []string) (string, error) {
	combinedArgs := append([]string{scriptFile}, args...)

	bashPath, err := exec.LookPath("bash")
	if err != nil {
		return "", fmt.Errorf("failed to find bash: %v", err)
	}

	output, err := sh.Output(bashPath, combinedArgs...)
	if err != nil {
		return "", fmt.Errorf("error executing script: %v", err)
	}

	return output, nil
}

func RunScript(scriptName string, args []string) {
	scriptFile, err := createTempScriptFile(scriptName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer os.Remove(scriptFile)

	output, err := executeScript(scriptFile, args)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(output)
}
