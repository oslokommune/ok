{{ template "doNotEdit" . }}
################################################################################
# Config
################################################################################
locals {

  ################################################################################
  # Basics
  ################################################################################
  main_container_name = "{{ .AppName }}"
{{- if dig "ExampleImage" "Enable" false . }}
  main_container_port = 80
{{- else }}
  main_container_port = 8080
{{- end }}

  ################################################################################
  # HealthCheck
  ################################################################################
{{- if dig "ExampleImage" "Enable" false . }}
  main_container_health_check_path = "/"
{{- else }}
  main_container_health_check_path = "/health"
{{- end }}
  main_container_health_check = {
    command = [
      "CMD-SHELL",
    {{- if dig "ExampleImage" "Enable" false . }}
      "wget -q -O- -t 1 http://localhost:${local.main_container_port}${local.main_container_health_check_path} || exit 1"
    {{- else }}
      "curl -f http://localhost:${local.main_container_port}${local.main_container_health_check_path} || exit 1"
    {{- end }}
    ]
    interval    = 20
    retries     = 3
    startPeriod = 5
    timeout     = 2
  }

  ################################################################################
  # External
  ################################################################################
{{- if dig "Subdomain" "Enable" false .AlbHostRouting }}
  # Your app will be available on a subdomain of this DNS zone
  route53_zone_name = "${local.environment}.oslo.systems"
{{- end }}

{{- if dig "ApexDomain" "Enable" false .AlbHostRouting }}
  # Your app will be available on exactly this address (the apex of this DNS zone)
  route53_apex_zone_name = "${local.main_container_name}.oslo.kommune.no"
{{- end }}

{{- if dig "DatabaseConnectivity" "Enable" false . }}
  db_name = "${local.environment}-main"
{{- end }}

  ################################################################################
  # Environment
  ################################################################################
  runtime_platform = {
    operating_system_family = "LINUX"
    cpu_architecture        = "X86_64"
  }
  cpu    = 512
  memory = 1024

  ################################################################################
  # Environment variables
  ################################################################################
  environment_variables = []
  secrets = []

  ################################################################################
  # Container details
  ################################################################################
  additional_main_container_dependencies = []

  additional_container_definitions = {}

  ################################################################################
  # Main container image
  ################################################################################
{{- $private_base_uri := "${local.account_id}.dkr.ecr.${local.region}.amazonaws.com/${local.environment}-ecr-public" }}
{{- $public_base_uri := "public.ecr.aws" }}
{{- $nginx_image_path := "nginx/nginx:1.27-alpine3.19-slim@sha256:44dbe45eb96afb7d83a08b738fc3b6218752e80546fa50830beae33cd54a2c70" }}
{{- $otel_image_path := "aws-observability/aws-otel-collector:v0.42.0@sha256:cd481b72f3b98710ba69c27dca5a329a9d57808d9c74b288cb26c177938cbcf1" }}

{{- if dig "ExampleImage" "Enable" false . }}
  {{- if dig "VpcEndpoints" "Enable" false . }}
    image_uri                = "{{ $private_base_uri }}/{{ $nginx_image_path }}"
  {{- else }}
    image_uri                = "{{ $public_base_uri }}/{{ $nginx_image_path }}"
  {{- end }}
{{- else }}
    image_uri = "${data.aws_ecr_repository.app.repository_url}@${var.main_container_image_digest}"
{{- end }}

{{- if dig "OpenTelemetrySidecar" "Enable" false . }}
  {{- if dig "VpcEndpoints" "Enable" false . }}
    otel_collector_image_uri = "{{ $private_base_uri }}/{{ $otel_image_path }}"
  {{- else }}
    otel_collector_image_uri = "{{ $public_base_uri }}/{{ $otel_image_path }}"
  {{- end }}
{{- end }}

  ################################################################################
  # Autoscaling
  ################################################################################

  desired_count            = 2
  autoscaling_min_capacity = 2
  autoscaling_max_capacity = 2

{{- if dig "IamForCicd" "Enable" false . }}

  ################################################################################
  # IAM for CI/CD
  ################################################################################
  app_gh_repo_name = "{{ .IamForCicd.AppGitHubRepo }}"
  additional_policies_for_github_actions_in_app_repo = []

  iac_gh_repo_name = "{{ .IamForCicd.IacGitHubRepo }}"
  additional_policies_for_github_actions_in_iac_repo = []

{{- end }}

{{- if and (hasKey .IamForCicd "AssumableCdRole") .IamForCicd.AssumableCdRole }}

  assume_cd_role_name = "${local.environment}-${local.main_container_name}-ecs-debug"

{{- end }}

  ################################################################################
  # Storage and volumes
  ################################################################################
  mount_points = []

{{- if and .AppEcsExec .AppReadOnlyRootFileSystem }}

  mount_points_aws_ecs_exec = [
    {
      sourceVolume  = "init-var-lib-amazon"
      containerPath = "/var/lib/amazon"
      readOnly      = false
    },
    {
      sourceVolume  = "init-var-log-amazon"
      containerPath = "/var/log/amazon"
      readOnly      = false
    }
  ]

{{- end }}

{{- if dig "Enable" false .AlbHostRouting  }}

  ################################################################################
  # Load balancing
  ################################################################################

  {{- if dig "Internal" false .AlbHostRouting }}
    load_balancer_name = substr("${local.environment}-main-internal", 0, 32)
  {{- else }}
    load_balancer_name = substr("${local.environment}-main-public", 0, 32)
  {{- end }}

  load_balancer = [
  {{- if dig "Subdomain" "Enable" false .AlbHostRouting }}
  {
    target_group_arn = module.alb_tg_host_routing_subdomain.target_group_arn
    container_name   = local.main_container_name
    container_port   = local.main_container_port
  },
  {{- end }}

  {{- if dig "ApexDomain" "Enable" false .AlbHostRouting }}
  {
    target_group_arn = module.alb_tg_host_routing_apex_domain.target_group_arn
    container_name   = local.main_container_name
    container_port   = local.main_container_port
  }
  {{- end }}
  ]

{{- end }}

{{- if dig "OpenTelemetrySidecar" "Enable" false . }}

  ################################################################################
  # Observability
  ################################################################################
  prometheus_scrape_target = "localhost:${local.main_container_port}"
  prometheus_metrics_path = "/metrics"

{{- end }}

  ################################################################################
  # Common
  ################################################################################
  {{ template "commonTags" }}

  {{ template "commonLocals" . }}

}
