{{ if dig "IamForCicd" "Enable" false . -}}
{{ template "doNotEdit" . }}
module "iam_policies_generic" {
  source = "git@github.com:oslokommune/golden-path-iac//terraform/modules/iam-policies-generic?ref=iam-policies-generic-v4.1.1"

  account_id  = local.account_id
  region      = local.region
  environment = local.environment
  app_names   = [local.main_container_name]
}

module "iam_policies_cicd" {
  source = "git@github.com:oslokommune/golden-path-iac//terraform/modules/iam-policies-cicd?ref=iam-policies-cicd-v0.2.1"

  s3_bucket_id       = local.remote_state_bucket_id
  dynamodb_table_arn = local.remote_state_dynamodb_table_arn

  ecr_repository_names = local.ecr_repository_names
}

resource "aws_iam_policy" "ecr_read_write" {

  # This is a map where the app name is the key and the value is a policy document
  for_each = module.iam_policies_cicd.policy_documents.ecr_push

  path        = "/${local.environment}/${local.main_container_name}/"
  name        = "${local.environment}_${each.key}_ecr_read_write"
  description = "Provides write access to Amazon Elastic Container Registry repository associated with the environment \"${local.environment}\" and the app \"${each.key}\""
  policy      = each.value

}

resource "aws_iam_policy" "terraform_remote_state_read_write" {
  path        = "/${local.environment}/${local.main_container_name}/"
  name        = "${local.environment}_${local.main_container_name}_remote_state_read_write"
  description = "Provides write access to the Terraform remote state bucket and DynamoDB table associated with the environment \"${local.environment}\" and the app \"${local.main_container_name}\""
  policy      = module.iam_policies_cicd.policy_documents.remote_state
}

resource "aws_iam_policy" "cd_part1" {
  path        = "/${local.environment}/${local.main_container_name}/"
  name        = "${local.environment}-${local.main_container_name}-cd-part-1"
  description = "Allows continuous deployment of the ECS service ${local.main_container_name} (part 1 of 2)"

  policy = data.aws_iam_policy_document.cd_part1.json
}

resource "aws_iam_policy" "cd_part2" {
  path        = "/${local.environment}/${local.main_container_name}/"
  name        = "${local.environment}-${local.main_container_name}-cd-part-2"
  description = "Allows continuous deployment of the ECS service ${local.main_container_name} (part 2 of 2)"

  policy = data.aws_iam_policy_document.cd_part2.json
}

data "aws_iam_policy_document" "cd_part1" {
  source_policy_documents = [
    module.iam_policies_generic.policy_documents.acm_read,
    module.iam_policies_generic.policy_documents.alb_wildcard_list,
    module.iam_policies_generic.policy_documents.alb_read,
    module.iam_policies_generic.policy_documents.autoscaling_wildcard_list,
    module.iam_policies_generic.policy_documents.autoscaling_read,
    module.iam_policies_generic.policy_documents.autoscaling_write,
    module.iam_policies_generic.policy_documents.cloud_map_service_discovery_read,
    module.iam_policies_generic.policy_documents.cloud_map_service_discovery_wildcard_list,
    module.iam_policies_generic.policy_documents.cloudwatch_dashboard_read,
    module.iam_policies_generic.policy_documents.cloudwatch_log_groups_wildcard_list,
    module.iam_policies_generic.policy_documents.cloudwatch_log_groups_read,
    module.iam_policies_generic.policy_documents.ecr_read,
    module.iam_policies_generic.policy_documents.ecs_read,
    module.iam_policies_generic.policy_documents.ecs_tagging,
    module.iam_policies_generic.policy_documents.ecs_service_read,
    module.iam_policies_generic.policy_documents.ecs_service_write,
    module.iam_policies_generic.policy_documents.ecs_task_definition_wildcard_read,
    module.iam_policies_generic.policy_documents.ecs_task_definition_wildcard_write,
  ]
}

data "aws_iam_policy_document" "cd_part2" {
  source_policy_documents = [
    module.iam_policies_generic.policy_documents.iam_pass_role,
    module.iam_policies_generic.policy_documents.iam_detatch_role,
    module.iam_policies_generic.policy_documents.iam_policy_read,
    module.iam_policies_generic.policy_documents.iam_policy_read_path_based,
    module.iam_policies_generic.policy_documents.iam_policy_create,
    module.iam_policies_generic.policy_documents.iam_policy_delete_path_based,
    module.iam_policies_generic.policy_documents.iam_role_read,
    module.iam_policies_generic.policy_documents.prometheus_read,
    module.iam_policies_generic.policy_documents.prometheus_wildcard_list,
    module.iam_policies_generic.policy_documents.rds_db_subnets_read,
    module.iam_policies_generic.policy_documents.route53_read,
    module.iam_policies_generic.policy_documents.route53_wildcard_list,
    module.iam_policies_generic.policy_documents.ssm_wildcard_list,
    module.iam_policies_generic.policy_documents.ssm_read,
    module.iam_policies_generic.policy_documents.ssm_tags_read,
    module.iam_policies_generic.policy_documents.vpc_wildcard_list,
    module.iam_policies_generic.policy_documents.vpc_read,
    module.iam_policies_generic.policy_documents.wildcard_resources_read,
  ]
}
{{- else }}
# x-boilerplate-delete
{{- end }}
