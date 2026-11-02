#!/usr/bin/env bash

# Select a cluster using fzf
cluster=$(aws ecs list-clusters | jq -r '.clusterArns[]' | fzf --prompt "Select a cluster: ")
if [ -z "$cluster" ]; then
  echo "Cluster selection cancelled."
  exit 1
fi

# Extract just the cluster name from the ARN
cluster_name=$(basename "$cluster")

echo "Selected cluster: $cluster_name"

# Select a service using fzf
service=$(aws ecs list-services --cluster "$cluster_name" | jq -r '.serviceArns[]' | fzf --prompt "Select a service: ")
if [ -z "$service" ]; then
  echo "Service selection cancelled."
  exit 1
fi

# Extract just the service name from the ARN
service_name=$(basename "$service")

echo "Selected service: $service_name"

# List the first task ARN
task=$(aws ecs list-tasks --cluster "$cluster_name" --service-name "$service_name" | jq -r '.taskArns[]' | fzf --prompt "Select a task: ")
if [ -z "$task" ]; then
  echo "No tasks found."
  exit 1
fi

echo "Selected task: $task"

# Select a container using fzf
container=$(aws ecs describe-task-definition --task-definition "$service_name" | jq -r '.taskDefinition.containerDefinitions[].name' | fzf --prompt "Select a container: ")
if [ -z "$container" ]; then
  echo "Container selection cancelled."
  exit 1
fi

echo "Selected container: $container"

# # Execute command
aws ecs execute-command --cluster "$cluster_name" \
    --task "$task" \
    --container "$container" \
    --interactive \
    --command "/bin/sh"
