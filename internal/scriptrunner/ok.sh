#!/usr/bin/env bash

VERSION="v2.3.1" # x-release-please-version
ENV_FILE=env.yaml
GOLDEN_PATH_IAC_REPO="oslokommune/golden-path-iac"
REUSABLE_WORKFLOWS_REPO="oslokommune/reusable-workflows"
PORT_FORWARD_SCRIPT="port-forward"
BINARIES_TAG_PREFIX="ok_binaries"
COLOR_CMD="\033[0;32m" # Green. Colors to use for commands.
COLOR_ERR="\033[0;31m" # Red. Colors to use for errors.
COLOR_DEFAULT="\033[0m" # Default color. Use for regular text.
TEMPLATES_RELEASE="templates-latest"

# Path to the directory where this script is located
SCRIPTPATH="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 || exit 1; pwd -P )"

function export_environment_variables() {
  require_env_file
  export TEAM_NAME=`yq .metadata.team $ENV_FILE`
  export ENVIRONMENT=`yq .metadata.environment $ENV_FILE`
  export AWS_ACCOUNT_ID=`yq .aws.accountID $ENV_FILE`
  export AWS_REGION=`yq .aws.region $ENV_FILE`
}

function require_env_file() {
  if [[ ! -f $ENV_FILE ]]; then
      echo Cannot find environment file $ENV_FILE
      exit 1
  fi
}

function require_dependencies() {
    if ! command -v fzf &> /dev/null
    then
        echo "fzf is not installed. Please install it on your system."
        echo "For more information, please visit: https://km.oslo.systems/fzf.html"
        exit 1
    fi

    if ! command -v yq &> /dev/null
    then
        echo "yq is not installed. Please install it on your system."
        echo "For more information, please visit: https://km.oslo.systems/yq.html"
        exit 1
    fi

    if ! command -v terraform &> /dev/null
    then
        echo "Terraform is not installed. Please install it on your system."
        echo "For more information, please visit: https://km.oslo.systems/terraform.html"
        exit 1
    fi

    if ! command -v jq &> /dev/null
    then
        echo "jq is not installed. Please install it on your system."
        echo "For more information, please visit: https://stedolan.github.io/jq/download/"
        exit 1
    fi

    if ! command -v gh &> /dev/null
    then
        echo "gh is not installed. Please install it on your system."
        echo "For more information, please visit: https://km.oslo.systems/gh.html"
        exit 1
    fi
}

# Prints a message in red
# $1: message
function print_err() {
  echo -e "${COLOR_ERR}${1}${COLOR_DEFAULT}"
}

# Tests that the user is logged in to GitHub.
function validate_golden_path_iac_repo_access() {
  local IS_LOGGED_IN

  gh auth status >/dev/null 2>&1
  IS_LOGGED_IN=$?

  if [[ ! $IS_LOGGED_IN -eq 0 ]]; then
      print_err "You are not logged in to GitHub."
      echo
      echo "Possible solutions:"

      if [[ -z $GITHUB_TOKEN ]]; then
        echo -e "- The environment variable GITHUB_TOKEN is set. Is it valid? Consider unsetting it, and use ${COLOR_CMD}gh auth login${COLOR_DEFAULT} instead."
      else
        echo -e "- Use ${COLOR_CMD}gh auth login${COLOR_DEFAULT} to log in."
      fi

      echo
      echo -e "Details from running ${COLOR_CMD}gh auth status${COLOR_DEFAULT}:"
      echo
      gh auth status
      exit 1
  fi
}

function getVersions() {
  local status
  local tags

  tags=$(\
    gh api "repos/$GOLDEN_PATH_IAC_REPO/releases" \
      | jq -r '.[] | select(.tag_name | startswith("'$BINARIES_TAG_PREFIX'") ) | [.tag_name, .published_at] | @tsv' | sort -rV
  )

  status=$?
  if [[ $status -ne 0 ]]; then
    echo -e "Internal error: Command failed: ${COLOR_CMD}gh api repos/$GOLDEN_PATH_IAC_REPO/releases${COLOR_DEFAULT}"
    echo
    echo "This shouldn't happen, and can be filed as a bug at https://github.com/$GOLDEN_PATH_IAC_REPO"
    exit 1
  fi

  echo "$tags"
}

function downloadVersion() {
  local tagToDownload
  tagToDownload=$1
  local tmpDir
  tmpDir=$(mktemp -d)
  local targetFile
  targetFile="$tmpDir/release.tar.gz"
  local outputPath
  outputPath=$SCRIPTPATH

  read -r -p "Download to a different directory than $outputPath? [y/N] " response
  if [[ "$response" =~ ^([yY][eE][sS]|[yY])+$ ]]; then
    read -r -p "Enter path: " outputPath

    if [[ ! -d $outputPath ]]; then
      print_err "Directory $outputPath does not exist, create it and run again."
      exit 1
    fi

    echo "Downloading to $outputPath"
  fi

  cd "$tmpDir" || (echo "Could not change directory to $tmpDir" && exit 1)
  local status
  status=$(gh release download "$tagToDownload" \
    --repo $GOLDEN_PATH_IAC_REPO \
    --archive=tar.gz \
    --output "$targetFile"
    )

  if [[ $status -ne 0 ]]; then
    print_err "Could not download version $tagToDownload"
    echo
    echo "Internal error: $status"
    exit 1
  fi

  tar xfz \
        "$targetFile" \
        --strip-components=1

  echo "Do you want to view the changelog? (y/n)"
  read -r answer
  if [[ "$answer" == "y" ]]; then
    less "$tmpDir/bin/CHANGELOG.md"
  fi

  doUpdate=false

  while true; do
    read -rp "Continue with installation? (y/n) " yn
      case $yn in
        [Yy]* ) doUpdate=true; break;;
        [Nn]* ) unset doUpdate; break;;
        * ) break;;
      esac
  done

  if [[ "$doUpdate" == true ]]; then
    echo "Updating files..."
    # Example files being copied:
    # ok
    # port-forward

    local unameOut
    unameOut="$(uname -s)"
    case "${unameOut}" in
        Linux*)
          find bin -maxdepth 1 -executable -type f -exec cp {} "$outputPath" \;
          ;;
        Darwin*)
          find bin -maxdepth 1 -perm +111 -type f -exec cp {} "$outputPath" \;
          ;;
        *)
          print_err "Your operation system is not supported: $unameOut"
          exit 1
          ;;
    esac
  else
    echo "Aborting"
    exit 0
  fi

  if [[ ! -f "$outputPath/ok" ]]; then
    print_err "Please try the update again or manually update the files by copying them from $tmpDir to $outputPath"
    exit 1
  fi

  if [[ ! -f "$outputPath/port-forward" ]]; then
    print_err "Please try the update again or manually update the files by copying them from $tmpDir to $outputPath"
    exit 1
  fi

  if [[ ! -x "$outputPath/ok" ]]; then
    chmod +x "$outputPath/ok"
  fi
  if [[ ! -x "$outputPath/port-forward" ]]; then
    chmod +x "$outputPath/port-forward"
  fi

  echo "Update complete"
}

function tagToSemver() {
    local tag
    tag=$1
    local semver
    semver=$(echo "$tag" | sed -e "s/$BINARIES_TAG_PREFIX-//g")
    echo "$semver"
}

function update() {
  validate_golden_path_iac_repo_access

  local versions
  local status

  versions=$(getVersions)
  status=$?

  if [[ $status -ne 0 ]]; then
    print_err "Could not get versions"
    echo
    echo "$versions"
    exit 1
  fi

  if [[ -z "$versions" ]]; then
    print_err "No versions found"
    exit 1
  fi

  local version
  version=$(echo "$versions" | fzf --header="Select version" | awk '{print $1}')

  if [[ -z "$version" ]]; then
    print_err "No version selected"
    exit 1
  fi

  downloadVersion "$version"
}

function checkVersion() {
  echo "Current version: $VERSION"

  local IS_LOGGED_IN
  gh auth status >/dev/null 2>&1
  IS_LOGGED_IN=$?
  if [[ ! $IS_LOGGED_IN -eq 0 ]]; then
    exit 0
  fi

  local latestVersion
  latestVersion=$(getVersions | head -n1 | awk '{print $1}')
  local currentVersion
  currentVersion=$(tagToSemver "$latestVersion")

  if [[ "$currentVersion" == "$VERSION" ]]; then
    echo "You are running the latest version of ok"
  else
    echo "Latest version: $currentVersion"
    echo "Run 'ok update' to update or change version"
  fi
}

function printHelp(){
  if [[ $1 == "bootstrap" ]]; then
    echo Usage: ok bootstrap
    echo

    echo "This command will create the necessary S3 bucket and DynamoDB table that will be used to store Terraform state."
    echo "For more information review: https://km.oslo.systems/setup/infrastructure/bootstrap-the-environment/"
  elif [[ $1 == "scaffold" ]]; then
    echo Usage: ok scaffold
    echo

    echo "Creates a new Terraform project with a _config.tf, _variables.tf, _versions.tf and _config.auto.tfvars.json file based on values configured in env.yml."
  elif [[ $1 == "env" ]]; then
    echo Usage: ok env
    echo

    echo "Creates a new env.yml file with placeholder values."
  elif [[ $1 == "envars" ]]; then
    echo Usage: ok envars
    echo

    echo "Exports the values in env.yml as environment variables."
  elif [[ $1 == "get-template" ]]; then
    echo Usage: ok get-template [template-name]
    echo

    echo "Downloads a template from the golden-path-iac repository."
    echo "https://github.com/oslokommune/golden-path-iac/tree/main/terraform/templates"
  elif [[ $1 == "forward" ]]; then
    echo Usage: ok forward
    echo

    echo "Starts a port forwarding session to a database. See "Connect to a database from your computer"."
  elif [[ $1 == "update" ]]; then
    echo Usage: ok update
    echo

    echo "List all available remote versions of ok and provide a menu to select which version to use."
  elif [[ $1 == "version" ]]; then
    echo Usage: ok version
    echo

    echo "Prints the version of the ok tool and the current latest version available."
  else
    echo Usage: ok [command]
    echo

    echo ok is a command line tool used in the Golden Path IaC project. It is used to bootstrap and scaffold new projects.
    echo

    echo Available commands:
    echo "  bootstrap"
    echo "  scaffold"
    echo "  env"
    echo "  envars"
    echo "  get-template"
    echo "  forward"
    echo "  update"
    echo "  version"
  fi
  exit 0
}

# Downloads a template from the $GOLDEN_PATH_IAC_REPO repo from the specified release.
# The second argument ($2) specifies the path to save the file, if omitted the template will be returned on stdout
#
# For debugging purposes you can fetch templates from disk by running:
#   `$ TEMPLATE_FROM_DISK=true ok bootstrap`
#
# $1 - The name of the template to download
# $2 - The name of the file to write the template to
function get_template() {
  validate_golden_path_iac_repo_access

  if [[ -z "$1" ]]; then
    print_err "First parameter must be a valid string name"
    exit 1
  fi

  local TEMPLATE_NAME
  # If not file extension is provided, assume .tf
  if [[ "$1" == *.* ]]; then
    TEMPLATE_NAME="$1"
  else
    TEMPLATE_NAME="$1.tf"
  fi

  local OUTPUT_FILE=$2

  # Get the template from disk
  if [[ -n "$TEMPLATE_FROM_DISK" ]]; then
    SCRIPTPATH="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
    TEMPLATE_ON_DISK="$SCRIPTPATH/../terraform/templates/$TEMPLATE_NAME"

    if [[ ! -f $TEMPLATE_ON_DISK ]]; then
      print_err "Template not found: $TEMPLATE_ON_DISK"
      exit 1
    fi

    if [[ -n "$OUTPUT_FILE" ]]; then
      cp "$TEMPLATE_ON_DISK" "$OUTPUT_FILE"
      return
    fi

    # Since output file is empty, we print to stdout
    cat "$TEMPLATE_ON_DISK"

    return
  fi

  # Get the template from web
  if [[ -z "$OUTPUT_FILE" ]]; then
    # Return the template on stdout
    gh release download "${TEMPLATES_RELEASE}" \
      --repo $GOLDEN_PATH_IAC_REPO \
      --pattern "$TEMPLATE_NAME" \
      --output -
  else
    gh release download "${TEMPLATES_RELEASE}" \
      --repo $GOLDEN_PATH_IAC_REPO \
      --pattern "$TEMPLATE_NAME" \
      --output "$OUTPUT_FILE"
  fi

  SUCCESS=$?
  if [[ $SUCCESS -ne 0 ]]; then
    exit 1
  fi
}

require_dependencies

for i in "$@"
do
case $i in
    --help)
    printHelp $1
    ;;
esac
done

if [[ $1 == "version" ]]; then
    checkVersion
    exit 0
fi

if [[ $1 == "workflows" ]]; then
    WORKFLOW_NAME=$(gh release view "${TEMPLATES_RELEASE}" --repo "${REUSABLE_WORKFLOWS_REPO}" --json assets --jq '.assets[].name' | fzf)

    gh release download "${TEMPLATES_RELEASE}" \
        --repo "${REUSABLE_WORKFLOWS_REPO}" \
        --pattern "${WORKFLOW_NAME}"
    exit 0
fi

if [[ $1 == "bootstrap" || $1 == "scaffold" ]]; then

    require_env_file

    yq $ENV_FILE &> /dev/null
    exit_status=$?
    if [ $exit_status -eq 1 ]; then
      echo Error: $ENV_FILE is not a valid yaml file
      echo
      echo Execute "'yq $ENV_FILE'" to see the error, then try again
      exit
    fi

    validate_golden_path_iac_repo_access
fi


if [[ $1 == "bootstrap" ]]; then
    DIR="remote_state"

    echo Creating directory $DIR
    mkdir -p $DIR

    export TEAM_NAME=`yq .metadata.team $ENV_FILE`
    export ENVIRONMENT=`yq .metadata.environment $ENV_FILE`
    export AWS_REGION=`yq .aws.region $ENV_FILE`
    export AWS_ACCOUNT_ID=`yq .aws.accountID $ENV_FILE`

    if [[ -f $DIR/_config.auto.tfvars.json ]]; then
        echo $DIR/_config.auto.tfvars.json already exists, not overwriting
    else
        echo Creating local variable file $DIR/_config.auto.tfvars.json
        TEMPLATE=$(get_template "_config.auto.tfvars.json.template")
        if [[ $? -ne 0 ]]; then
          print_err "$TEMPLATE"
          exit 1
        fi
        echo "$TEMPLATE" | envsubst '$TEAM_NAME,$ENVIRONMENT,$AWS_REGION,$AWS_ACCOUNT_ID' > $DIR/_config.auto.tfvars.json
    fi
    echo Copying common configuration file to $DIR/_config.auto.tfvars.json

    if [[ -f $DIR/_config.tf ]]; then
        echo $DIR/_config.tf already exists, not overwriting
    else
        echo Creating local variable file $DIR/_config.tf
        TEMPLATE=$(get_template "_config.tf.template")
        if [[ $? -ne 0 ]]; then
          print_err "$TEMPLATE"
          exit 1
        fi
        echo "$TEMPLATE" | envsubst '$TEAM_NAME,$ENVIRONMENT' > $DIR/_config.tf
    fi
    echo Copying common configuration file to $DIR/_config.tf

    if [[ -f $DIR/_variables.tf ]]; then
        echo $DIR/_variables.tf already exists, not overwriting
    else
        echo Creating local variable file $DIR/_variables.tf
        TEMPLATE=$(get_template "_variables.tf.template")
        if [[ $? -ne 0 ]]; then
          print_err "$TEMPLATE"
          exit 1
        fi
        echo "$TEMPLATE" | envsubst '$TEAM_NAME,$ENVIRONMENT,$AWS_REGION,$AWS_ACCOUNT_ID' > $DIR/_variables.tf
    fi
    echo Copying common configuration file to $DIR/_variables.tf.json

    if [[ -f $DIR/_versions.tf ]]; then
        echo $DIR/_versions_bootstrap.tf already exists, not overwriting
    else
        echo Creating local variable file $DIR/_versions.tf
        TEMPLATE=$(get_template "_versions_bootstrap.tf.template")
        if [[ $? -ne 0 ]]; then
          print_err "$TEMPLATE"
          exit 1
        fi
        echo "$TEMPLATE" | envsubst '$STATE_NAME,$ENVIRONMENT,$AWS_REGION,$AWS_ACCOUNT_ID' > $DIR/_versions_bootstrap.tf
    fi
    echo Copying common configuration file to $DIR/_versions_bootstrap.tf.json

    echo Copying S3 backend template to $DIR/terraform_s3_backend.tf
    get_template "terraform_s3_backend.tf" "$DIR/terraform_s3_backend.tf"

elif [[ $1 == "scaffold" ]]; then
    DIR=$2

    if [[ -z "$DIR" ]]; then
        echo "Missing directory name"
        echo
        echo "USAGE: ok scaffold <directory>"
        exit 1
    fi

    echo Creating directory $DIR
    mkdir -p $DIR

    export STATE_NAME=$DIR
    export_environment_variables

    if [[ -n "$DELETE_OLD_STYLE_CONFIG" ]]; then
        rm -f $DIR/config.tf
        rm -f $DIR/common.tf
    fi

    if [[ -f $DIR/_config.auto.tfvars.json ]]; then
        echo $DIR/_config.auto.tfvars.json already exists, not overwriting
    else
        echo Creating local variable file $DIR/_config.auto.tfvars.json
        TEMPLATE=$(get_template "_config.auto.tfvars.json.template")
        if [[ $? -ne 0 ]]; then
          print_err "$TEMPLATE"
          exit 1
        fi
        echo "$TEMPLATE" | envsubst '$TEAM_NAME,$ENVIRONMENT,$AWS_REGION,$AWS_ACCOUNT_ID' > $DIR/_config.auto.tfvars.json
    fi
    echo Copying common configuration file to $DIR/_config.auto.tfvars.json

    if [[ -f $DIR/_config.tf ]]; then
        echo $DIR/_config.tf already exists, not overwriting
    else
        echo Creating local variable file $DIR/_config.tf
        TEMPLATE=$(get_template "_config.tf.template")
        if [[ $? -ne 0 ]]; then
          print_err "$TEMPLATE"
          exit 1
        fi
        echo "$TEMPLATE" | envsubst '$TEAM_NAME,$ENVIRONMENT' > $DIR/_config.tf
    fi
    echo Copying common configuration file to $DIR/_config.tf

    if [[ -f $DIR/_variables.tf ]]; then
        echo $DIR/_variables.tf already exists, not overwriting
    else
        echo Creating local variable file $DIR/_variables.tf
        TEMPLATE=$(get_template "_variables.tf.template")
        if [[ $? -ne 0 ]]; then
          print_err "$TEMPLATE"
          exit 1
        fi
        echo "$TEMPLATE" | envsubst '$TEAM_NAME,$ENVIRONMENT,$AWS_REGION,$AWS_ACCOUNT_ID' > $DIR/_variables.tf
    fi
    echo Copying common configuration file to $DIR/_variables.tf.json

    if [[ -f $DIR/_versions.tf ]]; then
        echo $DIR/_versions.tf already exists, not overwriting
    else
        echo Creating local variable file $DIR/_versions.tf
        TEMPLATE=$(get_template "_versions.tf.template")
        if [[ $? -ne 0 ]]; then
          print_err "$TEMPLATE"
          exit 1
        fi
        echo "$TEMPLATE" | envsubst '$STATE_NAME,$ENVIRONMENT,$AWS_REGION,$AWS_ACCOUNT_ID' > $DIR/_versions.tf
    fi
    echo Copying common configuration file to $DIR/_versions.tf.json

elif [[ $1 == "env" ]]; then
    if [[ -f env.yaml ]]; then
        echo This would overwrite the file existing file env.yaml. Move or delete it and retry.
        exit 1
    fi

    cat << EOF > ./env.yaml
apiVersion: 1.0
kind: Environment

metadata:
  team: 'my-project'
  environment: 'my-project-dev'

aws:
  accountID: '123456789012'
  region: 'eu-west-1'
EOF
elif [[ $1 == "envars" ]]; then
  export_environment_variables
elif [[ $1 == "get-template" ]]; then
  TEMPLATE_NAME="$2"
  OUTPUT_FILE="$3"

  if [[ -z $TEMPLATE_NAME ]]; then
    validate_golden_path_iac_repo_access

    TEMPLATE_NAME=$(gh release view $TEMPLATES_RELEASE \
      --repo "$GOLDEN_PATH_IAC_REPO" \
      --json assets --jq '.assets[].name' \
      | fzf)

    if [[ -z $TEMPLATE_NAME ]]; then
      echo "No template selected, exiting."
      exit
    fi
  elif [[ ! "$TEMPLATE_NAME" =~ \. ]]; then
    TEMPLATE_NAME="$TEMPLATE_NAME.tf"
  fi

  if [[ -z "$OUTPUT_FILE" ]]; then
      OUTPUT_FILE="$TEMPLATE_NAME"
  fi

  get_template "$TEMPLATE_NAME" "$OUTPUT_FILE"

  if [[ -z "$2" ]]; then
    echo "Downloaded template $TEMPLATE_NAME"
  fi
elif [[ $1 == "forward" ]]; then
  "$SCRIPTPATH/$PORT_FORWARD_SCRIPT" "$@"
elif [[ $1 == "update" ]]; then
  update
elif [[ $2 != "tests" ]]; then
  printHelp
fi
