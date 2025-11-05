# ZOGTest-Golang

## Requirement
- Postgres for Database

## Library
- [godotenv](github.com/joho/godotenv) for create .ENV
- [Echo](github.com/labstack/echo/v4) for Routing
- [Testify](github.com/stretchr/testify) for Testing
- [Firebase Admin SDK](firebase.google.com/go/v4) for Push Notifications

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

- Run Swagger
```bash
http://localhost:8000/swagger/index.html
```
- Menggunakan Docker
```bash
wsl -d Ubuntu

docker build -t zogtest .
```

- Untuk menjalankan container
```bash
docker run --rm -p 8000:8000 \
  -e APP_HOST=0.0.0.0 \
  -e APP_PORT=8000 \
  -e DATABASE_URL="postgres://postgres:aero1996@host.docker.internal:5432/zogtest-golang" \
  zogtest
```

- Untuk menjalankan unit tests
```bash
go test ./...
```

## Firebase Integration

This application includes Firebase Cloud Messaging (FCM) for push notifications. See [FIREBASE.md](FIREBASE.md) for detailed setup instructions.

**Quick Setup:**
1. Download Firebase service account credentials from Firebase Console
2. Save as `firebase-credentials.json` in project root
3. Add `FIREBASE_CREDENTIALS_PATH=./firebase-credentials.json` to `.env`
4. Restart the application

Firebase endpoints will be available at `/api/v1/notifications/*` if configured.
