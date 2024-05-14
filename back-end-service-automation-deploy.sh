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

if [ "$#" -lt 8 ]; then
cat <<EOF | cat -


>> Script usage: $0 images_name app_name app_namespace replica ssh_user ssh_host ssh_port target_dir

Where:
    - images_name: specify the name of the image will be built
    - app_name: the app name when upgrading helm chart
    - app_namespace: K8S namespace to build on
    - replica: number of application instance will be run
    - ssh_user: username of target host to run
    - ssh_host: ip or domain of target host to run
    - ssh_port: port of target host to run
    - target_dir: the directory that command will be execute each time as well as the directory that will be store chart folder

EOF
exit 1
fi

lastest_git_commit_hash_id=$(git log -n 1 --pretty=format:'%h')


images_name="$1"
app_name="$2"
app_namespace="$3"
replica="$4"
ssh_user="$5"
ssh_host="$6"
ssh_port="$7"
target_dir="$8"
current_time=$(date -d "$b 0 min" "+%Y%m%d%H%M%S")
images_tag="${current_time}_${lastest_git_commit_hash_id}"

print_message "Deloying new version of service with images tag: ${images_tag}"

cd ./service
print_message_and_run_input_command "Downloading library with go mod" "go mod download"
print_message_and_run_input_command "Building go project to executable program" "CGO_ENABLED=0 GOOS=linux go build -o ./go_app"
cd ..

print_message "Updating chart information"
cat <<EOF | cat - | tee ./charts/fast-storage-back-end-helm-chart/Chart.yaml
apiVersion: v2
name: ${app_name}
description: A Helm chart for Kubernetes to deploy the ${app_name} service
type: application
version: 1.0.0
appVersion: ${images_tag}
EOF

print_message "Add helm chart note"
git_commit_message_and_commit_id=$(git log -1)
cat <<EOF | cat - | tee ./charts/fast-storage-back-end-helm-chart/templates/NOTES.txt
${git_commit_message_and_commit_id}
EOF

print_message "Uploading necessary file to target host $ssh_host"
ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no $ssh_user@$ssh_host -p $ssh_port "mkdir -p ${target_dir}"
scp -P $ssh_port -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -r ./charts/fast-storage-back-end-helm-chart/ $ssh_user@$ssh_host:$target_dir
scp -P $ssh_port -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -r ./service/additional_source_code/ $ssh_user@$ssh_host:$target_dir
scp -P $ssh_port -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no ./service/Dockerfile $ssh_user@$ssh_host:$target_dir
scp -P $ssh_port -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no ./service/go_app $ssh_user@$ssh_host:$target_dir

print_message "Reverting chart information"
cat <<EOF | cat - | tee ./charts/fast-storage-back-end-helm-chart/Chart.yaml
apiVersion: v2
name: app_name
description: A Helm chart for Kubernetes to deploy the app_name service
type: application
version: 1.0.0
appVersion: app_version
EOF

print_message "Removing chart note"
try_catch "rm -f ./charts/fast-storage-back-end-helm-chart/templates/NOTES.txt"

print_message "Removing built app"
try_catch "rm -f ./service/go_app"

print_message "Starting to build on remote host"

docker_build_command="docker build -f ./Dockerfile -t ${images_name}:${images_tag} ."
print_message_and_command_with_out_execute "Building image with Docker" "${docker_build_command}"
ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no $ssh_user@$ssh_host -p $ssh_port "cd ${target_dir} ; source ~/.bash_profile ; echo ${docker_build_command} > docker_build_command.txt ; eval ${docker_build_command}"

docker_push_command="docker push ${images_name}:${images_tag}"
print_message_and_command_with_out_execute "Pushing images to image registry" "${docker_push_command}"
ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no $ssh_user@$ssh_host -p $ssh_port "cd ${target_dir} ; source ~/.bash_profile ; echo ${docker_push_command} > docker_push_command.txt ; eval ${docker_push_command}"

helm_upgrade_command="helm upgrade -i --force --set image.name=${images_name},image.tag=${images_tag},replica=${replica},port=6060 ${app_name} -n ${app_namespace} --create-namespace ./fast-storage-back-end-helm-chart"
print_message_and_command_with_out_execute "Upgrading helm chart of application" "${helm_upgrade_command}"
try_catch "ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no $ssh_user@$ssh_host -p $ssh_port \"cd ${target_dir} ; source ~/.bash_profile ; echo ${helm_upgrade_command} > helm_upgrade_command.txt ; eval ${helm_upgrade_command}\""

docker_remove_image_command="docker rmi ${images_name}:${images_tag}"
print_message_and_command_with_out_execute "Removing built images" "${docker_remove_image_command}"
try_catch "ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no $ssh_user@$ssh_host -p $ssh_port \"cd ${target_dir} ; source ~/.bash_profile ; eval ${docker_remove_image_command}\""

print_message "Deployment process has been done"
