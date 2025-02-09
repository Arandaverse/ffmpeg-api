# FFMPEG Serverless API

A robust and scalable serverless API for processing videos using FFMPEG. This service provides a RESTful interface for video processing operations with comprehensive job management, storage options, and user authentication.

## Features

- **Advanced Video Processing**:

  - Asynchronous FFMPEG command execution
  - Real-time progress tracking (0-100%)
  - Support for multiple input and output files
  - Custom FFMPEG command templating
  - Detailed job status monitoring
  - Processing time metrics

- **Comprehensive Storage System**:

  - Local storage for development
  - MinIO integration for production
  - Automatic bucket creation and policy management
  - Support for external URLs
  - User-specific storage paths
  - Temporary file management

- **Authentication & Security**:

  - User registration and login
  - Secure password hashing with bcrypt
  - API token-based authentication
  - User-specific resource isolation
  - Usage tracking and statistics

- **Job Management**:

  - Unique job UUID tracking
  - Detailed job status and progress
  - File metadata collection
  - Processing time statistics
  - Error handling and reporting

- **File Processing Features**:
  - Multi-file input support
  - Format detection
  - Image dimension extraction
  - File size tracking
  - Progress monitoring
  - Cleanup of temporary files

## Prerequisites

- Go 1.23.3 or higher
- FFMPEG installed on the system
- SQLite (for development)
- MinIO (optional, for production storage)

## Configuration

The application uses environment variables for configuration. Create a `.env` file in the root directory:

```env
# Server Configuration
SERVER_PORT=8000
API_TOKEN_LENGTH=32

# Database Configuration
DB_DRIVER=sqlite
DB_DSN=ffmpeg_api.db

# FFMPEG Configuration
FFMPEG_PATH=/usr/bin/ffmpeg
TEMP_DIR=tmp
PROGRESS_UPDATE_INTERVAL=5

# Storage Configuration
STORAGE_PROVIDER=local  # or 'minio'
MINIO_ENDPOINT=127.0.0.1
MINIO_PORT=56732
MINIO_ACCESS_KEY=your_access_key
MINIO_SECRET_KEY=your_secret_key
MINIO_USE_SSL=false
MINIO_BUCKET_NAME=ffmpeg-files
MINIO_REGION=us-east-1
MINIO_BUCKET_URL=http://your-minio-url
```

## Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/ffmpeg-serverless.git
cd ffmpeg-serverless
```

2. Install dependencies:

```bash
go mod download
```

3. Build the application:

```bash
go build -o ffmpeg-api ./cmd/main.go
```

## API Endpoints

### Authentication

- **Register User**

  ```http
  POST /register
  Content-Type: application/json

  {
    "username": "user",
    "password": "password",
    "email": "user@example.com"
  }
  ```

  Response:

  ```json
  {
    "api_token": "your_api_token"
  }
  ```

- **Login**

  ```http
  POST /login
  Content-Type: application/json

  {
    "username": "user",
    "password": "password"
  }
  ```

  Response:

  ```json
  {
    "api_token": "your_api_token"
  }
  ```

### Video Processing

- **Process Video**

  ```http
  POST /ffmpeg
  X-API-Token: your_api_token
  Content-Type: application/json

  {
    "ffmpeg_command": "ffmpeg -i {{input}} {{output}}",
    "input_files": {
      "input": "https://example.com/input.mp4"
    },
    "output_files": {
      "output": "output.mp4"
    }
  }
  ```

  Response:

  ```json
  {
    "uuid": "job_uuid",
    "status": "pending"
  }
  ```

- **Check Job Status**
  ```http
  GET /ffmpeg/progress/{uuid}
  X-API-Token: your_api_token
  ```
  Response:
  ```json
  {
    "uuid": "job_uuid",
    "status": "PROCESSING",
    "progress": 45,
    "result": "",
    "created_at": "2024-03-21T10:00:00Z",
    "updated_at": "2024-03-21T10:01:00Z",
    "total_processing_seconds": 60,
    "ffmpeg_command_run_seconds": 45,
    "output_files": {
      "output": {
        "file_id": "unique_id",
        "file_type": "video",
        "file_format": "mp4",
        "size_mbytes": 10.5,
        "storage_url": "http://storage/path/to/file.mp4"
      }
    }
  }
  ```

## Project Structure

```
.
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── config/                 # Configuration management
│   ├── domain/                 # Domain models and interfaces
│   ├── dto/                    # Data Transfer Objects
│   ├── handlers/               # HTTP handlers and routing
│   │   └── routes/            # Route definitions
│   ├── logger/                # Logging utilities
│   ├── repository/            # Data access layer
│   ├── response/              # HTTP response utilities
│   ├── service/               # Business logic
│   │   ├── auth_service.go    # Authentication service
│   │   ├── ffmpeg_service.go  # FFMPEG processing service
│   │   └── storage_service.go # File storage service
│   └── validation/            # Request validation
├── docs/                      # API documentation
├── go.mod
├── go.sum
└── README.md
```

## Error Handling

The API uses standardized error responses:

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Error description"
  }
}
```

Common error codes:

- `BAD_REQUEST`: Invalid request parameters
- `UNAUTHORIZED`: Missing or invalid API token
- `FORBIDDEN`: Insufficient permissions
- `NOT_FOUND`: Resource not found
- `INTERNAL_SERVER_ERROR`: Server-side error
- `VALIDATION_ERROR`: Request validation failed
- `INVALID_CREDENTIALS`: Wrong username or password
- `STORAGE_ERROR`: File storage/retrieval error
- `FFMPEG_ERROR`: Video processing error

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
