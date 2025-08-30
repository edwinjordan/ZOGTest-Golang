# ZOGTest-Golang

## Requirement
- Postgres for Database

## Library
- [godotenv](github.com/joho/godotenv) for create .ENV
- [Echo](github.com/labstack/echo/v4) for Routing
- [Testify](github.com/stretchr/testify) for Testing

## How To Setup on Local Environment
- Git clone this repository to your local environment
- Copy env file to .env
- Change configuration with your own configuration
- Open file main.go
- If you want to run database migration you can remove comment on this line

- go to terminal and run 

```go
go run main.go
```

- Set golang environtment variable

```bash
# GOROOT is the location where Go package is installed on your system
export GOROOT=/usr/lib/go
# GOPATH is the location of your work directory
export GOPATH=$HOME/go
```

- Move to go directory and clone repository from git

```bash
git clone https://github.com/edwinjordan/ZOGTest-Golang.git
cd ZOGTest-Golang
```

- Copy env to .env and change configuration with your own configuration

```bash
cp env .env
```

- Check if the app is running normally

```bash
go run main.go
```

- Build golang app

```bash
go build
```

- Test golang app

```bash
go test ./... -v
```

- Test golang app