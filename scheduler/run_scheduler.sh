#!/bin/bash

export ENCRYPT_FOLDER_API_URL="http://localhost:8080/fast_storage/api/v1/storage/crypto_every_folder"
export ENCRYPT_FOLDER_API_KEY="DnX521Wks684SoF097rAP925TjU555Dp"

python3.12 service_scheduler.py
