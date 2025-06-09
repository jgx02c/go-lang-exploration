# 📁 Go File Upload/Download Service

A simple and efficient file upload and download service built with Go and Echo framework.

## 🚀 Features

- File upload with progress tracking
- Secure file downloads
- File metadata storage
- Basic file validation
- Simple REST API interface

## 🛠️ Tech Stack

- Go (Golang)
- Echo Framework
- PostgreSQL for metadata storage
- Docker for containerization

## 📋 Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- PostgreSQL (provided via Docker)

## 🏗️ Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── handlers/
│   ├── middleware/
│   ├── models/
│   └── storage/
├── docker-compose.yml
└── README.md
```

## 🚀 Getting Started

1. Clone the repository:
```bash
git clone <repository-url>
cd go-file-service
```

2. Start the services using Docker Compose:
```bash
docker-compose up -d
```

3. Run the application:
```bash
go run cmd/server/main.go
```

## 📡 API Endpoints

### Upload File
```http
POST /api/v1/upload
Content-Type: multipart/form-data

file: <file>
```

### Download File
```http
GET /api/v1/download/:fileId
```

### List Files
```http
GET /api/v1/files
```

## 🔧 Development

### Database Access
To access the PostgreSQL database:
```bash
docker exec -it echo_postgres psql -U postgres -d echo_api
```

To exit the PostgreSQL shell:
```bash
\q
```

### Running Tests
```bash
go test ./...
```

## 📝 License

MIT License

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
