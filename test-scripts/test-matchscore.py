import os
import requests
from dotenv import load_dotenv

# Load environment variables from the .env file
load_dotenv()

# Get the Hugging Face API token from the environment variables
API_TOKEN = os.getenv('HF_API_KEY')

# Check if the token is loaded correctly
if not API_TOKEN:
    print("Error: Hugging Face API token is not set. Please check your .env file.")
    exit(1)

# List all models from the environment variables
models = [
    os.getenv('HF_MODEL_1'),
    os.getenv('HF_MODEL_2'),
    os.getenv('HF_MODEL_3'),
    os.getenv('HF_MODEL_4'),
    os.getenv('HF_MODEL_5'),
    os.getenv('HF_MODEL_6'),
    os.getenv('HF_MODEL_7'),
    os.getenv('HF_MODEL_8'),
    os.getenv('HF_MODEL_9'),
    os.getenv('HF_MODEL_10')
]

# Define the Hugging Face API URL
API_URL = "https://api-inference.huggingface.co/models/"

# Function to get match score from Hugging Face API
def get_match_score(model, text_1, text_2):
    headers = {
        "Authorization": f"Bearer {API_TOKEN}"
    }
    payload = {
        "inputs": [text_1, text_2]
    }
    
    # Print payload to verify correct format
    print(f"Request Payload: {payload}")
    
    response = requests.post(API_URL + model, headers=headers, json=payload)

    if response.status_code == 200:
        result = response.json()
        return result[0]['score']  # The similarity score returned by Hugging Face
    else:
        print(f"Error with model {model}: {response.status_code}")
        print(f"Response: {response.text}")  # Print response body for debugging
        return None


# Example texts (Seeker Profile and Job Description)
seeker_profile = "Software Engineer with 5 years of experience in backend development, working with Java, Spring Boot, and Microservices."
job_description = "We are looking for a Backend Developer with expertise in Java and Spring Boot to work on our enterprise-level microservices platform."

# Test match score for each model
for model in models:
    print(f"Calculating match score using model: {model}")
    score = get_match_score(model, seeker_profile, job_description)
    
    if score is not None:
        print(f"Match score for {model}: {score:.4f}")
    else:
        print(f"Failed to calculate score for {model}")
