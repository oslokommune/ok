{{ template "doNotEdit" . }}
################################################################################
# Security group for the ECS service
################################################################################

module "sg_ecs_app" {

  # https://github.com/terraform-aws-modules/terraform-aws-security-group
  source  = "terraform-aws-modules/security-group/aws"
  version = "5.3.0"

  name        = "${local.environment}-${local.main_container_name}-ecs-app"
  description = "Used by the ECS service named ${local.main_container_name}"
  vpc_id      = module.data_networking.vpc_id

{{- if dig "AlbHostRouting" "Enable" false . }}

  ingress_with_source_security_group_id = [
    {
      description              = "Allow inbound traffic from the security group associated with the public load balancer"
      from_port                = local.main_container_port
      to_port                  = local.main_container_port
      protocol                 = "tcp"
      source_security_group_id = data.aws_security_group.alb.id
    }
  ]

{{- end }}

{{- $description := "" }}
{{- $cidr_blocks := "" }}

{{- if dig "VpcEndpoints" "Enable" false . }}
  {{ $description = "Allow outbound traffic to private subnets in order to reach VPC interface endpoints" }}
  {{ $cidr_blocks = "local.csv_private_cidr_blocks" }}
{{- else }}
  {{ $description = "Allow outbound traffic to internet" }}
  {{ $cidr_blocks = "\"0.0.0.0/0\"" }}
{{- end -}}


  egress_with_cidr_blocks = [
    {
      description = "{{ $description }}"
      from_port   = 443
      to_port     = 443
      protocol    = "tcp"
      cidr_blocks = {{ $cidr_blocks }}
    }
  ]

{{- if dig "DatabaseConnectivity" "Enable" false . }}

  egress_with_source_security_group_id = [
    {
      description              = "Allow outbound traffic to the database"
      source_security_group_id = data.aws_security_group.db_main.id
      rule                     = "postgresql-tcp"
    }
  ]

{{- end }}

}

{{- if dig "VpcEndpoints" "Enable" false . }}

resource "aws_security_group_rule" "ecs_service_s3_https" {
  type              = "egress"
  description       = "Allow traffic to the S3 VPC gateway endpoint."
  from_port         = 443
  to_port           = 443
  protocol          = "tcp"
  security_group_id = module.sg_ecs_app.security_group_id
  # This is like pl-6da54004 (com.amazonaws.eu-west-1.s3) which is the same one used in the route table of the private subnets
  prefix_list_ids = [data.aws_ec2_managed_prefix_list.s3.id]
}

{{- end }}

{{- if dig "AlbHostRouting" "Enable" false . }}

################################################################################
# Add an ingress rule to the existing ALB security group
################################################################################
module "sg_rules_alb_public" {

  # https://github.com/terraform-aws-modules/terraform-aws-security-group
  source  = "terraform-aws-modules/security-group/aws"
  version = "5.3.0"

  create_sg         = false
  security_group_id = data.aws_security_group.alb.id
  egress_with_source_security_group_id = [
    {
      description              = "Allow outbound traffic to the ECS service ${local.main_container_name}"
      source_security_group_id = module.sg_ecs_app.security_group_id
      protocol                 = "tcp"
      from_port                = local.main_container_port
      to_port                  = local.main_container_port
    }
  ]

}

{{- end }}

{{- if dig "DatabaseConnectivity" "Enable" false . }}

################################################################################
# Add an ingress rule to the existing database security group
################################################################################

module "sg_rules_db_main" {

  # https://github.com/terraform-aws-modules/terraform-aws-security-group
  source  = "terraform-aws-modules/security-group/aws"
  version = "5.3.0"

  create_sg         = false
  security_group_id = data.aws_security_group.db_main.id
  ingress_with_source_security_group_id = [
    {
      description              = "Allow inbound traffic from the ECS service ${local.main_container_name} to the database"
      source_security_group_id = module.sg_ecs_app.security_group_id
      rule                     = "postgresql-tcp"
    }
  ]

}

{{- end }}
