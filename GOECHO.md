## Create New Go
mkdir echo-api
cd echo-api
go mod init echo-api

## Install Echo

go get github.com/labstack/echo/v4

## Run the server

go run main.go



## Register User

curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"123"}'


## Login User

curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"123"}'


## Access Profile

curl -H "Authorization: Bearer <your_token>" \
  http://localhost:8080/profile


## Routes

e := echo.New()

// Root routes
e.GET("/", homeHandler)
e.POST("/register", registerHandler)

// Grouped routes (e.g. /api/v1)
api := e.Group("/api/v1")

// Protected routes
secured := e.Group("/api/v1/user")
secured.Use(middleware.JWT([]byte("secret")))

secured.GET("/profile", profileHandler)

## Common Middleware

e.Use(middleware.Logger())
e.Use(middleware.Recover())
e.Use(middleware.CORS()) // optional


## JWT Example

userGroup := e.Group("/user")
userGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
    SigningKey: []byte(os.Getenv("JWT_SECRET")),
}))
userGroup.GET("/profile", profileHandler)


## Reccommended Project Scructure

/myapp
│
├── main.go                # app entry point
├── routes/
│   └── routes.go          # route definitions
├── handlers/
│   └── user.go            # controller functions
├── models/
│   └── user.go            # user struct and db logic
├── middleware/
│   └── auth.go            # custom middleware (e.g., auth)
├── utils/
│   └── jwt.go             # helper to create/parse tokens
├── db/
│   └── connect.go         # db connection
├── go.mod
├── .env
