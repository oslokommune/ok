{{ template "doNotEdit" . }}
terraform {
  required_version = "{{ .TerraformVersion }}"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "{{ .AwsProviderVersion }}"
    }
  }

{{- if dig "S3Backend" false . }}

  backend "s3" {
    region = "{{ .Region }}"

    # WARNING: Do not change these values - it can cause a lot of problems
    bucket  = "ok-iac-config-{{ .AccountId }}-{{ .Region }}-{{ .Environment }}"
    key     = "terraform/{{ .StackName }}/terraform.tfstate"
    encrypt = true

    dynamodb_table = "terraform-state-lock-{{ .Environment }}"
  }

{{- end }}
}

provider "aws" {
  region              = local.region
  allowed_account_ids = [local.account_id]

{{- if and (hasKey .IamForCicd "AssumableCdRole") .IamForCicd.AssumableCdRole }}

  dynamic "assume_role" {
    for_each = var.assume_cd_role ? [1] : []

    content {
      role_arn = "arn:aws:iam::${local.account_id}:role/${local.assume_cd_role_name}"
    }
  }

{{- end }}

  default_tags {
    tags = local.common_tags
  }
}
