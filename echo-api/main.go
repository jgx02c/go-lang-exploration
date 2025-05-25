package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/joho/godotenv"
)

var db *pgxpool.Pool

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"` // plaintext for demo only — hash in prod
}

func main() {
	godotenv.Load()

	// DB setup
	var err error
	db, err = pgxpool.New(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Echo setup
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/register", register)
	e.POST("/login", login)

	// Protected group
	r := e.Group("/profile")
	r.Use(echojwt.JWT([]byte(os.Getenv("JWT_SECRET"))))

	r.GET("", profile)

	// Start
	e.Logger.Fatal(e.Start(":8080"))
}

// Register user
func register(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}

	_, err := db.Exec(context.Background(),
		"INSERT INTO users (username, password) VALUES ($1, $2)", u.Username, u.Password)
	if err != nil {
		log.Printf("register error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "could not create user"})
	}
	return c.JSON(http.StatusCreated, u)
}

// Login and return JWT
func login(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}

	row := db.QueryRow(context.Background(),
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
func profile(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["user_id"].(float64)

	var username string
	err := db.QueryRow(context.Background(), "SELECT username FROM users WHERE id=$1", int(id)).Scan(&username)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "user not found"})
	}

	return c.JSON(http.StatusOK, map[string]string{"user": username})
}
