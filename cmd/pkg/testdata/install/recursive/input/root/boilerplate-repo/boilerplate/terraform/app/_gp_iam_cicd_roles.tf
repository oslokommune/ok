{{ if dig "IamForCicd" "Enable" false . -}}
{{ template "doNotEdit" . }}
locals {

  app_gh_env_name = "${local.environment}-app-${local.main_container_name}-ecr"
  iac_gh_env_name = "${local.environment}-app-${local.main_container_name}-cicd"

  remote_state_bucket_id          = "ok-iac-config-${local.account_id}-${local.region}-${local.environment}"
  remote_state_dynamodb_table_arn = "arn:aws:dynamodb:${local.region}:${local.account_id}:table/terraform-state-lock-${local.environment}"

  oidc_client_ids    = ["sts.amazonaws.com"]
  oidc_issuer_domain = "token.actions.githubusercontent.com"
  github_org         = "oslokommune"

  ecr_repository_names = [
    data.aws_ecr_repository.app.name
  ]

  # These policies are added to all roles
  policies_common = [
    aws_iam_policy.terraform_remote_state_read_write.arn,
  ]

  role_policy_arns_cd = concat(
    [
      aws_iam_policy.cd_part1.arn,
      aws_iam_policy.cd_part2.arn,
    ],
    local.policies_common,
    local.additional_policies_for_github_actions_in_iac_repo
  )

}

# Role for pushing to ECR
module "iam_assumable_role_github_oidc_ecr" {
  # https://github.com/terraform-aws-modules/terraform-aws-iam/tree/v5.5.5/modules/iam-assumable-role-with-oidc
  source  = "terraform-aws-modules/iam/aws//modules/iam-assumable-role-with-oidc"
  version = "5.52.2"

  create_role = true

  role_path        = "/${local.environment}/${local.main_container_name}/"
  role_name_prefix = "gh_ecr_${local.main_container_name}-"
  role_description = "Used by GitHub Actions. Repository: ${local.app_gh_repo_name}. Environment: ${local.app_gh_env_name}"

  provider_url = local.oidc_issuer_domain

  oidc_fully_qualified_audiences = local.oidc_client_ids

  # Trust condition like repo:octo-org/octo-repo:environment:Production
  # See https://docs.github.com/en/actions/deployment/security-hardening-your-deployments/about-security-hardening-with-openid-connect#filtering-for-a-specific-environment
  oidc_fully_qualified_subjects = [
    "repo:${local.github_org}/${local.app_gh_repo_name}:environment:${local.app_gh_env_name}"
  ]

  force_detach_policies = true

  role_policy_arns = concat(
    local.policies_common,
    [aws_iam_policy.ecr_read_write["${local.environment}-${local.main_container_name}"].arn],
    local.additional_policies_for_github_actions_in_app_repo
  )

}

# Role for running `terraform apply`
module "iam_assumable_role_github_oidc_cicd" {
  # https://github.com/terraform-aws-modules/terraform-aws-iam/tree/v5.5.5/modules/iam-assumable-role-with-oidc
  source  = "terraform-aws-modules/iam/aws//modules/iam-assumable-role-with-oidc"
  version = "5.52.2"

  create_role = true

  role_path        = "/${local.environment}/${local.main_container_name}/"
  role_name_prefix = "gh_cicd_${local.main_container_name}-"
  role_description = "Used by GitHub Actions. Repository: ${local.app_gh_repo_name}. Environment: ${local.app_gh_env_name}"

  provider_url = local.oidc_issuer_domain

  oidc_fully_qualified_audiences = local.oidc_client_ids

  # Trust condition like repo:octo-org/octo-repo:environment:Production
  # See https://docs.github.com/en/actions/deployment/security-hardening-your-deployments/about-security-hardening-with-openid-connect#filtering-for-a-specific-environment
  oidc_fully_qualified_subjects = [
    "repo:${local.github_org}/${local.iac_gh_repo_name}:environment:${local.iac_gh_env_name}"
  ]

  force_detach_policies = true

  role_policy_arns = local.role_policy_arns_cd

}

output "iam_assumable_role_github_oidc_ecr_arn" {
  value = module.iam_assumable_role_github_oidc_ecr.iam_role_arn
}

output "iam_assumable_role_github_oidc_cicd_arn" {
  value = module.iam_assumable_role_github_oidc_cicd.iam_role_arn
}


output "github_org" {
  value = local.github_org
}

output "app_gh_repo_name" {
  value = local.app_gh_repo_name
}

output "app_gh_env_name" {
  value = local.app_gh_env_name
}

output "iac_gh_env_name" {
  value = local.iac_gh_env_name
}

output "iac_gh_repo_name" {
  value = local.iac_gh_repo_name
}

{{- else }}
# x-boilerplate-delete
{{- end }}
