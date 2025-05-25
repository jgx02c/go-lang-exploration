package handlers

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"echo-api/db"
	"echo-api/models"
)

// Register user
func Register(c echo.Context) error {
	u := new(models.User)
	if err := c.Bind(u); err != nil {
		return err
	}

	_, err := db.DB.Exec(context.Background(),
		"INSERT INTO users (username, password) VALUES ($1, $2)", u.Username, u.Password)
	if err != nil {
		log.Printf("register error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "could not create user"})
	}
	return c.JSON(http.StatusCreated, u)
}

// Login and return JWT
func Login(c echo.Context) error {
	u := new(models.User)
	if err := c.Bind(u); err != nil {
		return err
	}

	row := db.DB.QueryRow(context.Background(),
		"SELECT id, password FROM users WHERE username=$1", u.Username)

	var dbID int
	var dbPass string
	err := row.Scan(&dbID, &dbPass)
	if err != nil || dbPass != u.Password {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "invalid credentials"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": dbID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"token": t})
}

// Protected profile route
func Profile(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["user_id"].(float64)

	var username string
	err := db.DB.QueryRow(context.Background(), "SELECT username FROM users WHERE id=$1", int(id)).Scan(&username)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "user not found"})
	}

	return c.JSON(http.StatusOK, map[string]string{"user": username})
} 