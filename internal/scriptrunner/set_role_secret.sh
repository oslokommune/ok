#!/usr/bin/env bash
OK_AWS_ROLE_ARN=$(terraform output -raw 'iam_assumable_role_github_oidc_ecs_deploy_arn')
GITHUB_ORG=$(jq -r '.github_org' _config.auto.tfvars.json)
IAC_REPO_NAME=$(jq -r '.iac_repo_name' _config.auto.tfvars.json)
IAC_REPO="$GITHUB_ORG/$IAC_REPO_NAME"
IAC_ENV=$(jq -r '.main_container_name' _config.auto.tfvars.json)

# You must create the GitHub environment for this to work
echo "$OK_AWS_ROLE_ARN" | gh secret set --repo "$IAC_REPO" --env "$IAC_ENV" AWS_ROLE_ARN
