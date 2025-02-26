{{ if dig "VpcFlowLogs" "Enable" false . }}
{{ template "doNotEdit" . }}

data "aws_cloudwatch_log_group" "vpc_flow_logs" {
  name = data.aws_ssm_parameter.vpc_flow_logs_log_group_name.value
}

data "aws_ssm_parameter" "vpc_flow_logs_log_group_name" {
  name  = "/${local.environment}/vpc/log-group-name"
}
{{ else }}
# x-boilerplate-delete
{{ end -}}