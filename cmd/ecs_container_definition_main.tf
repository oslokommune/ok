locals {

  main_container = {

    ################################################################################
    # Container details
    ################################################################################
    name      = var.main_container_name
    image     = local.image_uri
    essential = true
    dependencies = [
      {
        containerName = "aws-opentelemetry-collector"
        condition     = "START"
      }
    ]

    ################################################################################
    # Port mappings
    ################################################################################
    port_mappings = [
      {
        name : "${var.main_container_name}-${var.main_container_port}"
        protocol : "tcp",
        containerPort : var.main_container_port
        hostPort : var.main_container_port
      }
    ],

    ################################################################################
    # Environment variables
    # See https://docs.aws.amazon.com/AmazonECS/latest/developerguide/taskdef-envfiles.html
    ################################################################################
    environment = [
      {
        name : "KTOR_LOG_LEVEL",
        value : "INFO"
      },
      {
        name : "BAR_SERVICE_URL",
        value : "http://km:8080"
      },
      {
        name  = "DB_NAME"
        value = "tootikkidb"
      },
      {
        name  = "WHATEVER"
        value = "sswhaiteiiisvers"
      }
    ]

    ################################################################################
    # Environment variables (sensitive)
    # See https://docs.aws.amazon.com/AmazonECS/latest/developerguide/secrets-envvar-ssm-paramstore.html#secrets-envvar-ssm-paramstore-update-container-definition
    ################################################################################
    secrets = [
      {
        name      = "DB_ENDPOINT"
        valueFrom = data.aws_ssm_parameter.db_endpoint.arn
      },
      {
        name      = "DB_USERNAME"
        valueFrom = aws_ssm_parameter.db_username.arn
      },
      {
        name      = "DB_PASSWORD"
        valueFrom = aws_ssm_parameter.db_password.arn
      }
    ]

    ################################################################################
    # HealthCheck
    ################################################################################
    health_check = {
      # WARNING: Curl MUST be installed inside the container image
      command = [
        "CMD-SHELL",
        "curl -f http://localhost:${var.main_container_port}${var.main_container_health_check_path} || exit 1"
      ]
      interval    = 20 # How often to check
      retries     = 3
      startPeriod = var.main_container_health_check_start_period
      timeout     = 2
    }

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
    readonly_root_filesystem = true

    ################################################################################
    # Monitoring and logging
    ################################################################################
    create_cloudwatch_log_group = false # Created outside of the module
    log_configuration = {
      logDriver = "awslogs"
      options = {
        "awslogs-group"         = aws_cloudwatch_log_group.ecs_service.name
        "awslogs-region"        = var.region
        "awslogs-stream-prefix" = "ecs"
      }
    }

  }

}