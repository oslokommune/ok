/*
 * Terraform stores the state of its operations in a S3 bucket with state lock
 * stored in a dynamoDB table.
 *
 * This bootstraps the required backend infrastructure for terraform, and must
 * be run before any other operation.
 *
 * See documentation on https://www.terraform.io/language/settings/backends/s3
 */

locals {
  # this leaves 25 characters for your environment before it is truncated
  bucket_name = substr("ok-iac-config-${var.account_id}-${var.region}-${var.environment}", 0, 63)

}

/* S3 bucket for storing Terraform state */
module "s3_bucket" {
  source = "git@github.com:oslokommune/golden-path-iac//terraform/modules/s3_bucket?ref=s3_bucket-v0.1.0"

  bucket_name = local.bucket_name
  # Highly recommended in the case of accidental deletions and human error.
  # Without versioning you can wipe out your entire terraform state.
  versioning_enabled = true

  tags = local.common_tags
}

/* DynamoDB lock table for Terraform state */
#tfsec:ignore:aws-dynamodb-enable-at-rest-encryption tfsec:ignore:aws-dynamodb-table-customer-key
resource "aws_dynamodb_table" "terraform_state_lock" {
  #checkov:skip=CKV2_AWS_16:State lock does not need auto scaling
  #checkov:skip=CKV_AWS_119:State lock does not need customer managed KMS for encryption
  name = "terraform-state-lock-${var.environment}"

  hash_key       = "LockID"
  read_capacity  = 1
  write_capacity = 1

  attribute {
    name = "LockID"
    type = "S"
  }

  point_in_time_recovery {
    enabled = true
  }

  tags = local.common_tags
}
