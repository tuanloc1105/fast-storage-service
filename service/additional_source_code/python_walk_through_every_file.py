import os
import sys

if __name__ == "__main__":
    input_path_of_directory = sys.argv[1]
    if not os.path.isdir(input_path_of_directory):
        print('input path is not directory')
        sys.exit(1)
    for root, dirs, files in os.walk(input_path_of_directory):
        if not dirs and not files:
            print('empty directory')
            sys.exit(0)
        for file in files:
            file_path = os.path.join(root, file)
            try:
                print(f"{file_path}")
            except (IOError, OSError):
                pass
