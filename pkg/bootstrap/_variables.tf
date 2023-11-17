variable "team_name" {
  description = "The name of the team that owns the infrastructure"
  type        = string
}

variable "environment" {
  description = "The name of the environment (e.g. dev, staging, prod)"
  type        = string
}

variable "region" {
  description = "The AWS region"
  type        = string
}

variable "account_id" {
  description = "The AWS account ID"
  type        = string
}
