#!/bin/bash

# Example usage: ./back-end-service-automation-deploy.sh "tuanloc/fast-storage-service" "fast-storage-backend" "fs-service" "3" "ae403" "0.tcp.ap.ngrok.io" "17742" "/home/ae403/fs-service"

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
images_tag="${current_time}_${lastest_git_commit_hash_id}"

final_image_name="$images_name:$images_tag"

clear
print_message "Deloying new version of service with images tag: ${images_tag}"
cd ./schedule

print_message_and_run_input_command "Change docker compose image name" "sed -i \"s|image_name_of_encrypt_folder_sheduler|$final_image_name|\" docker-compose.yml"

print_message "Uploading necessary file to target host $ssh_host"
ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no $ssh_user@$ssh_host -p $ssh_port "mkdir -p ${target_dir}"
scp -P $ssh_port -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -r ./encrypt-folder/ $ssh_user@$ssh_host:$target_dir
scp -P $ssh_port -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no ./docker-compose.yaml $ssh_user@$ssh_host:$target_dir

down="docker compose down"
print_message_and_command_with_out_execute "Down running docker service" "$down"
ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no $ssh_user@$ssh_host -p $ssh_port "cd ${target_dir} ; source ~/.bash_profile ; echo ${down} > down.txt ; eval ${down}"

rmi="docker rmi $final_image_name"
print_message_and_command_with_out_execute "Remove built image" "$rmi"
try_catch "ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no $ssh_user@$ssh_host -p $ssh_port \"cd ${target_dir} ; source ~/.bash_profile ; echo ${rmi} > rmi.txt ; eval ${rmi}\""

build="docker buildx build -f ./encrypt-folder/Dockerfile -t $final_image_name ./encrypt-folder"
print_message_and_command_with_out_execute "Build image" "$build"
ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no $ssh_user@$ssh_host -p $ssh_port "cd ${target_dir} ; source ~/.bash_profile ; echo ${build} > build.txt ; eval ${build}"

up="docker compose up -d"
print_message_and_command_with_out_execute "Start docker service" "$up"
ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no $ssh_user@$ssh_host -p $ssh_port "cd ${target_dir} ; source ~/.bash_profile ; echo ${up} > up.txt ; eval ${up}"

print_message_and_run_input_command "Restore file changed" "git reset . ; git restore ."
cd ..
