{{ template "doNotEdit" . }}
locals {

  main_container = {

    ################################################################################
    # Container details
    ################################################################################
    name      = local.main_container_name
    image     = local.image_uri
    essential = true

    dependencies = concat([
{{- if dig "OpenTelemetrySidecar" "Enable" false . }}
    {
      containerName = module.aws_opentelemetry_collector.otel_container_definition_terraform_aws_modules_format.name
      condition     = "START"
    }
{{- end }}
    ], local.additional_main_container_dependencies)

    ################################################################################
    # Port mappings
    ################################################################################
    port_mappings = [
      {
        name : "${local.main_container_name}-${local.main_container_port}"
        protocol : "tcp",
        containerPort : local.main_container_port
        hostPort : local.main_container_port
      }
    ]

    ################################################################################
    # Environment variables
    # See https://docs.aws.amazon.com/AmazonECS/latest/developerguide/taskdef-envfiles.html
    ################################################################################
{{- if dig "ExampleImage" "Enable" false . }}
    environment = local.environment_variables
{{- else }}
    environment = concat(local.environment_variables, [
      {
        name  = "IMAGE_DIGEST"
        value = var.main_container_image_digest
      }
  {{- if dig "Xray" "Enable" false . }},
      {
        name  = "OTEL_TRACES_SAMPLER"
        value = "xray"
      },
      {
        name  = "OTEL_RESOURCE_ATTRIBUTES"
        value = "service.name=${local.main_container_name},service.namespace=${local.environment}"
      }
  {{- end }}
    ])
{{- end }}

    ################################################################################
    # Environment variables from secrets (sensitive)
    # See https://docs.aws.amazon.com/AmazonECS/latest/developerguide/secrets-envvar-ssm-paramstore.html#secrets-envvar-ssm-paramstore-update-container-definition
    ################################################################################
    secrets = concat(local.secrets, [
{{- if dig "DatabaseConnectivity" "Enable" false . }}
      {
        name      = "DB_ENDPOINT"
        valueFrom = data.aws_ssm_parameter.db_endpoint.arn
      }
{{- end }}
    ])

    ################################################################################
    # HealthCheck
    ################################################################################
    health_check = local.main_container_health_check

    ################################################################################
    # Container size (this is optional) (should be lower than the task size, if set)
    ################################################################################
    # The total amount of cpu of all containers in a task will need to be lower than the task-level cpu value
    # See https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_definition_parameters.html#task_size
    # Container size (should be lower than the task size)
    cpu    = null
    memory = null

    ################################################################################
    # Storage and volumes
    ################################################################################
    readonly_root_filesystem = {{ .AppReadOnlyRootFileSystem }}
{{- if and .AppEcsExec .AppReadOnlyRootFileSystem }}
    mount_points             = concat(local.mount_points, local.mount_points_aws_ecs_exec)
{{- else }}
    mount_points             = local.mount_points
{{- end }}

    ################################################################################
    # Monitoring and logging
    ################################################################################
    create_cloudwatch_log_group = false # Created outside of the module
    log_configuration = {
      logDriver = "awslogs"
      options = {
        "awslogs-group"         = data.aws_ssm_parameter.ecs_service_log_group_name.insecure_value
        "awslogs-region"        = local.region
        "awslogs-stream-prefix" = "ecs"
      }
    }

  }

}
