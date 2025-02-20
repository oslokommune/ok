{{ if dig "VpcFlowLogs" "Enable" false . }}
{{ template "doNotEdit" . }}

resource "aws_cloudwatch_log_group" "vpc_flow_logs" {
  name_prefix       = "/${local.environment}/vpc/"
  retention_in_days = 90
  skip_destroy      = local.skip_destroy
}

resource "aws_ssm_parameter" "vpc_flow_logs_log_group_name" {
  name  = "/${local.environment}/vpc/log-group-name"
  type  = "String"
  value = aws_cloudwatch_log_group.vpc_flow_logs.name
}

output "vpc_flow_logs_log_group_name" {
  value = nonsensitive(aws_ssm_parameter.vpc_flow_logs_log_group_name.value)
}
{{ else }}
# x-boilerplate-delete
{{ end -}}