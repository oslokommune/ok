#!/bin/bash
# shellcheck disable=SC2034
VERSION="v2.3.1" # x-release-please-version

configFile="$HOME/.ok-port-forward.conf"

if [[ "$1" == "-h" || "$1" == "--help" ]]; then
    echo "Usage: $0 [options]"
    echo "Options:"
    echo "  -h, --help            Show this help message and exit"
    echo "  -k, --keep            Saves the selection you make for the next time by storing it in $configFile"
    echo "  -c, --clear           Clears the saved selections, e.g. it deletes $configFile"
    exit 0
fi

if [[ "$1" == "-k" || "$1" == "--keep" ]]; then
    WRITE_CONFIG=true
fi

colorRed="\033[31m";
colorGreen="\033[32m";
colorYellow="\033[33m";
colorReset="\033[0m";

function printMessage {
    _message=$1
    _color=$2
    _date=$(date +"%H:%M:%S")
    if [ "$_color" == "red" ]; then
        echo -e "$_date: $colorRed$_message$colorReset"
    elif [ "$_color" == "yellow" ]; then
        echo -e "$_date: $colorYellow$_message$colorReset"
    elif [ "$_color" == "green" ]; then
        echo -e "$_date: $colorGreen$_message$colorReset"
    else
        echo -e "$_date: $_message"
    fi
}

function checkRequirements() {
    if ! command -v aws &> /dev/null
    then
        printMessage "aws-cli is not installed. Please install it on your system." "red"
        printMessage "For more information, please visit: https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html" "red"
        exit 1
    fi
    if ! command -v jq &> /dev/null
    then
        printMessage "jq is not installed. Please install it on your system." "red"
        printMessage "For more information, please visit: https://stedolan.github.io/jq/download/" "red"
        exit 1
    fi
    # session manager plugin
    if ! command -v session-manager-plugin &> /dev/null
    then
        printMessage "session-manager-plugin is not installed. Please install it on your system." "red"
        printMessage "For more information, please visit: https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html" "red"
        exit 1
    fi
    if ! command -v fzf &> /dev/null
    then
        printMessage "fzf is not installed. Please install it on your system." "red"
        printMessage "For more information, please visit: https://github.com/junegunn/fzf#installation" "red"
        exit 1
    fi
}

checkRequirements

tagPrefix="OkTool:"
taskDefinitionFamilyPrefix="ssm-portfwd"

if [ -f "$configFile" ]; then
    # shellcheck source=/dev/null
    source "$configFile"
else
    if [ "$WRITE_CONFIG" = true ]; then
        touch "$configFile"
    fi
fi

if [ -n "$TASK_DEF_FAMILY" ]; then
    while true; do
        read -rp "$(echo -e "Using task def:" "$TASK_DEF_FAMILY"? "(n to reset)?") " yn
        case $yn in
            [Yy]* ) break;;
            [Nn]* ) unset TASK_DEF_FAMILY; break;;
            * ) break;;
        esac
    done
fi

if [ -z "$TASK_DEF_FAMILY" ]; then
    printMessage "Loading task definition families ..."
    families=$(aws ecs list-task-definition-families --status ACTIVE --family-prefix="$taskDefinitionFamilyPrefix" --query 'families[]' | jq -r '.[]')
    if [ -z "$families" ] || [ "$families" == "null" ]; then
        printMessage "No task definitions found" red
        exit 1
    fi
    TASK_DEF_FAMILY=$(echo "$families" | fzf)
    if [ -z "$TASK_DEF_FAMILY" ]; then
        printMessage "You must select a task definition family" yellow
        exit 1
    fi
fi

printMessage "Loading task definition"
TASK_DEF_ARN=$(aws ecs list-task-definitions --family-prefix "$TASK_DEF_FAMILY" --sort DESC --query 'taskDefinitionArns[0]' | jq -r '.')
if [ -z "$TASK_DEF_ARN" ] || [ "$TASK_DEF_ARN" = "null" ]; then
    printMessage "No task definition found for $TASK_DEF_FAMILY. This is probably a bug." red
    exit 1
fi

printMessage "Loading task definition details ..."
taskDefTags=$(aws ecs list-tags-for-resource --resource-arn "$TASK_DEF_ARN" | jq -r '.tags[] | select(.key | startswith("'$tagPrefix'"))')

if [ -z "$taskDefTags" ] || [ "$taskDefTags" == "null" ]; then
    printMessage "No tags found for $TASK_DEF_ARN. This is probably a bug." red
    exit 1
fi

securityGroups=$(echo "$taskDefTags" | jq -r 'select(.key == "OkTool:SecurityGroup") | .value')
subnets=$(echo "$taskDefTags" | jq -r 'select(.key == "OkTool:Subnets") | .value')
# Convert forward slashes to commas (aws tags set by terraform does allow comma separators in tag values)
subnets=$(echo "$subnets" | tr '/' ',')
clusterName=$(echo "$taskDefTags" | jq -r 'select(.key == "OkTool:ClusterName") | .value')

# Validate that we have all the required tags
if [ -z "$securityGroups" ] || [ "$securityGroups" == "null" ]; then
    printMessage "No security groups found for $TASK_DEF_ARN. This is probably a bug." red
    exit 1
fi
if [ -z "$subnets" ] || [ "$subnets" == "null" ]; then
    printMessage "No subnets found for $TASK_DEF_ARN. This is probably a bug." red
    exit 1
fi
if [ -z "$clusterName" ] || [ "$clusterName" == "null" ]; then
    printMessage "No cluster name found for $TASK_DEF_ARN. This is probably a bug." red
    exit 1
fi

runTaskConfig="awsvpcConfiguration={subnets=[$subnets],securityGroups=[$securityGroups]}"

# Run the task and get the task ARN
printMessage "Starting ECS task ..."
taskArn=$(aws ecs run-task --enable-execute-command --cluster "$clusterName" --task-definition "$TASK_DEF_ARN" --count 1 --launch-type FARGATE --network-configuration "$runTaskConfig" --query 'tasks[0].taskArn' --output text)
if [[ ! $taskArn =~ ^arn:aws:ecs:[a-z0-9-]+:[0-9]+:task/.*$ ]]; then
    printMessage "Failed to run task" red
    exit 1
fi

# Parse the task id from the task arn
taskId=$(echo "$taskArn" | cut -d "/" -f 3)
if [[ ! $taskId =~ ^[a-z0-9-]+$ ]]; then
    printMessage "Failed to parse task id from task arn. This is probably a bug." red
    exit 1
fi

# Keep checking the task status until it's running
taskStatusCounter=0
_saidTaskText=false
_saidContainerText=false
_saidAgentText=false
while true; do
    if [ $taskStatusCounter -gt 45 ]; then
        printMessage "Task status check timed out. Please run the program again." red
        exit 1
    fi

    # Overall task information
    taskDescription=$(aws ecs describe-tasks --cluster "$clusterName" --tasks "$taskId" --output json)
    taskStatus=$(echo "$taskDescription" | jq -r '.tasks[0].lastStatus')

    # The portforward task
    # Do a separate `describe-tasks` here and not parse out from the json in $taskDescription above
    # to make the code easier to read and follow: jq can query for someting that "startswith", but
    # the complexity of the query makes it less readable.
    # We have to query for a explicit container that starts with "ssm-portfwd-" because GuardDuty
    # injects a separate container in the task list
    portfwdTaskDescription=$(aws ecs describe-tasks --cluster "$clusterName" --tasks "$taskId" --query "tasks[0].containers[?starts_with(name,'ssm-portfwd-')]" --output json)
    containerStatus=$(echo "$portfwdTaskDescription" | jq -r ".[0].lastStatus")
    agentLastStatus=$(echo "$portfwdTaskDescription" | jq -r ".[0].managedAgents[0].lastStatus")

    if [ "$taskStatus" == "RUNNING" ]; then
        if [ "$containerStatus" == "RUNNING" ]; then
            if [ "$agentLastStatus" == "RUNNING" ]; then
                echo
                break
            else
                if [ "$_saidAgentText" = false ]; then
                    echo
                    printMessage "Waiting for agent to start " yellow
                    _saidAgentText=true
                fi
            fi
        else
            if [ "$_saidContainerText" = false ]; then
                echo
                printMessage "Waiting for container to start " yellow
                _saidContainerText=true
            fi
        fi
    else
        if [ "$_saidTaskText" = false ]; then
            printMessage "Waiting for task to start " yellow
            _saidTaskText=true
        fi
    fi
    echo -n "."
    sleep 2
done

printMessage "Task started successfully. Retrieving runtimeId ..."

runtimeId=$(aws ecs describe-tasks --cluster "$clusterName" --tasks "$taskId" --query "tasks[0].containers[?starts_with(name,'ssm-portfwd-')].runtimeId" --output text)

if [[ ! $runtimeId =~ ^[a-z0-9-]+$ ]]; then
    printMessage "Failed to retrieve runtimeId. This is probably a bug." red
    exit 1
fi

remotePortDefault=${REMOTE_PORT:-5432}
localPortDefault=${LOCAL_PORT:-4812}

if [ -n "$LOCAL_PORT" ] && [ -n "$REMOTE_PORT" ] && [ -n "$RDS_ENDPOINT" ]; then
    printMessage "Setting up forwarding to $colorGreen$RDS_ENDPOINT$colorReset on remote $colorGreen$REMOTE_PORT$colorReset and local $colorGreen$LOCAL_PORT$colorReset ..."
    while true; do
        read -rp "Continue? (n to re-configure)? " yn
        echo
        case $yn in
            [Yy]* ) break;;
            [Nn]* ) unset LOCAL_PORT; unset REMOTE_PORT; unset RDS_ENDPOINT; break;;
            * ) break;;
        esac
    done
fi

if [ -z "$LOCAL_PORT" ] || [ -z "$REMOTE_PORT" ] || [ -z "$RDS_ENDPOINT" ]; then

    printMessage "Loading RDS endpoints ..."
    endpoints=$(aws rds describe-db-instances | jq -r '.DBInstances[].Endpoint.Address')
    if [ -z "$endpoints" ]; then
        printMessage "No RDS endpoints found" red
        exit 1
    fi
    RDS_ENDPOINT=$(echo "$endpoints" | fzf)

    if [ -z "$RDS_ENDPOINT" ]; then
        printMessage "No RDS endpoint selected" red
        exit 1
    fi

    while true; do
        read -rp "$(echo -e "Enter local port or press enter to use:" "$colorGreen$localPortDefault$colorReset")" LOCAL_PORT
        if [ -z "$LOCAL_PORT" ]; then
            LOCAL_PORT=$localPortDefault
            break
        elif [[ "$LOCAL_PORT" =~ ^[0-9]+$ ]]; then
            break
        else
            printMessage "Port must be a number. Try again." red
        fi
    done

    while true; do
        read -rp "$(echo -e "Enter remote port or press enter to use:" "$colorGreen$remotePortDefault$colorReset")" REMOTE_PORT
        if [ -z "$REMOTE_PORT" ]; then
            REMOTE_PORT=$remotePortDefault
            break
        elif [[ "$REMOTE_PORT" =~ ^[0-9]+$ ]]; then
            break
        else
            printMessage "Port must be a number. Try again." red
        fi
    done
fi

if [ -n "$WRITE_CONFIG" ]; then
    { echo "LOCAL_PORT=$LOCAL_PORT";\
        echo "TASK_DEF_ARN=$TASK_DEF_ARN";\
        echo "REMOTE_PORT=$REMOTE_PORT";\
        echo "RDS_ENDPOINT=$RDS_ENDPOINT"
    }> "$configFile"
fi

docName="AWS-StartPortForwardingSessionToRemoteHost"
printMessage "To exit, press CTRL+C."

aws ssm start-session --target ecs:"$clusterName"_"$taskId"_"$runtimeId" --document-name "$docName" --parameters host="$RDS_ENDPOINT",portNumber="$REMOTE_PORT",localPortNumber="$LOCAL_PORT"

printMessage "Stopping ECS task, please wait ..." yellow

stopTaskCmd=$(aws ecs stop-task --cluster "$clusterName" --task "$taskId" --reason "Session ended by CLI" --output json | jq -r '.task.desiredStatus')

if [ "$stopTaskCmd" != "STOPPED" ]; then
    printMessage "Failed to stop task. It will automatically be shut down after a while." yellow
    exit 1
fi

printMessage "Task stopped successfully" green

exit 0
