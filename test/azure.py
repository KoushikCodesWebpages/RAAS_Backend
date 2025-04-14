from azure.storage.blob import BlobServiceClient
import os

# Azure Blob Storage connection string
connection_string = os.environ.get("AZURE_STORAGE_CONNECTION_STRING")
container_name = os.environ.get("AZURE_BLOB_CONTAINER")

def upload_blob(file_path, blob_name):
    try:
        # Create a blob service client
        blob_service_client = BlobServiceClient.from_connection_string(connection_string)
        
        # Get a client to interact with the specified container
        blob_client = blob_service_client.get_blob_client(container=container_name, blob=blob_name)
        
        # Upload the file to Azure Blob Storage
        with open(file_path, "rb") as data:
            blob_client.upload_blob(data)
        
        print(f"File {blob_name} uploaded successfully.")
        
    except Exception as e:
        print(f"Error uploading file: {e}")

def download_blob(blob_name, download_path):
    try:
        # Create a blob service client
        blob_service_client = BlobServiceClient.from_connection_string(connection_string)
        
        # Get a client to interact with the specified container
        blob_client = blob_service_client.get_blob_client(container=container_name, blob=blob_name)
        
        # Download the blob to a local file
        with open(download_path, "wb") as download_file:
            blob_data = blob_client.download_blob()
            blob_data.readinto(download_file)
        
        print(f"File {blob_name} downloaded successfully.")
        
    except Exception as e:
        print(f"Error downloading file: {e}")

def list_blobs():
    try:
        # Create a blob service client
        blob_service_client = BlobServiceClient.from_connection_string(connection_string)
        
        # Get a client to interact with the specified container
        container_client = blob_service_client.get_container_client(container_name)
        
        # List the blobs in the container
        blob_list = container_client.list_blobs()
        
        for blob in blob_list:
            print(blob.name)
        
    except Exception as e:
        print(f"Error listing blobs: {e}")

def delete_blob(blob_name):
    try:
        # Create a blob service client
        blob_service_client = BlobServiceClient.from_connection_string(connection_string)
        
        # Get a client to interact with the specified container
        blob_client = blob_service_client.get_blob_client(container=container_name, blob=blob_name)
        
        # Delete the blob
        blob_client.delete_blob()
        
        print(f"Blob {blob_name} deleted successfully.")
        
    except Exception as e:
        print(f"Error deleting blob: {e}")

# Test the functions
file_path = "test.txt"
blob_name = "test-blob.txt"
download_path = "downloaded-test.txt"

# Create a test file
with open(file_path, "w") as f:
    f.write("This is a test file.")

# Upload the file
upload_blob(file_path, blob_name)

# List the blobs
print("Blobs in the container:")
list_blobs()

# Download the blob
download_blob(blob_name, download_path)

# Delete the blob
delete_blob(blob_name)