package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Blank import for PostgreSQL driver
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "file-service/download-service/proto"
)

type server struct {
	pb.UnimplementedFileDownloadServer
	uploadDir string
	db        *sql.DB
}

func (s *server) DownloadFile(req *pb.DownloadFileRequest, stream pb.FileDownload_DownloadFileServer) error {
	// Get user ID from JWT token
	userID, err := getUserIDFromContext(stream.Context())
	if err != nil {
		return err
	}

	// Open file
	filePath := filepath.Join(s.uploadDir, req.FileId)
	file, err := os.Open(filePath)
	if err != nil {
		return status.Error(codes.NotFound, "file not found")
	}
	defer file.Close()

	// Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		return status.Error(codes.Internal, "failed to get file info")
	}

	// Send metadata first
	err = stream.Send(&pb.DownloadFileResponse{
		Data: &pb.DownloadFileResponse_Metadata{
			Metadata: &pb.FileMetadata{
				FileId:      req.FileId,
				Filename:    fileInfo.Name(),
				ContentType: "application/octet-stream", // TODO: Implement proper content type detection
				Size:        fileInfo.Size(),
				CreatedAt:   fileInfo.ModTime().String(),
				UserId:      userID,
			},
		},
	})
	if err != nil {
		return status.Error(codes.Internal, "failed to send metadata")
	}

	// Send file in chunks
	buffer := make([]byte, 1024*1024) // 1MB chunks
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Error(codes.Internal, "failed to read file")
		}

		err = stream.Send(&pb.DownloadFileResponse{
			Data: &pb.DownloadFileResponse_Chunk{
				Chunk: buffer[:n],
			},
		})
		if err != nil {
			return status.Error(codes.Internal, "failed to send chunk")
		}
	}

	return nil
}

func (s *server) ListFiles(ctx context.Context, req *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {
	// Get user ID from JWT token
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	log.Printf("Listing files for user ID: %s", userID)

	// Query files from database
	query := `
		SELECT 
			id::text, 
			filename, 
			content_type, 
			size, 
			created_at 
		FROM files 
		WHERE user_id = $1::integer
	`
	log.Printf("Executing query: %s with user_id=%s", query, userID)
	
	rows, err := s.db.Query(query, userID)
	if err != nil {
		log.Printf("Database query error: %v", err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to query files: %v", err))
	}
	defer rows.Close()

	var files []*pb.FileMetadata
	for rows.Next() {
		var file pb.FileMetadata
		var createdAt time.Time
		err := rows.Scan(&file.FileId, &file.Filename, &file.ContentType, &file.Size, &createdAt)
		if err != nil {
			log.Printf("Row scan error: %v", err)
			return nil, status.Error(codes.Internal, fmt.Sprintf("failed to scan file row: %v", err))
		}
		file.UserId = userID
		file.CreatedAt = createdAt.Format(time.RFC3339)
		files = append(files, &file)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Row iteration error: %v", err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("error iterating file rows: %v", err))
	}

	log.Printf("Found %d files for user %s", len(files), userID)
	return &pb.ListFilesResponse{
		Files: files,
	}, nil
}

func getUserIDFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "metadata is not provided")
	}

	tokens := md.Get("authorization")
	if len(tokens) == 0 {
		return "", status.Error(codes.Unauthenticated, "authorization token is not provided")
	}

	token := tokens[0]
	// Remove "Bearer " prefix if present
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}
	claims := jwt.MapClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !parsedToken.Valid {
		return "", status.Error(codes.Unauthenticated, "invalid token")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return "", status.Error(codes.Internal, "invalid user ID in token")
	}

	return fmt.Sprintf("%d", int(userID)), nil
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Create upload directory (shared with upload service)
	uploadDir := "uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}

	// Initialize database connection
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Start gRPC server
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterFileDownloadServer(s, &server{
		uploadDir: uploadDir,
		db:        db,
	})

	log.Printf("Download service listening on :50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
} 