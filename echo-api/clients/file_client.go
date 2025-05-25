package clients

import (
	"context"
	"io"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	uploadpb "echo-api/proto/upload"
    downloadpb "echo-api/proto/download"
)

type FileClient struct {
	uploadClient   uploadpb.FileUploadClient
	downloadClient downloadpb.FileDownloadClient
}

func NewFileClient() (*FileClient, error) {
	// Connect to upload service
	uploadConn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	// Connect to download service
	downloadConn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &FileClient{
		uploadClient:   uploadpb.NewFileUploadClient(uploadConn),
		downloadClient: downloadpb.NewFileDownloadClient(downloadConn),
	}, nil
}

// UploadFile uploads a file to the upload service
func (c *FileClient) UploadFile(ctx context.Context, filePath string, token string) (*uploadpb.UploadFileResponse, error) {
	// Add token to context
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", token)

	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	// Create upload stream
	stream, err := c.uploadClient.UploadFile(ctx)
	if err != nil {
		return nil, err
	}

	// Send metadata
	err = stream.Send(&uploadpb.UploadFileRequest{
		Data: &uploadpb.UploadFileRequest_Metadata{
			Metadata: &uploadpb.FileMetadata{
				Filename:    fileInfo.Name(),
				ContentType: "application/octet-stream", // TODO: Implement proper content type detection
				Size:       fileInfo.Size(),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	// Send file in chunks
	buffer := make([]byte, 1024*1024) // 1MB chunks
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		err = stream.Send(&uploadpb.UploadFileRequest{
			Data: &uploadpb.UploadFileRequest_Chunk{
				Chunk: buffer[:n],
			},
		})
		if err != nil {
			return nil, err
		}
	}

	// Get response
	return stream.CloseAndRecv()
}

// DownloadFile downloads a file from the download service
func (c *FileClient) DownloadFile(ctx context.Context, fileID string, token string, outputPath string) error {
	// Add token to context
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", token)

	// Create download request
	req := &downloadpb.DownloadFileRequest{
		FileId: fileID,
	}

	// Start download stream
	stream, err := c.downloadClient.DownloadFile(ctx, req)
	if err != nil {
		return err
	}

	// Create output file
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Receive file chunks
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Handle metadata
		if metadata := resp.GetMetadata(); metadata != nil {
			log.Printf("Downloading file: %s (%d bytes)", metadata.Filename, metadata.Size)
			continue
		}

		// Write chunk to file
		if chunk := resp.GetChunk(); chunk != nil {
			_, err = file.Write(chunk)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ListFiles lists all files for a user
func (c *FileClient) ListFiles(ctx context.Context, token string) (*downloadpb.ListFilesResponse, error) {
	// Add token to context
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", token)

	// Create list request
	req := &downloadpb.ListFilesRequest{}

	// Get file list
	return c.downloadClient.ListFiles(ctx, req)
} 