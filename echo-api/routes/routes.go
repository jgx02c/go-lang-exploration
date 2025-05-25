package routes

import (
	"os"

	"github.com/labstack/echo/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"echo-api/handlers"
)

func Setup(e *echo.Echo) {
	// Public routes
	e.POST("/register", handlers.Register)
	e.POST("/login", handlers.Login)

	// Protected group
	r := e.Group("/profile")
	r.Use(echojwt.JWT([]byte(os.Getenv("JWT_SECRET"))))
	r.GET("", handlers.Profile)

	// File handling routes
	files := e.Group("/files")
	files.Use(echojwt.JWT([]byte(os.Getenv("JWT_SECRET"))))
	files.POST("/upload", handlers.UploadFile)
	files.GET("/download/:id", handlers.DownloadFile)
	files.GET("/list", handlers.ListFiles)
}
