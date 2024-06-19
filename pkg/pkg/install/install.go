package install

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"
)

func Run() error {
	cmd, err := CreateBoilerplateCommand("packages.yml")
	if err != nil {
		return fmt.Errorf("creating boilerplate command: %w", err)
	}

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("running boilerplate command: %w", err)
	}

	return err
}

func CreateBoilerplateCommand(filePath string) (*exec.Cmd, error) {
	fmt.Println("Installing packages...")

	fileContents, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}

	err = yaml.Unmarshal(fileContents, &PackageManifest{})
	if err != nil {
		return nil, fmt.Errorf("unmarshalling YAML: %w", err)
	}

	cmd := exec.Command(
		"boilerplate",
		"--template-url", "%s/app/ref=v1.0.2",
		"--var-file", "common-config.yml",
		"--template-url", "app-hello.yml",
		"--output-folder", "app-hello",
		"--non-interactive",
	)

	return cmd, nil
}
