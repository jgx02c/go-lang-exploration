package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "file-service/upload-service/proto"
)

type server struct {
	pb.UnimplementedFileUploadServer
	uploadDir string
}

func (s *server) UploadFile(stream pb.FileUpload_UploadFileServer) error {
	// Get user ID from JWT token
	userID, err := getUserIDFromContext(stream.Context())
	if err != nil {
		return err
	}

	// First message should contain metadata
	req, err := stream.Recv()
	if err != nil {
		return status.Error(codes.InvalidArgument, "failed to receive metadata")
	}

	metadata := req.GetMetadata()
	if metadata == nil {
		return status.Error(codes.InvalidArgument, "first message must contain metadata")
	}

	// Generate unique file ID
	fileID := uuid.New().String()
	filePath := filepath.Join(s.uploadDir, fileID)

	// Create file
	file, err := os.Create(filePath)
	if err != nil {
		return status.Error(codes.Internal, "failed to create file")
	}
	defer file.Close()

	// Process file chunks
	var totalSize int64
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Error(codes.Internal, "failed to receive chunk")
		}

		chunk := req.GetChunk()
		if chunk == nil {
			continue
		}

		_, err = file.Write(chunk)
		if err != nil {
			return status.Error(codes.Internal, "failed to write chunk")
		}
		totalSize += int64(len(chunk))
	}

	// Send response
	return stream.SendAndClose(&pb.UploadFileResponse{
		FileId:    fileID,
		Filename:  metadata.Filename,
		Size:      totalSize,
		CreatedAt: fmt.Sprintf("%d", totalSize), // TODO: Add proper timestamp
		UserId:    userID, // Add user ID to the response
	})
}

func (s *server) GetFileMetadata(ctx context.Context, req *pb.GetFileMetadataRequest) (*pb.FileMetadata, error) {
	// Get user ID from JWT token
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// TODO: Implement file metadata retrieval from database
	// For now, return dummy data with the actual user ID
	return &pb.FileMetadata{
		Filename:    "example.txt",
		ContentType: "text/plain",
		Size:        1024,
		UserId:      userID, // Use the actual user ID from the token
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

	// Create upload directory
	uploadDir := "uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}

	// Start gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterFileUploadServer(s, &server{uploadDir: uploadDir})

	log.Printf("Upload service listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
} 