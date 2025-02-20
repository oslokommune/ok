{{ if .IamForCicd.Enable -}}
#!/usr/bin/env bash

# -----------------------------------------------------------------------------
# Common stuff
# -----------------------------------------------------------------------------

# Check if user is logged in with gh
if ! gh auth status > /dev/null 2>&1; then
  read -p "You are not logged in with gh. Do you want to run 'gh auth login' now? (y/n): " response
  if [[ "$response" == "y" ]]; then
    gh auth login
  else
    echo -e "\e[31mError: You are not logged in with gh. Please run 'gh auth login' to log in.\e[0m"
    exit 1
  fi
fi

cd ..
GITHUB_ORG=$(terraform output --raw "github_org")
if [[ -z "$GITHUB_ORG"  ]]; then
  echo -e "\e[31mError: One or more required variables are empty:\e[0m"
  echo "GITHUB_ORG: $GITHUB_ORG"
  exit 1
fi

# -----------------------------------------------------------------------------
# Fetch data
# -----------------------------------------------------------------------------
echo "Fetching necessary Terraform outputs for application repository"

REPO_NAME=$(terraform output --raw "app_gh_repo_name")
REPO_ENV_NAME=$(terraform output --raw "app_gh_env_name")
AWS_ROLE_ARN=$(terraform output --raw "iam_assumable_role_github_oidc_ecr_arn")

if [[ -z "$REPO_NAME" || -z "$REPO_ENV_NAME" ]]; then
  echo -e "\e[31mError: One or more required variables are empty:\e[0m"
  echo "REPO_NAME: $REPO_NAME"
  echo "REPO_ENV_NAME: $REPO_ENV_NAME"
  exit 1
fi

if [[ "$AWS_ROLE_ARN" != arn:* ]]; then
  echo -e "\e[31mError: AWS_ROLE_ARN does not start with 'arn:'.\e[0m"
  exit 1
fi

# -----------------------------------------------------------------------------
# Set AWS_ROLE_ARN for repository
# -----------------------------------------------------------------------------
echo

echo "Command to run:"
echo -e "\e[32mecho $AWS_ROLE_ARN | gh secret set --repo $GITHUB_ORG/$REPO_NAME --env $REPO_ENV_NAME AWS_ROLE_ARN\e[0m"
printf "\nGitHub environment '$REPO_ENV_NAME' in '$GITHUB_ORG/$REPO_NAME' will be created if it is not available\n\n"

# Ask for user confirmation
read -p "Are you sure you want to run the command above? (y/n): " confirm

if [[ $confirm == [yY] ]]; then
  gh api --silent --method PUT -H "Accept: application/vnd.github+json" repos/$GITHUB_ORG/$REPO_NAME/environments/$REPO_ENV_NAME
  echo "$AWS_ROLE_ARN" | gh secret set --repo "$GITHUB_ORG/$REPO_NAME" --env "$REPO_ENV_NAME" AWS_ROLE_ARN
  echo "AWS_ROLE_ARN secret has been set."
else
  echo "Operation cancelled."
fi

echo

# -----------------------------------------------------------------------------
# Done
# -----------------------------------------------------------------------------

echo "Done!"
{{ else -}}
# x-boilerplate-delete
{{ end -}}
