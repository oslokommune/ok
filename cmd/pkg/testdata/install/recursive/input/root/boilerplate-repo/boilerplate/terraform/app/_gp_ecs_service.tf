{{ template "doNotEdit" . }}
module "ecs_service" {

  # https://github.com/terraform-aws-modules/terraform-aws-ecs
  # https://github.com/terraform-aws-modules/terraform-aws-ecs/tree/master/modules/service
  source  = "terraform-aws-modules/ecs/aws//modules/service"
  version = "5.12.0"

  create = true

  ################################################################################
  # Environment
  # https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_definition_parameters.html#task_size
  ################################################################################
  cluster_arn      = data.aws_ecs_cluster.this.id
  launch_type      = "FARGATE" # https://docs.aws.amazon.com/AmazonECS/latest/developerguide/AWS_Fargate.html
  runtime_platform = local.runtime_platform
  cpu              = local.cpu
  memory           = local.memory
{{- if and .AppEcsExec .AppReadOnlyRootFileSystem }}
  volume = {
    for mount_point in concat(local.mount_points, local.mount_points_aws_ecs_exec) : mount_point.sourceVolume => {}
  }
{{- else }}
  volume = {
    for mount_point in local.mount_points: mount_point.sourceVolume => {}
  }
{{- end }}

  ################################################################################
  # Task definition
  # See https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_definitions.html
  ################################################################################
  ignore_task_definition_changes = false
{{- if dig "OpenTelemetrySidecar" "Enable" false . }}
  container_definitions = merge({
    main_container = local.main_container
    otel_container = module.aws_opentelemetry_collector.otel_container_definition_terraform_aws_modules_format
  }, local.additional_container_definitions)
{{- else }}
  container_definitions = merge({
    main_container = local.main_container
  }, local.additional_container_definitions)
{{- end }}

  ################################################################################
  # Deployment configuration
  # See https://docs.aws.amazon.com/AmazonECS/latest/developerguide/deployment-type-ecs.html
  ################################################################################
  name                               = local.main_container_name
  scheduling_strategy                = "REPLICA" # Must be REPLICA for FARGATE
  deployment_minimum_healthy_percent = 100
  deployment_maximum_percent         = 200
  # If Amazon ECS Exec is enabled, readonly_root_filesystem in the container definition must be false
  # See https://docs.aws.amazon.com/AmazonECS/latest/developerguide/ecs-exec.html
  # See https://github.com/aws-containers/amazon-ecs-exec-checker/issues/21
  enable_execute_command = {{ .AppEcsExec }} # Amazon ECS Exec

  ################################################################################
  # Service auto scaling
  # See https://docs.aws.amazon.com/AmazonECS/latest/developerguide/service-auto-scaling.html
  ################################################################################
  enable_autoscaling       = true
  desired_count            = local.desired_count
  autoscaling_min_capacity = local.autoscaling_min_capacity
  autoscaling_max_capacity = local.autoscaling_max_capacity
{{- if dig "DailyShutdown" "Enable" false . }}
  autoscaling_scheduled_actions = {
    daily_shutdown = {
      min_capacity = 0
      max_capacity = 0
      schedule     = "cron(0 16 * * ? *)" # At 16:00 UTC every day (18:00 UTC+2)
    }
    daily_startup = {
      min_capacity = local.autoscaling_min_capacity
      max_capacity = local.autoscaling_max_capacity
      schedule     = "cron(0 6 * * ? *)" # At 06:00 UTC every day (08:00 UTC+2)
    }
  }
{{- end }}

  ################################################################################
  # Networking
  # See https://docs.aws.amazon.com/AmazonECS/latest/userguide/fargate-task-networking.html
  ################################################################################
  subnet_ids            = module.data_networking.aws_private_subnet_ids
  assign_public_ip      = false
  create_security_group = false
  security_group_ids    = [module.sg_ecs_app.security_group_id]

{{- if dig "AlbHostRouting" "Enable" false . }}

  ################################################################################
  # Load balancing
  # See https://docs.aws.amazon.com/AmazonECS/latest/developerguide/service-load-balancing.html
  ################################################################################
  load_balancer = local.load_balancer
{{- end }}

{{- if dig "ServiceConnect" "Enable" false . }}

  ################################################################################
  # Service Connect
  # See https://docs.aws.amazon.com/AmazonECS/latest/developerguide/service-connect.html
  ################################################################################
  service_connect_configuration = {
    enabled   = true
    namespace = local.environment

    service = {
      discovery_name = local.main_container_name

      client_alias = {
        dns_name = local.main_container_name
        port     = local.main_container_port
      }

      port_name = "${local.main_container_name}-${local.main_container_port}"
    }

    log_configuration = {
      log_driver = "awslogs"
      options    = {
        "awslogs-group"         = data.aws_ssm_parameter.service_connect_log_group_name.insecure_value
        "awslogs-region"        = local.region
        "awslogs-stream-prefix" = "ecs"
      }
    }

  }

{{- end }}

  ################################################################################
  # Task execution IAM role
  # See https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_execution_IAM_role.html
  ################################################################################
  create_task_exec_iam_role          = true
  task_exec_iam_role_use_name_prefix = true
  task_exec_iam_role_path            = "/${local.environment}/${local.main_container_name}/"
  task_exec_iam_role_name            = substr("ecs-task-exec-${local.main_container_name}", 0, 36)
  task_exec_iam_role_description     = "Grants the Amazon ECS container and Fargate agents permissions to make AWS API calls on the behalf of the application \"${local.main_container_name}\" in the environment \"${local.environment}\"."
{{- if dig "OpenTelemetrySidecar" "Enable" false . }}
  task_exec_iam_role_policies        = {
    opentelemetry_collector = module.aws_opentelemetry_collector.task_exec_policy_arn
  }
{{- end }}
  task_exec_ssm_param_arns = concat(
    [for s in local.secrets : s.valueFrom if strcontains(s.valueFrom, "ssm")]
{{- if dig "DatabaseConnectivity" "Enable" false . }},
    [data.aws_ssm_parameter.db_endpoint.arn]
{{- end }}
    )
  task_exec_secret_arns = [for s in local.secrets : join(":", slice(split(":", s.valueFrom), 0, 7)) if strcontains(s.valueFrom, "secretsmanager")]

  ################################################################################
  # Task IAM role
  # See https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-iam-roles.html
  ################################################################################
  create_tasks_iam_role          = true
  tasks_iam_role_use_name_prefix = true
  tasks_iam_role_path            = "/${local.environment}/${local.main_container_name}/"
  tasks_iam_role_name            = substr("ecs-task-${local.main_container_name}", 0, 36)
  tasks_iam_role_description     = "Grants the containers running inside the ECS task permissions to make AWS API calls on the behalf of the application \"${local.main_container_name}\" in the environment \"${local.environment}\"."
{{- if dig "OpenTelemetrySidecar" "Enable" false . }}
  tasks_iam_role_policies        = {
    opentelemetry_collector = module.aws_opentelemetry_collector.task_policy_arn
  }
{{- end }}


}

################################################################################
# Outputs
################################################################################

{{ if not .ExampleImage.Enable -}}
output "image_digest" {
  value = var.main_container_image_digest
}

output "image_tag" {
  value = var.main_container_image_tag
}
{{- end }}

output "ecr_repository_url" {
  value = data.aws_ecr_repository.app.repository_url
}

{{- if dig "Subdomain" "Enable" false .AlbHostRouting }}

output "service_url" {
  value = "https://${local.main_container_name}.${local.route53_zone_name}"
}

{{- end }}

{{- if dig "ApexDomain" "Enable" false .AlbHostRouting }}

output "service_url_apex_domain" {
  value = "https://${local.route53_apex_zone_name}"
}

{{- end }}
