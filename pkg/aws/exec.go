package aws

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/charmbracelet/huh"
)

func listClusters() ([]string, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		fmt.Printf("error loading configuration: %v\n", err)
		return nil, err
	}

	ecsSvc := ecs.NewFromConfig(cfg)

	result, err := ecsSvc.ListClusters(context.TODO(), &ecs.ListClustersInput{})
	if err != nil {
		return nil, fmt.Errorf("listing clusters: %w", err)
	}

	var clusters []string
	clusters = append(clusters, result.ClusterArns...)

	return clusters, nil
}

func Exec() (string, error) {
	var cluster string

	clusters, _ := listClusters()

	options := make([]huh.Option[string], 0)

	for _, c := range clusters {
		options = append(options, huh.NewOption(c, c))
	}

	s := huh.NewSelect[string]().
		Options(options...).
		Title("Select cluster").
		Value(&cluster)

	err := huh.NewForm(huh.NewGroup(s)).Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return "", nil
		} else {
			return "", fmt.Errorf("running multi select form: %w", err)
		}
	}

	return cluster, nil
}
