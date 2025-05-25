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
