#!/bin/bash

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
    print_message_with_echo "Try to execute command: $1"
    if eval $1
    then
        print_message_with_echo "Execute command successfuly"
    else
        print_message_with_echo "Execute command failed. Tryting to execute the second command"
        if eval $2
        then
            print_message_with_echo "Execute second command successfuly"
        else
            print_message_with_echo "Execute second command failed"
        fi
    fi
}
