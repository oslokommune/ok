{{ if dig "VpcEndpoints" "Enable" false . }}
{{ template "doNotEdit" . }}
################################################################################
# VPC endpoints
# See https://docs.aws.amazon.com/whitepapers/latest/aws-privatelink/what-are-vpc-endpoints.html
################################################################################

module "vpc_endpoints" {

  # https://github.com/terraform-aws-modules/terraform-aws-vpc/tree/master/modules/vpc-endpoints
  source  = "terraform-aws-modules/vpc/aws//modules/vpc-endpoints"
  version = "5.9.0"

  create = true

  vpc_id = module.vpc.vpc_id

  ################################################################################
  # Security group
  # The security group attached to the VPC endpoints **must** allow incoming connections on port 443 from the private subnets of the VPC.
  ################################################################################
  create_security_group      = true
  security_group_name_prefix = "${local.environment}-vpc-endpoints-"
  security_group_description = "VPC endpoints security group"
  security_group_rules = {
    ingress_https = {
      description = "HTTPS from VPC"
      cidr_blocks = [module.vpc.vpc_cidr_block]
    }
  }

  endpoints = {

    ################################################################################
    # Interface endpoints
    # See https://docs.aws.amazon.com/vpc/latest/privatelink/create-interface-endpoint.html
    ################################################################################

    {{ if dig "VpcEndpoints" "Ecr" false . }}
    # Elastic Container Registry (ECR) API endpoint.
    # Used for calls to the Amazon ECR API. API actions such as DescribeImages and CreateRepository go to this endpoint.
    ecr_api = {
      service             = "ecr.api"
      private_dns_enabled = true # Optionally enabled
      subnet_ids          = module.vpc.private_subnets
      tags                = { Name = "${local.environment}-ecr-api" }
    }
    {{ end }}
    {{ if dig "VpcEndpoints" "Dkr" false . }}
    # Elastic Container Registry (ECR) Docker ("dkr") Registry API endpoint.
    # Docker client commands such as push and pull use this endpoint.
    ecr_dkr = {
      service             = "ecr.dkr"
      private_dns_enabled = true # MUST be enabled
      subnet_ids          = module.vpc.private_subnets
      tags                = { Name = "${local.environment}-ecr-dkr" }
    }
    {{ end -}}
    {{ if dig "VpcEndpoints" "Logs" false . }}
    # Amazon ECS tasks hosted on Fargate that pull container images from Amazon ECR that also use the `awslogs` log driver to send log information to CloudWatch Logs require the CloudWatch Logs VPC endpoint.
    logs = {
      service             = "logs"
      private_dns_enabled = true
      subnet_ids          = module.vpc.private_subnets
      tags                = { Name = "${local.environment}-logs" }
    }
    {{ end -}}
    {{ if dig "VpcEndpoints" "Ssm" false . }}
    # For SSM parameter store (secrets used by ECS containers).
    ssm = {
      service             = "ssm"
      private_dns_enabled = true
      subnet_ids          = module.vpc.private_subnets
      tags                = { Name = "${local.environment}-ssm" }
    }
    {{ end -}}
    {{ if dig "VpcEndpoints" "SsmMessages" false . }}
    # For ECS Exec to run commands in or get a shell to a ECS container running on AWS Fargate.
    ssmmessages =  {
      service             = "ssmmessages"
      private_dns_enabled = true
      subnet_ids          = module.vpc.private_subnets
      tags                = { Name = "${local.environment}-ssmmessages" }
    }
    {{ end -}}
    {{ if dig "VpcEndpoints" "Prometheus" false . }}
    # For Prometheus. See https://docs.aws.amazon.com/prometheus/latest/userguide/AMP-and-interface-VPC.html
    aps_workspaces = {
      service             = "aps-workspaces"
      private_dns_enabled = true
      subnet_ids          = module.vpc.private_subnets
      tags                = { Name = "${local.environment}-aps-workspace" }
    }
    {{ end -}}
    {{ if dig "VpcEndpoints" "Xray" false . }}
    # For X-ray. See TODO
    xray = {
      service             = "xray"
      private_dns_enabled = true
      subnet_ids          = module.vpc.private_subnets
      tags                = { Name = "${local.environment}-xray" }
    }
    {{ end -}}
    {{ if dig "VpcEndpoints" "Sqs" false . }}
    # For Amazon SQS
    sqs = {
      service             = "sqs"
      private_dns_enabled = true
      subnet_ids          = module.vpc.private_subnets
      tags                = { Name = "${local.environment}-sqs" }
    }
    {{ end -}}
    {{ if dig "VpcEndpoints" "SecretsManager" false . }}
    # For AWS Secrets Manager
    secretsmanager = {
      service             = "secretsmanager"
      private_dns_enabled = true
      subnet_ids          = module.vpc.private_subnets
      tags                = { Name = "${local.environment}-secretsmanager" }
    }
    {{ end -}}
    {{ if dig "VpcEndpoints" "Sts" false . }}
    # For AWS Security Token Service
    sts = {
      service             = "sts"
      private_dns_enabled = true
      subnet_ids          = module.vpc.private_subnets
      tags                = { Name = "${local.environment}-sts" }
    }
    {{ end -}}
    {{ if dig "VpcEndpoints" "Lambda" false . }}
    # For AWS Lambda
    lambda = {
      service             = "lambda"
      private_dns_enabled = true
      subnet_ids          = module.vpc.private_subnets
      tags                = { Name = "${local.environment}-lambda" }
    }
    {{ end -}}

    {{ if dig "VpcEndpoints" "S3" false . }}
    ################################################################################
    # Gateway endpoints
    # See https://docs.aws.amazon.com/vpc/latest/privatelink/gateway-endpoints.html
    ################################################################################

    s3 = {
      # Amazon S3 gateway endpoint for ECS tasks to access private ECR images.
      # Necessary since ECR utilizes S3 for storing image layers.
      # Important: For services using this endpoint, configure the security group to allow S3 IP traffic using the S3 prefix list.
      # See https://docs.aws.amazon.com/vpc/latest/privatelink/vpc-endpoints-s3.html
      service             = "s3"
      service_type        = "Gateway"
      private_dns_enabled = true
      subnet_ids          = module.vpc.private_subnets
      route_table_ids     = flatten([module.vpc.intra_route_table_ids, module.vpc.private_route_table_ids, module.vpc.public_route_table_ids, module.vpc.default_route_table_id])
      tags                = { Name = "${local.environment}-s3" }
    }
    {{ end -}}
  }
}

# Recommended reading
#
# - https://docs.aws.amazon.com/AmazonECR/latest/userguide/vpc-endpoints.html
# - https://docs.aws.amazon.com/vpc/latest/privatelink/aws-services-privatelink-support.html
# - https://aws.amazon.com/blogs/storage/introducing-private-dns-support-for-amazon-s3-with-aws-privatelink/
# - https://github.com/terraform-aws-modules/terraform-aws-vpc/issues/982
{{ else }}
# x-boilerplate-delete
{{ end -}}
