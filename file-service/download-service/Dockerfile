FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install protoc and required tools
RUN apk add --no-cache protobuf-dev make git

# Copy go mod files
COPY go.mod ./
RUN go mod download

# Copy source code
COPY . .

# Generate protobuf code
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
RUN protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/file.proto

# Update dependencies and generate go.sum
RUN go mod tidy

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/download-service

# Create final image
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/download-service .

# Create upload directory
RUN mkdir -p /app/uploads

# Run the service
CMD ["./download-service"] 