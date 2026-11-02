package scriptrunner

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
)

//go:embed port-forward.sh
var scripts embed.FS

// Loads a script from the embedded filesystem, writes it to a temp file, and returns the path to the temp file...
func createTempScriptFile(scriptName string) (string, error) {
	scriptContent, err := scripts.ReadFile(scriptName)
	if err != nil {
		return "", fmt.Errorf("reading script: %w", err)
	}

	tmpFileName := fmt.Sprintf("*-%s", scriptName)
	tmpFile, err := os.CreateTemp("", tmpFileName)
	if err != nil {
		return "", fmt.Errorf("creating temp file: %w", err)
	}

	if _, err := tmpFile.Write(scriptContent); err != nil {
		return "", fmt.Errorf("writing to temp file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return "", fmt.Errorf("closing temp file: %w", err)
	}

	return tmpFile.Name(), nil
}

// Executes a script with the given arguments. Without interactive, the script's stdin is empty, so
// a prompt in it gets end of input instead of waiting.
func executeScript(scriptFile string, args []string, interactive bool) error {

	bashPath, err := exec.LookPath("bash")
	if err != nil {
		return fmt.Errorf("failed to find bash: %v", err)
	}

	combinedArgs := append([]string{scriptFile}, args...)

	cmd := exec.Command(bashPath, combinedArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if interactive {
		cmd.Stdin = os.Stdin
	}

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("executing script: %v", err)
	}

	return nil
}

func RunScript(scriptName string, args []string, interactive bool) error {
	scriptFile, err := createTempScriptFile(scriptName)
	if err != nil {
		return err
	}
	defer os.Remove(scriptFile)

	return executeScript(scriptFile, args, interactive)
}
