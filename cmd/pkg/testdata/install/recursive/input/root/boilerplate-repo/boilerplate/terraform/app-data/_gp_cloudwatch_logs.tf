{{ template "doNotEdit" . }}
###############################################################################
# CloudWatch Logs
# See https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/Working-with-log-groups-and-streams.html
################################################################################

resource "aws_cloudwatch_log_group" "ecs_service" {
  name_prefix       = "/${local.environment}/ecs/${local.main_container_name}/"
  retention_in_days = 90
  skip_destroy      = true
}

resource "aws_ssm_parameter" "ecs_service_log_group_name" {
  name  = "/${local.environment}/ecs/${local.main_container_name}/log-group-name"
  type  = "String"
  value = aws_cloudwatch_log_group.ecs_service.name
}

{{- if dig "OpenTelemetrySidecar" "Enable" false . }}

resource "aws_cloudwatch_log_group" "otel_collector" {
  name_prefix       = "/${local.environment}/ecs/${local.main_container_name}/otel-collector/"
  retention_in_days = 90
  skip_destroy      = true
}

resource "aws_ssm_parameter" "otel_collector_log_group_name" {
  name  = "/${local.environment}/ecs/${local.main_container_name}/otel-collector/log-group-name"
  type  = "String"
  value = aws_cloudwatch_log_group.otel_collector.name
}

{{- end }}

{{- if dig "ServiceConnect" "Enable" false . }}

resource "aws_cloudwatch_log_group" "service_connect" {
  name_prefix       = "/${local.environment}/ecs/${local.main_container_name}/service-connect/"
  retention_in_days = 90
  skip_destroy      = true
}

resource "aws_ssm_parameter" "service_connect_log_group_name" {
  name  = "/${local.environment}/ecs/${local.main_container_name}/service-connect/log-group-name"
  type  = "String"
  value = aws_cloudwatch_log_group.service_connect.name
}

{{- end }}

