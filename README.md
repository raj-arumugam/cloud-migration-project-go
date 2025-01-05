# Cloud Migration Project (Go)
A robust Go utility for seamless photo migration between cloud services (Amazon Photos to Google Drive). This tool automates the secure transfer of photos while maintaining folder structures and metadata.

## Prerequisites
- Go 1.21 or higher
- AWS Account with access to Amazon Photos
- Google Cloud Project with Drive API enabled
- AWS CLI (optional, for configuration)

## Configuration

### AWS Setup
1. Create an AWS account if you don't have one
2. Create IAM credentials with access to Amazon Photos
3. Note down your:
   - AWS Access Key ID
   - AWS Secret Access Key
   - AWS Region
   - Bucket Name

### Google Setup
1. Create a Google Cloud Project
2. Enable the Google Drive API
3. Create OAuth 2.0 credentials
4. Note down your:
   - Client ID
   - Client Secret

## Installation 

##Clone the repository
```
git clone https://github.com/yourusername/cloud-migration-project-go
cd cloud-migration-project-go
```

## Environment Variables
Create a `.env` file in the project root:

##AWS Configuration
```
AWS_ACCESS_KEY_ID=your_aws_access_key
AWS_SECRET_ACCESS_KEY=your_aws_secret_key
AWS_REGION=your_region
AWS_BUCKET_NAME=your_bucket_name
```

##Google Configuration
```
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
GOOGLE_TOKEN_PATH=path_to_your_token_file
``` 

## Running the Application
```
go run main.go
```

## Project Structure
- `main.go`: The main entry point for the application.
- `internal/cloud/aws/`: AWS-specific implementations.
- `internal/cloud/google/`: Google Drive-specific implementations.
- `internal/config/`: Configuration loading and struct definitions.
- `internal/utils/`: Utility functions for logging and error handling.

## Notes
- This tool is designed to migrate photos from Amazon Photos to Google Drive.
- It assumes that the photos are stored in the specified AWS bucket.
- The tool will create a new folder in Google Drive with the same name as the AWS bucket.
- The tool will migrate all photos in the specified AWS bucket.
- The tool will not migrate albums or other metadata.
- The tool will not migrate videos or other file types.
- The tool will not migrate photos that are already in the Google Drive folder.

## License
This project is open-sourced under the MIT License - see the LICENSE file for details.

## Contact
For any questions or feedback, please contact [rajkumar.arumugam@yahoo.com](mailto:rajkumar.arumugam@yahoo.com).
