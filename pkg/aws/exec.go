package aws

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/magefile/mage/sh"
	"strings"
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

func listTasks(clusterArn, serviceArn string) ([]string, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		fmt.Printf("error loading configuration: %v\n", err)
		return nil, err
	}
	ecsSvc := ecs.NewFromConfig(cfg)
	result, err := ecsSvc.ListTasks(context.TODO(), &ecs.ListTasksInput{
		Cluster:     &clusterArn,
		ServiceName: &serviceArn,
	})
	if err != nil {
		return nil, fmt.Errorf("listing tasks: %w", err)
	}
	var tasks []string
	tasks = append(tasks, result.TaskArns...)
	return tasks, nil
}

func getTaskDetails(clusterName, taskId string) (*types.Task, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error loading configuration: %w", err)
	}

	ecsSvc := ecs.NewFromConfig(cfg)

	input := &ecs.DescribeTasksInput{
		Cluster: &clusterName,
		Tasks:   []string{taskId},
	}

	result, err := ecsSvc.DescribeTasks(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("error describing task: %w", err)
	}

	if len(result.Tasks) == 0 {
		return nil, fmt.Errorf("no task found")
	}

	return &result.Tasks[0], nil
}

func outputExecuteCommand(clusterName, taskId, containerName string) error {
	combinedArgs := []string{
		"ecs",
		"execute-command",
		"--cluster", clusterName,
		"--task", taskId,
		"--container", containerName,
		"--command", "/bin/sh",
		"--interactive",
	}
	green := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	fmt.Println("------------------------------------------------------------------------------------------")
	fmt.Println("Running aws command:")
	fmt.Println(green.Render("aws " + strings.Join(combinedArgs, " ")))
	fmt.Println("------------------------------------------------------------------------------------------")

	_ = sh.RunV("aws", combinedArgs...)
	return nil
}

func Exec() (string, error) {
	var cluster string
	var service string
	var task string
	var container string

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

	tasks, err := listTasks(cluster, service)
	if err != nil {
		return "", fmt.Errorf("listing tasks: %w", err)
	}

	taskOptions := make([]huh.Option[string], 0)
	for _, t := range tasks {
		taskOptions = append(taskOptions, huh.NewOption(t, t))
	}

	taskSelect := huh.NewSelect[string]().
		Options(taskOptions...).
		Title("Select running task").
		Value(&task)

	err = huh.NewForm(huh.NewGroup(taskSelect)).Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return "", nil
		} else {
			return "", fmt.Errorf("running task select form: %w", err)
		}
	}

	taskDetail, err := getTaskDetails(cluster, task)
	if err != nil {
		return "", fmt.Errorf("getting task details: %w", err)
	}

	containerOptions := make([]huh.Option[string], 0)
	for _, c := range taskDetail.Containers {
		containerOptions = append(containerOptions, huh.NewOption(*c.Name, *c.Name))
	}

	containerSelect := huh.NewSelect[string]().
		Options(containerOptions...).
		Title("Select container").
		Value(&container)

	err = huh.NewForm(huh.NewGroup(containerSelect)).Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return "", nil
		} else {
			return "", fmt.Errorf("running container select form: %w", err)
		}
	}

	outputExecuteCommand(cluster, task, container)
	return "", nil
}
