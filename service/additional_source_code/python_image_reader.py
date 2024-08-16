import sys
from PIL import Image # pip install Pillow
import io
import base64
import os

# Function to check if a file is an image
def is_image(file_path):
    try:
        Image.open(file_path)
        return True
    except IOError:
        return False

# Main function
def main(file_path):
    # Check if the file is an image
    if is_image(file_path):
        print("true")
        # Open the image file
        img = Image.open(file_path)

        # Create a BytesIO object to hold the image data in memory
        buffered = io.BytesIO()

        # Save the image to the BytesIO object in WebP format
        img.save(buffered, format="WEBP", optimize = True, quality = 10)

        # Get the byte data from the BytesIO object
        img_byte_data = buffered.getvalue()

        # Encode the byte data to base64
        img_base64 = base64.b64encode(img_byte_data).decode("utf-8")

        # Print the base64 string
        print(img_base64)
    else:
        print("false")

# Entry point of the script
if __name__ == "__main__":
    file_path = sys.argv[1]
    main(file_path)
