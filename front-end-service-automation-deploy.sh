#!/bin/bash

# Example usage: ./front-end-service-automation-deploy.sh "tuanloc/fast-storage-service-web" "ae403" "0.tcp.ap.ngrok.io" "17742" "/home/ae403/fs-web"

handle_error() {
    echo -e "\n\n >> An error occurred on line $1 \n\n"
    exit 1
}

trap 'handle_error $LINENO' ERR

print_message_and_run_input_command() {
cat <<EOF | cat -

>>> $1
>>> Command: $2

EOF
eval $2
}

print_message_and_command_with_out_execute() {
cat <<EOF | cat -

>>> $1
>>> Command: $2

EOF
}

print_message() {
cat <<EOF | cat -

>>> $1

EOF
}

try_catch() {
    print_message "Try to execute command: $1"
    if eval $1
    then
        print_message "Execute command successfuly"
    else
        print_message "Execute command failed. Tryting to execute the second command"
        if eval $2
        then
            print_message "Execute second command successfuly"
        else
            print_message "Execute second command failed"
        fi
    fi
}

if [ "$#" -lt 5 ]; then
cat <<EOF | cat -


>> Script usage: $0 images_name app_name app_namespace replica ssh_user ssh_host ssh_port target_dir

Where:
    - images_name: specify the name of the image will be built
    - ssh_user: username of target host to run
    - ssh_host: ip or domain of target host to run
    - ssh_port: port of target host to run
    - target_dir: the directory that command will be execute each time as well as the directory that will be store chart folder

EOF
exit 1
fi

lastest_git_commit_hash_id=$(git log -n 1 --pretty=format:'%h')


images_name="$1"
ssh_user="$2"
ssh_host="$3"
ssh_port="$4"
target_dir="$5"
current_time=$(date -d "$b 0 min" "+%Y%m%d%H%M%S")
# images_tag="${current_time}_${lastest_git_commit_hash_id}"
images_tag="latest"

# print_message "Deloying new version of service with images tag: ${images_tag}"
clear
cd ./web

print_message_and_run_input_command "Remove node modules" "rm -rf ./node_modules"
print_message_and_run_input_command "Remove built app" "rm -rf ./app-run"
print_message_and_run_input_command "Remove nx cache" "rm -rf ./.nx"
print_message_and_run_input_command "Remove angular cache" "rm -rf ./.angular"

print_message_and_run_input_command "Install dependencies" "pnpm install"
print_message_and_run_input_command "Build web" "npm run build:angular"

print_message_and_run_input_command "Rename built folder" "mv ./dist/ ./app-run/"

print_message "Uploading necessary file to target host $ssh_host"
ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no $ssh_user@$ssh_host -p $ssh_port "mkdir -p ${target_dir}"
scp -P $ssh_port -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -r ./app-run/ $ssh_user@$ssh_host:$target_dir
scp -P $ssh_port -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no ./nginx.conf $ssh_user@$ssh_host:$target_dir
scp -P $ssh_port -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no ./docker-compose.yml $ssh_user@$ssh_host:$target_dir
scp -P $ssh_port -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no ./Dockerfile.copy_from_locally_built_app $ssh_user@$ssh_host:$target_dir


down="docker compose down"
print_message_and_command_with_out_execute "Down running docker service" "$down"
ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no $ssh_user@$ssh_host -p $ssh_port "cd ${target_dir} ; source ~/.bash_profile ; echo ${down} > down.txt ; eval ${down}"

rmi="docker rmi fast-storage-service:latest"
print_message_and_command_with_out_execute "Remove built image" "$rmi"
try_catch "ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no $ssh_user@$ssh_host -p $ssh_port \"cd ${target_dir} ; source ~/.bash_profile ; echo ${rmi} > rmi.txt ; eval ${rmi}\""

build="docker buildx build -f ./Dockerfile.copy_from_locally_built_app -t fast-storage-service:latest ."
print_message_and_command_with_out_execute "Build image" "$build"
ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no $ssh_user@$ssh_host -p $ssh_port "cd ${target_dir} ; source ~/.bash_profile ; echo ${build} > build.txt ; eval ${build}"

up="docker compose up -d"
print_message_and_command_with_out_execute "Start docker service" "$up"
ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no $ssh_user@$ssh_host -p $ssh_port "cd ${target_dir} ; source ~/.bash_profile ; echo ${up} > up.txt ; eval ${up}"

print_message_and_run_input_command "Remove node modules" "rm -rf ./node_modules"
print_message_and_run_input_command "Remove built app" "rm -rf ./app-run"
print_message_and_run_input_command "Remove nx cache" "rm -rf ./.nx"
print_message_and_run_input_command "Remove angular cache" "rm -rf ./.angular"
cd ..

print_message_and_run_input_command "Restore file changed" "git reset . ; git restore ."
