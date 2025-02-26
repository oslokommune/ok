{{ if dig "OpenTelemetrySidecar" "Enable" false . -}}
{{ template "doNotEdit" . }}
module "aws_opentelemetry_collector" {

  # https://github.com/oslokommune/golden-path-iac/tree/main/terraform/modules/aws-otel-collector-ecs-sidecar
  source = "git@github.com:oslokommune/golden-path-iac//terraform/modules/aws-otel-collector-ecs-sidecar?ref=aws-otel-collector-ecs-sidecar-v1.1.1"

  name     = local.environment
  app_name = local.main_container_name

  otel_image                     = local.otel_collector_image_uri
  otel_collector_log_group_name  = data.aws_ssm_parameter.otel_collector_log_group_name.insecure_value
  otel_config_ssm_parameter_name = "/${local.environment}/${local.main_container_name}/otel/config"
  otel_enable_resource_detection = true

  ################################################################################
  # Receiver configuration for AWS ECS Container Metrics
  ################################################################################
  enable_awsecscontainermetrics = true

  ################################################################################
  # Receiver configuration for Prometheus
  ################################################################################
  prometheus_scrape_target = local.prometheus_scrape_target
  prometheus_metrics_path  = local.prometheus_metrics_path

  ################################################################################
  # Exporter configuration for Prometheus
  ################################################################################
  prometheus_workspace = local.environment

{{- if dig "Xray" "Enable" false . }}

  ################################################################################
  # AWS X-Ray tracing
  ################################################################################
  enable_xray = true
{{- end }}

}
{{ else -}}
# x-boilerplate-delete
{{ end -}}
