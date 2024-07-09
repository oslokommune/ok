package install

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"
	"strings"
)

const DefaultBaseUrl = "git@github.com:oslokommune/golden-path-boilerplate.git//boilerplate/terraform"

func Run() error {
	cmds, err := CreateBoilerplateCommands("packages.yml")
	if err != nil {
		return fmt.Errorf("creating boilerplate command: %w", err)
	}

	curDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting current directory: %w", err)
	}

	log.Debug().Msgf("Current working directory: %s", curDir)

	for _, cmd := range cmds {
		args := strings.Join(cmd.Args[1:], " ")
		log.Debug().Msgf("Running boilerplate command: %s %s", cmd.Path, args)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("running boilerplate command: %w", err)
		}
	}

	return nil
}

func CreateBoilerplateCommands(filePath string) ([]*exec.Cmd, error) {
	var BaseUrl = os.Getenv("BASE_URL")
	if BaseUrl == "" {
		BaseUrl = DefaultBaseUrl	
	}

	fmt.Println("Installing packages...")

	fileContents, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}

	var manifest PackageManifest

	err = yaml.Unmarshal(fileContents, &manifest)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling YAML: %w", err)
	}

	var cmds []*exec.Cmd
	for _, pkg := range manifest.Packages {
		templateURL := fmt.Sprintf("%s/%s?ref=%s", BaseUrl, pkg.Template, pkg.Ref)

		cmdArgs := []string{
			"--template-url", templateURL,
			"--output-folder", pkg.OutputFolder,
			"--non-interactive",
		}

		for _, varFile := range pkg.VarFiles {
			cmdArgs = append(cmdArgs, "--var-file", varFile)
		}

		cmd := exec.Command("boilerplate", cmdArgs...)
		cmds = append(cmds, cmd)
	}

	return cmds, nil
}
