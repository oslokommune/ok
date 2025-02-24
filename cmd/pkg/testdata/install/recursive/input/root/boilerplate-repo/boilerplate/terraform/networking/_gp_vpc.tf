{{ template "doNotEdit" . }}
module "args_vpc" {
  # https://github.com/oslokommune/golden-path-iac/tree/main/terraform/modules/args-vpc
  source = "git@github.com:oslokommune/golden-path-iac//terraform/modules/args-vpc?ref=args-vpc-v0.2.0"
  cidr   = local.vpc_cidr_block
}

locals {
  args_vpc = module.args_vpc.data.default
}

data "aws_availability_zones" "vpc_azs" {}

#tfsec:ignore:aws-ec2-no-excessive-port-access tfsec:ignore:aws-ec2-no-public-ingress-acl tfsec:ignore:
module "vpc" {

  # https://github.com/terraform-aws-modules/terraform-aws-vpc
  source  = "terraform-aws-modules/vpc/aws"
  version = "5.9.0"

  name        = local.environment # Replace if multiple VPCs
  cidr        = local.args_vpc.cidr
  enable_ipv6 = local.args_vpc.enable_ipv6

  enable_dns_hostnames = local.args_vpc.enable_dns_hostnames

  azs = data.aws_availability_zones.vpc_azs.names

  map_public_ip_on_launch            = local.args_vpc.map_public_ip_on_launch
  manage_default_network_acl         = local.args_vpc.manage_default_network_acl
  manage_default_security_group      = local.args_vpc.manage_default_security_group
  manage_default_route_table         = local.args_vpc.manage_default_route_table
  default_route_table_name           = local.environment
  create_database_subnet_route_table = local.args_vpc.create_database_subnet_route_table

  public_subnets                                               = local.args_vpc.public_subnets
  public_subnet_enable_resource_name_dns_aaaa_record_on_launch = local.args_vpc.public_subnet_enable_resource_name_dns_aaaa_record_on_launch
  public_subnet_enable_dns64                                   = local.args_vpc.public_subnet_enable_dns64
  public_subnet_ipv6_prefixes                                  = local.args_vpc.public_subnet_ipv6_prefixes
  public_subnet_assign_ipv6_address_on_creation                = local.args_vpc.public_subnet_assign_ipv6_address_on_creation
  public_subnet_tags                                           = local.args_vpc.public_subnet_tags

  private_subnets                                               = local.args_vpc.private_subnets
  private_subnet_enable_resource_name_dns_aaaa_record_on_launch = local.args_vpc.private_subnet_enable_resource_name_dns_aaaa_record_on_launch
  private_subnet_enable_dns64                                   = local.args_vpc.private_subnet_enable_dns64
  private_subnet_ipv6_prefixes                                  = local.args_vpc.private_subnet_ipv6_prefixes
  private_subnet_assign_ipv6_address_on_creation                = local.args_vpc.private_subnet_assign_ipv6_address_on_creation
  private_subnet_tags                                           = local.args_vpc.private_subnet_tags

  database_subnets                                               = local.args_vpc.database_subnets
  database_subnet_enable_resource_name_dns_aaaa_record_on_launch = local.args_vpc.database_subnet_enable_resource_name_dns_aaaa_record_on_launch
  database_subnet_enable_dns64                                   = local.args_vpc.database_subnet_enable_dns64
  database_subnet_ipv6_prefixes                                  = local.args_vpc.database_subnet_ipv6_prefixes
  database_subnet_assign_ipv6_address_on_creation                = local.args_vpc.database_subnet_assign_ipv6_address_on_creation
  database_subnet_tags                                           = local.args_vpc.database_subnet_tags

  intra_subnets                                               = local.args_vpc.intra_subnets
  intra_subnet_enable_resource_name_dns_aaaa_record_on_launch = local.args_vpc.intra_subnet_enable_resource_name_dns_aaaa_record_on_launch
  intra_subnet_enable_dns64                                   = local.args_vpc.intra_subnet_enable_dns64
  intra_subnet_ipv6_prefixes                                  = local.args_vpc.intra_subnet_ipv6_prefixes
  intra_subnet_assign_ipv6_address_on_creation                = local.args_vpc.intra_subnet_assign_ipv6_address_on_creation
  intra_subnet_tags                                           = local.args_vpc.intra_subnet_tags

  enable_nat_gateway = local.args_vpc.enable_nat_gateway
  single_nat_gateway = local.vpc_single_nat_gateway

  tags = merge(
    local.common_tags,
    { VpcName = local.environment }
  )

{{ if dig "VpcFlowLogs" "Enable" false . }}
  enable_flow_log                      = true
  create_flow_log_cloudwatch_log_group = false
  create_flow_log_cloudwatch_iam_role  = true
  flow_log_max_aggregation_interval    = 60

  flow_log_destination_type        = "cloud-watch-logs"
  flow_log_destination_arn         = data.aws_cloudwatch_log_group.vpc_flow_logs.arn
{{ end }}
}
