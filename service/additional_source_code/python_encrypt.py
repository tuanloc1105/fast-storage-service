import os
import sys
from cryptography.fernet import Fernet


def load_key(secret_key_directory: str):
    return open(secret_key_directory, "rb").read()


def is_input_path_a_file(path: str):
    if os.path.isdir(path):
        return False
    elif os.path.isfile(path):
        return True
    else:
        raise Exception("The path is neither a file nor a directory.")


def encrypt_file(file_name):
    key = load_key()
    fernet = Fernet(key)
    
    with open(file_name, "rb") as file:
        # read all file data
        file_data = file.read()
        
    # encrypt data
    header = b'ENCRYPTED'
    data_to_encrypt = header + file_data
    encrypted_data = fernet.encrypt(data_to_encrypt)
    
    # write the encrypted file
    with open(file_name, "wb") as file:
        file.write(encrypted_data)


if __name__=="__main__":
    path_of_secret_key = sys.argv[1]
    path_of_file = sys.argv[2]
    if is_input_path_a_file(path_of_file):
        encrypt_file(path_of_secret_key, path_of_file)
    else:
        for root, dirs, files in os.walk(path_of_file):
            # Sort files by name
            files.sort()
            for file in files:
                file_path = os.path.join(root, file)
                try:
                    print(f'Working on file: {file_path}')
                    encrypt_file(path_of_secret_key, file_path)
                except (IOError, OSError):
                    pass
