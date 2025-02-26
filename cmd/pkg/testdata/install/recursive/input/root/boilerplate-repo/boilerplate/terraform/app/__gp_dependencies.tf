{{ template "doNotEdit" . }}
################################################################################
# AWS
################################################################################

data "aws_ec2_managed_prefix_list" "s3" {
  name = "com.amazonaws.${local.region}.s3"
}

################################################################################
# networking
################################################################################

module "data_networking" {

  # https://github.com/oslokommune/golden-path-iac/tree/main/terraform/modules/data-networking
  source = "git@github.com:oslokommune/golden-path-iac//terraform/modules/data-networking?ref=data-networking-v0.2.0"

  project_name = local.common_tags.Environment

}

data "aws_subnet" "cidr" {
  for_each = toset(module.data_networking.aws_private_subnet_ids)
  id       = each.value
}

locals {
  private_cidr_blocks     = [for s in data.aws_subnet.cidr : s.cidr_block]
  csv_private_cidr_blocks = join(",", local.private_cidr_blocks)
}

################################################################################
# security-groups
################################################################################

{{- if dig "AlbHostRouting" "Enable" false . }}

data "aws_security_group" "alb" {
  tags = merge(
    local.common_tags,
    {
      "Name" = local.load_balancer_name
    }
  )
  vpc_id = module.data_networking.vpc_id
}

{{- end }}

{{- if dig "DatabaseConnectivity" "Enable" false . }}

data "aws_security_group" "db_main" {
  tags = merge(
    local.common_tags,
    {
      "Name"     = local.db_name
      "Database" = local.db_name
    }
  )
  vpc_id = module.data_networking.vpc_id
}

{{- end }}

################################################################################
# app-common
################################################################################

data "aws_ecs_cluster" "this" {
  cluster_name = local.environment
}

{{- if dig "Enable" false .AlbHostRouting }}

################################################################################
# dns
################################################################################

  {{- if dig "Subdomain" "Enable" false .AlbHostRouting }}

data "aws_route53_zone" "this" {
  name = local.route53_zone_name
}

  {{- end }}

  {{- if dig "ApexDomain" "Enable" false .AlbHostRouting }}

data "aws_route53_zone" "apex" {
  name = local.route53_apex_zone_name
}

  {{- end }}

################################################################################
# load-balancing
################################################################################

data "aws_lb" "this" {
  name = local.load_balancer_name
}

data "aws_lb_listener" "this" {
  load_balancer_arn = data.aws_lb.this.arn
  port              = 443
}

{{- end }}

{{- if dig "DatabaseConnectivity" "Enable" false . }}

################################################################################
# databases
################################################################################

data "aws_ssm_parameter" "db_endpoint" {
  name            = "/${local.environment}/database/${local.db_name}/db_endpoint"
  with_decryption = false
}

{{- end }}

################################################################################
# data stores
################################################################################

data "aws_ssm_parameter" "ecs_service_log_group_name" {
  name = "/${local.environment}/ecs/${local.main_container_name}/log-group-name"
}

{{- if dig "OpenTelemetrySidecar" "Enable" false . }}

data "aws_ssm_parameter" "otel_collector_log_group_name" {
  name = "/${local.environment}/ecs/${local.main_container_name}/otel-collector/log-group-name"
}

{{- end }}

{{- if dig "ServiceConnect" "Enable" false . }}

data "aws_ssm_parameter" "service_connect_log_group_name" {
  name = "/${local.environment}/ecs/${local.main_container_name}/service-connect/log-group-name"
}

{{- end }}

data "aws_ecr_repository" "app" {
  name                 = "${local.environment}-${local.main_container_name}"
}
