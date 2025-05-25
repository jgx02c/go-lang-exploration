package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	

	"github.com/labstack/echo/v4"
	"echo-api/clients"
)

// UploadFile handles file upload requests
func UploadFile(c echo.Context) error {
	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No file uploaded"})
	}

	// Create temp file
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to open uploaded file"})
	}
	defer src.Close()

	// Create temp directory if it doesn't exist
	if err := os.MkdirAll("temp", 0755); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create temp directory"})
	}

	// Save to temp file
	tempPath := filepath.Join("temp", file.Filename)
	dst, err := os.Create(tempPath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create temp file"})
	}
	defer dst.Close()
	defer os.Remove(tempPath) // Clean up temp file

	if _, err = io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save uploaded file"})
	}

	// Initialize file client
	fileClient, err := clients.NewFileClient()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to connect to file service"})
	}

	// Upload file
	token := c.Request().Header.Get("Authorization")
	resp, err := fileClient.UploadFile(c.Request().Context(), tempPath, token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to upload file: %v", err)})
	}

	return c.JSON(http.StatusOK, resp)
}

// DownloadFile handles file download requests
func DownloadFile(c echo.Context) error {
	// Get file ID from URL
	fileID := c.Param("id")
	if fileID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "File ID is required"})
	}

	// Initialize file client
	fileClient, err := clients.NewFileClient()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to connect to file service"})
	}

	// Create temp directory if it doesn't exist
	if err := os.MkdirAll("temp", 0755); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create temp directory"})
	}

	// Download file
	tempPath := filepath.Join("temp", fileID)
	token := c.Request().Header.Get("Authorization")
	if err := fileClient.DownloadFile(c.Request().Context(), fileID, token, tempPath); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to download file: %v", err)})
	}
	defer os.Remove(tempPath) // Clean up temp file

	// Send file to client
	return c.File(tempPath)
}

// ListFiles handles file listing requests
func ListFiles(c echo.Context) error {
	// Initialize file client
	fileClient, err := clients.NewFileClient()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to connect to file service"})
	}

	// List files
	token := c.Request().Header.Get("Authorization")
	resp, err := fileClient.ListFiles(c.Request().Context(), token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to list files: %v", err)})
	}

	return c.JSON(http.StatusOK, resp)
} 