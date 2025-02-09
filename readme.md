# FFMPEG Serverless API

A robust and scalable serverless API for processing videos using FFMPEG. This service provides a RESTful interface for video processing operations with support for both local and MinIO storage backends.

## Features

- **Video Processing**: Process videos using FFMPEG with customizable commands
- **Storage Options**:
  - Local storage for development
  - MinIO integration for production-ready object storage
- **Authentication**:
  - User registration and login
  - API token-based authentication
- **Job Management**:
  - Asynchronous video processing
  - Job status tracking and progress monitoring
- **Scalable Architecture**:
  - Modular design
  - Clean separation of concerns
  - Interface-based dependencies

## Prerequisites

- Go 1.23.3 or higher
- FFMPEG installed on the system
- SQLite (for development)
- MinIO (optional, for production storage)
- Access to S3-compatible storage:
  - Either MinIO instance running locally/remotely
  - Or accessible S3 bucket with proper credentials
  - Source videos must be stored in S3-compatible storage
  - Output videos will be stored in the same storage

## Configuration

The application uses environment variables for configuration. Create a `.env` file in the root directory with the following variables:

```env
# Server Configuration
SERVER_PORT=
API_TOKEN_LENGTH=

# Database Configuration
DB_DRIVER=
POSTGRES_USER=
POSTGRES_PASSWORD=
POSTGRES_DB=
POSTGRES_PORT=
DB_URI=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DB}

# FFMPEG Configuration
FFMPEG_PATH=
TEMP_DIR=
PROGRESS_UPDATE_INTERVAL=

# Storage Configuration
STORAGE_PROVIDER=
MINIO_ENDPOINT=
MINIO_PORT=
MINIO_ACCESS_KEY=
MINIO_SECRET_KEY=
MINIO_USE_SSL=
MINIO_REGION=
MINIO_BUCKET_NAME=
MINIO_BUCKET_URL=
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

## Usage

### Starting the Server

```bash
./ffmpeg-api
```

### API Endpoints

#### Authentication

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

- **Login**

  ```http
  POST /login
  Content-Type: application/json

  {
    "username": "user",
    "password": "password"
  }
  ```

#### Video Processing

- **Process Video**

  ```http
  POST /ffmpeg
  X-API-Token: your_api_token
  Content-Type: application/json

  {
    "command": "custom_ffmpeg_command",
    "s3_file_url": "https://example.com/video.mp4",
    "format": "mp4",
    "quality": "high"
  }
  ```

- **Check Job Status**

  ```http
  GET /ffmpeg/progress/{uuid}
  X-API-Token: your_api_token
  ```

## Project Structure

```
.
├── cmd/
│   └── main.go           # Application entry point
├── internal/
│   ├── config/           # Configuration management
│   ├── domain/           # Domain models
│   ├── dto/              # Data Transfer Objects
│   ├── handlers/         # HTTP handlers
│   ├── logger/           # Logging utilities
│   ├── repository/       # Data access layer
│   ├── response/         # HTTP response utilities
│   ├── service/          # Business logic
│   └── validation/       # Request validation
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

- `BAD_REQUEST`
- `UNAUTHORIZED`
- `FORBIDDEN`
- `NOT_FOUND`
- `INTERNAL_SERVER_ERROR`
- `VALIDATION_ERROR`
- `INVALID_CREDENTIALS`

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
