package migrate

import (
	"fmt"
	"github.com/oslokommune/ok/pkg/pkg/common"
	"os/exec"
)

func Run(packages []common.Package) error {
	for _, p := range packages {
		for _, varFile := range p.VarFiles {
			fmt.Printf("Migrating YAML file: %s\n", varFile)

			// Example transformation command using yq - this assumes you want to perform some operation like changing a key's value
			// You will need to customize this command to your specific transformation requirements.
			// For instance: setting a value, renaming keys, or adding/removing fields.
			yqCmd := exec.Command("yq", "eval", `"your-transformation-expression"`, "-i", varFile)

			yqCmd.Stdout = nil
			yqCmd.Stderr = nil

			err := yqCmd.Run()
			if err != nil {
				return fmt.Errorf("error transforming YAML file %s: %w", varFile, err)
			}

			fmt.Printf("Successfully migrated YAML file: %s\n", varFile)
		}
	}

	return nil
}
