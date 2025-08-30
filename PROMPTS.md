# PROMPTS.md

## 1. Swagger dan Github Workflow 
- "Buatkan swagger json dari code berikut
```go
func init() { 
    config.LoadEnv() 
} 
func main() { 
    config.SetupLogging() dbPool, err := database.SetupPgxPool() if err != nil { logging.LogError(context.Background(), err, "database_setup") os.Exit(1) 
}
```
- "membuat github workflow main.yml"
- "go run test error in github action"
- "golang httprouter"
- "cannot use adaptHandler(httpSwager.WrapHandler) (value of func type httprouter.Handle) as echo.HandlerFunc value in argument to e.GET"
- "Refused to load the script 'http://localhost:8000/documentation/swagger-ui-bundle.js' because it violates the following Content Security Policy directive: "default-src 'none'". Note that 'script-src-elem' was not explicitly set, so 'default-src' is used as a fallback."
- "fixing because it violates the following Content Security Policy directive: "default-src 'none'". Note that 'img-src' was not explicitly set, so 'default-src' is used as a fallback"
- "how implement in main.go"
- "Refused to apply inline style because it violates the following Content Security Policy directive: "default-src 'none'". Either the 'unsafe-inline' keyword, a hash ('sha256-ezdv1bOGcoOD7FKudKN0Y2Mb763O6qVtM8LT2mtanIU='), or a nonce ('nonce-...') is required to enable inline execution. Note that hashes do not apply to event handlers, style attributes and javascript: navigations unless the 'unsafe-hashes' keyword is present. Note also that 'style-src' was not explicitly set, so 'default-src' is used as a fallback."
- "Refused to apply inline style because it violates the following Content Security Policy directive: "default-src 'self'". Either the 'unsafe-inline' keyword, a hash ('sha256-ezdv1bOGcoOD7FKudKN0Y2Mb763O6qVtM8LT2mtanIU='), or a nonce ('nonce-...') is required to enable inline execution. Note that hashes do not apply to event handlers, style attributes and javascript: navigations unless the 'unsafe-hashes' keyword is present. Note also that 'style-src' was not explicitly set, so 'default-src' is used as a fallback."

## 2. Swagger dan Github Workflow 
- "buatkan seperti itu tapi untuk post
```bash
// GetNews godoc 
// @Summary List news 
// @Description get string by ID 
// @Tags news
// @Produce json 
// @Param id path int true "Account ID" 
// @Success 200 {array} domain.News 
// @Failure 500 {object} domain.ResponseSingleData[domain.Empty]
// @Security ApiKeyAuth 
// @Router /news/{id} [get]
```
"
- "buatkan dengan isian title , slug ,content"
- "buatkan fungsi untuk update"
- "buatkan untuk delete"
- "create github Creating PostgreSQL service containers"
- "tambah pengaturan postgres
```bash
name: Go CI 
on: push: branches: [ "main" ] 
pull_request: branches: [ "main" ]
...
```
"
- "ini app_test , app_user, app_pass di ambil dari mana ya"
- "how set variable for connection in github workflow"
- "untuk settingan connection girhub workflow postgre yang sesuai dan berikan instruksi yang jelas"
- "kalau di github itu di gitignore .env bagaiamana?"
- "how fix?
```bash
2025-08-30 09:01:09.063 UTC [98] ERROR: relation "topik" does not exist at character 16 2025-08-30 09:01:09.063 UTC [98] STATEMENT: INSERT INTO topik (name, slug, created_at, updated_at) VALUES ($1, $2, NOW(), NOW()) RETURNING id The database cluster will be initialized with locale "en_US.utf8". The default database encoding has accordingly been set to "UTF8". The default text search configuration will be set to "english".
```
"
- "bagaimana saya membuat db/schema.sql?"
- "error while importing github.com/prometheus/client_golang/prometheus: missing go.sum entry for module providing package github.com/beorn7/perks/quantile (imported by github.com/prometheus/client_golang/prometheus); to add:"
- "how install github.com/prometheus/client_golang/prometheus"

## 3. Dockerfile
- "how to run dockerfile"
- "=> ERROR [skeleton 2/2] RUN moon docker scaffold go-app 1.0s"
- "I use WSL to run the docker"
- "yes please give me the fully patcher Dockerfile"
- "please provide the complete Dockerfile with the update"
- "can you give me step by step guide to run moon init on the docker WSL"
- "RUN --mount=type=cache,target=/gomod-cache --mount=type=cache,target=/go-cache go build -o /srv/go-app . how to fixed it?"
- "check wsl is running"
- "The distribution failed to start. Error code: 6, failure step: 2
Error code: Wsl/Service/CreateInstance/E_FAIL"
- "
```bash
Restart-Service : Service 'LxssManager (LxssManager)' cannot be stopped due to the following error: Cannot open LxssManager service on computer '.'. At line:1 char:1 + Restart-Service LxssManager + ~~~~~~~~~~~~~~~~~~~~~~~~~~~ + CategoryInfo : CloseError: (System.ServiceProcess.ServiceController:ServiceController) [Restart-Service ], ServiceCommandException + FullyQualifiedErrorId : CouldNotStopService,Microsoft.PowerShell.Commands.RestartServiceCommand
```
"