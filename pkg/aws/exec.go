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

func listServices(clusterArn string) ([]string, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		fmt.Printf("error loading configuration: %v\n", err)
		return nil, err
	}
	ecsSvc := ecs.NewFromConfig(cfg)
	result, err := ecsSvc.ListServices(context.TODO(), &ecs.ListServicesInput{
		Cluster: &clusterArn,
	})
	if err != nil {
		return nil, fmt.Errorf("listing services: %w", err)
	}
	var services []string
	services = append(services, result.ServiceArns...)
	return services, nil
}

func Exec() (string, error) {
	var cluster string
	var service string

	clusters, err := listClusters()
	if err != nil {
		return "", fmt.Errorf("listing clusters: %w", err)
	}

	clusterOptions := make([]huh.Option[string], 0)
	for _, c := range clusters {
		clusterOptions = append(clusterOptions, huh.NewOption(c, c))
	}

	clusterSelect := huh.NewSelect[string]().
		Options(clusterOptions...).
		Title("Select cluster").
		Value(&cluster)

	err = huh.NewForm(huh.NewGroup(clusterSelect)).Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return "", nil
		} else {
			return "", fmt.Errorf("running cluster select form: %w", err)
		}
	}

	services, err := listServices(cluster)
	if err != nil {
		return "", fmt.Errorf("listing services: %w", err)
	}

	serviceOptions := make([]huh.Option[string], 0)
	for _, s := range services {
		serviceOptions = append(serviceOptions, huh.NewOption(s, s))
	}

	serviceSelect := huh.NewSelect[string]().
		Options(serviceOptions...).
		Title("Select service").
		Value(&service)

	err = huh.NewForm(huh.NewGroup(serviceSelect)).Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return "", nil
		} else {
			return "", fmt.Errorf("running service select form: %w", err)
		}
	}

	return fmt.Sprintf("Cluster: %s, Service: %s", cluster, service), nil
}
