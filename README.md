# 📁 Go File Storage

[![Go Version](https://img.shields.io/badge/Go-1.21%2B-blue)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue)](https://docker.com/)

> A simple, fast, and secure file storage server written in Go. Perfect for self-hosted file storage solutions!

## ✨ Features

- 🚀 **Fast & Lightweight** - Single binary, minimal dependencies
- 🔒 **Secure** - Path traversal protection, file validation
- 📦 **Easy Setup** - Works out of the box with Docker or binary
- 🌐 **REST API** - Simple HTTP endpoints for file operations
- 🔗 **Direct Links** - Get instant download links for uploaded files
- 📊 **File Management** - Upload, download, and organize files
- 🛡️ **Production Ready** - Built with security and performance in mind

## 🚀 Quick Start

### Using Docker (Recommended)

```bash
# Run with Docker
docker run -p 8080:8080 -v ./storage:/app/storage your-username/go-file-storage

# Access the API
curl -X POST -F "file=@example.txt" http://localhost:8080/api/upload
```

### From Source

```bash
# Clone the repository
git clone https://github.com/rapir-0/go-file-storage.git
cd go-file-storage

# Run the server
go run main.go

# Server will start on http://localhost:8080
```

## 📖 API Documentation

### Upload File
Upload a file to the storage server.

**Endpoint:** `POST /api/upload`  
**Content-Type:** `multipart/form-data`

**Request:**
```bash
curl -X POST \
  -F "file=@/path/to/your/file.jpg" \
  http://localhost:8080/api/upload
```

**Response:**
```json
{
  "success": true,
  "file_id": "a1b2c3d4e5f6",
  "original_name": "file.jpg",
  "download_url": "http://localhost:8080/api/download?filename=a1b2c3d4e5f6.jpg",
  "size": 1024567
}
```

### Download File
Download a file using the provided download link.

**Endpoint:** `GET /api/download`  
**Parameters:** `filename` (required)

**Request:**
```bash
curl -O "http://localhost:8080/api/download?filename=a1b2c3d4e5f6.jpg"
```

**Response:** 
File content with appropriate headers for download.

## 🛠️ Installation

### Option 1: Docker Compose

```yaml
# docker-compose.yml
version: '3.8'
services:
  file-storage:
    image: rapir-0/go-file-storage:latest
    ports:
      - "8080:8080"
    volumes:
      - ./storage:/app/storage
    environment:
      - SERVER_PORT=8080
      - MAX_FILE_SIZE=10MB
```

```bash
docker-compose up -d
```

### Option 2: Binary Release

1. Download the latest release from [Releases](https://github.com/rapir-0/go-file-storage/releases)
2. Extract and run:
   ```bash
   ./go-file-storage
   ```

### Option 3: Build from Source

```bash
# Prerequisites: Go 1.21+
git clone https://github.com/rapir-0/go-file-storage.git
cd go-file-storage

# Build
go build -o go-file-storage main.go

# Run
./go-file-storage
```

## ⚙️ Configuration

Configure the server using environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | `8080` | Port for the HTTP server |
| `STORAGE_PATH` | `./storage` | Directory to store uploaded files |
| `MAX_FILE_SIZE` | `10MB` | Maximum allowed file size |
| `SERVER_HOST` | `localhost` | Server hostname/IP address |

**Example:**
```bash
export SERVER_PORT=3000
export MAX_FILE_SIZE=50MB
./go-file-storage
```

## 🏗️ Project Structure

```
go-file-storage/
├── main.go              # Application entry point
├── handlers/            # HTTP request handlers
│   ├── upload.go       # File upload logic
│   └── download.go     # File download logic
├── storage/            # File storage directory
├── docker-compose.yml  # Docker setup
├── Dockerfile         # Docker image
└── README.md          # This file
```

## 🔒 Security Features

- **Path Traversal Protection** - Prevents access to files outside storage directory
- **File Validation** - Validates uploaded files and generates unique identifiers
- **Size Limits** - Configurable maximum file size limits
- **Safe File Names** - Generates unique file names to prevent conflicts

## 🧪 Testing

```bash
# Test file upload
curl -X POST \
  -F "file=@test.txt" \
  http://localhost:8080/api/upload

# Test file download (use the download_url from upload response)
curl -O "http://localhost:8080/api/download?filename=unique-id.txt"
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the project
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🌟 Support

If you find this project helpful, please give it a ⭐ on GitHub!

**Issues & Questions:** [GitHub Issues](https://github.com/rapir-0/go-file-storage/issues)

---

<div align="center">
  <sub>Built with ❤️ using Go</sub>
</div>
