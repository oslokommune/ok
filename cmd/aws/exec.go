package aws

import (
	"fmt"

	"github.com/oslokommune/ok/pkg/aws"
	"github.com/spf13/cobra"
)

var EcsExecCommand = &cobra.Command{

	Use:   "ecs-exec",
	Short: "Get a shell to a running ECS task",

	RunE: func(cmd *cobra.Command, args []string) error {
		cluster, err := aws.Exec()
		if err != nil {
			return fmt.Errorf("list clusters: %w", err)
		}
		fmt.Println("Selected cluster: ", cluster)
		return nil
	},
}
