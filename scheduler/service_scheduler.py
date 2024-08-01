import schedule
import time
import os
import requests
from datetime import datetime


def call_fast_storage_api_to_encrypt_all_folder():
    now = datetime.now()
    time_formatted_time_vietnam_format = now.strftime("%d/%m/%Y %H:%M:%S.%f")[:-3]
    encrypt_folder_api_url = os.getenv("ENCRYPT_FOLDER_API_URL")
    encrypt_folder_api_key = os.getenv("ENCRYPT_FOLDER_API_KEY")
    print(f"Task is running at {time_formatted_time_vietnam_format}...")
    request_payload: dict[str, any] = {
        "request": {
            "encrypt": True
        }
    }
    request_header: dict[str, any] = {
        "api-key": encrypt_folder_api_key
    }
    response: requests.Response = requests.post(encrypt_folder_api_url, headers=request_header, json=request_payload, verify=False)
    print(f"    - Request URL: {encrypt_folder_api_url}")
    print(f"    - Request Headers: {request_header}")
    print(f"    - Request Data: {request_payload}")
    print(f"    - Response Status Code: {response.status_code}")
    print(f"    - Response Headers: {response.headers}")
    print(f"    - Response Text: {response.text}")



# Schedule the task to run every hour
# schedule.every().hour.do(call_fast_storage_api_to_encrypt_all_folder)
schedule.every().minute.at(":00").do(call_fast_storage_api_to_encrypt_all_folder)

# Keep the script running
while True:
    schedule.run_pending()
    time.sleep(1)
