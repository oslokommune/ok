{{- if and (hasKey .IamForCicd "AssumableCdRole") .IamForCicd.AssumableCdRole }}
{{ template "doNotEdit" . }}
module "iam_assumable_role_cd_debug" {
  # https://github.com/terraform-aws-modules/terraform-aws-iam/tree/v5.5.5/modules/iam-assumable-role
  source  = "terraform-aws-modules/iam/aws//modules/iam-assumable-role"
  version = "5.52.2"

  trusted_role_arns = [
    "{{ shell "aws" "iam" "list-roles" "--query" "Roles[?starts_with(RoleName,'AWSReservedSSO_AWSAdministratorAccess_')].Arn" "--output" "text" "--no-cli-pager" | trim }}"
  ]

  create_role = true

  role_name  = local.assume_cd_role_name
  role_description = "Role that can be assumed by SSO users to debug IAM permission issues with continuous deployment"

  role_requires_mfa = false

  custom_role_policy_arns = local.role_policy_arns_cd
}
{{ else -}}
# x-boilerplate-delete
{{ end -}}
