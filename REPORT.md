# AARCSX Forge Ultimate Repository Report

Generated on: 2026-05-30 16:06:32 UTC

## 1. Scope

- Root analyzed: /home/akshansh/AARCSX_Forge
- Includes: source code, docs, configs, migrations, commands, tests, package files
- Excludes from deep symbol parsing: .git internals and binary artifacts
- Primary language: Go

## 2. Root File Inventory

- .env
- .gitignore
- CLAUDE.md
- OSS_VS_ENTERPRISE.md
- README.md
- REPORT.md
- VERSIONING.md
- cmd/api/main.go
- cmd/forge/create.go
- cmd/forge/doctor.go
- cmd/forge/root.go
- cmd/forge/version.go
- cmd/worker/main.go
- docs/ARCHITECTURE.md
- go.mod
- go.sum
- internal/app/bootstrap.go
- internal/cli/doctor/checks.go
- internal/cli/doctor/runner.go
- internal/cli/runtime/version_manager.go
- internal/cli/scaffold/scaffold.go
- internal/cli/templates/github_releases.go
- internal/cli/templates/provider.go
- internal/config/config.go
- internal/database/postgres.go
- internal/database/redis.go
- internal/health/handler.go
- internal/identity/contracts.go
- internal/identity/dto.go
- internal/identity/entity.go
- internal/identity/handler.go
- internal/identity/permission.go
- internal/identity/repository.go
- internal/identity/role_entity.go
- internal/identity/route.go
- internal/identity/service.go
- internal/logger/logger.go
- internal/notifications/contracts.go
- internal/notifications/handler.go
- internal/notifications/provider/smtp.go
- internal/notifications/repository.go
- internal/notifications/route.go
- internal/notifications/service.go
- internal/observability/contracts.go
- internal/observability/dto.go
- internal/observability/entity.go
- internal/observability/handler.go
- internal/observability/metrics.go
- internal/observability/repository.go
- internal/observability/route.go
- internal/observability/service.go
- internal/observability/tracing.go
- internal/platform/contextx/context.go
- internal/platform/contextx/middleware/context.go
- internal/platform/events/bus.go
- internal/platform/events/contracts.go
- internal/platform/events/events.go
- internal/platform/httpx/response.go
- internal/platform/middleware/auth.go
- internal/platform/middleware/rate_limit.go
- internal/platform/middleware/secure_headers.go
- internal/platform/middleware/tenant.go
- internal/platform/security/argon2.go
- internal/platform/security/contracts.go
- internal/platform/security/jwt.go
- internal/platform/security/password.go
- internal/shared/README.md
- internal/storage/contracts.go
- internal/storage/handler.go
- internal/storage/provider/minio.go
- internal/storage/provider/s3.go
- internal/storage/repository.go
- internal/storage/route.go
- internal/storage/service.go
- internal/storage/service_test.go
- internal/tenants/contracts.go
- internal/tenants/dto.go
- internal/tenants/entity.go
- internal/tenants/handler.go
- internal/tenants/repository.go
- internal/tenants/route.go
- internal/tenants/service.go
- main
- migrations/000001_init_schema.down.sql
- migrations/000001_init_schema.up.sql
- migrations/000002_2.down.sql
- migrations/000002_2.up.sql
- migrations/000003_000003_notification_deliveries.down.sql
- migrations/000003_000003_notification_deliveries.up.sql
- migrations/000004_audit_logs.down.sql
- migrations/000004_audit_logs.up.sql
- pkg/version/version.go

## 3. Module Map (Top-Level)

- .agents: 0 files
- .codex: 0 files
- cmd: 6 files
- cmd/api: 1 files
- cmd/forge: 4 files
- cmd/worker: 1 files
- docs: 1 files
- internal: 66 files
- internal/app: 1 files
- internal/cli: 6 files
- internal/config: 1 files
- internal/database: 2 files
- internal/enterprise: 0 files
- internal/health: 1 files
- internal/identity: 9 files
- internal/logger: 1 files
- internal/notifications: 6 files
- internal/observability: 9 files
- internal/platform: 14 files
- internal/shared: 1 files
- internal/storage: 8 files
- internal/tenants: 7 files
- migrations: 8 files
- pkg: 1 files
- pkg/version: 1 files

## 4. Feature Overview

- API runtime (`cmd/api`) with Gin router and module route registration.
- Worker runtime (`cmd/worker`) with background processing (email task handling).
- CLI (`cmd/forge`) with project creation, version reporting, and environment doctor checks.
- Core dependency bootstrap (`internal/app`) that wires config, logger, DB, redis, metrics, tracing, services.
- Platform primitives (`internal/platform`) for context propagation, middleware, eventing, security, response envelopes.
- Domain modules:
  - `internal/identity`: authentication, user repository/service/handler, RBAC permission checker.
  - `internal/tenants`: tenant lifecycle CRUD.
  - `internal/storage`: signed upload/download URL generation and metadata persistence.
  - `internal/notifications`: email queueing and delivery tracking.
  - `internal/observability`: audit logging, metrics, tracing abstractions and handlers.
- Infra modules:
  - `internal/config`: env-based configuration loading + validation.
  - `internal/database`: postgres + redis connection setup.
  - `internal/logger`: zap logger construction helpers.
- Schema management via SQL migrations (`migrations/`).
- Project docs and product positioning (`README.md`, `docs/ARCHITECTURE.md`, `OSS_VS_ENTERPRISE.md`, `VERSIONING.md`).

## 5. Go Source Deep Dive (Functions, Types, Consts, Vars, Dependencies, Call Sites)

### File: cmd/api/main.go
- Package: package main
- Imports:
  - "context"
  - "log"
  - "net/http"
  - "os"
  - "os/signal"
  - "time"
  - "github.com/AARCSX/AARCSX_Forge/internal/app"
  - "github.com/AARCSX/AARCSX_Forge/internal/config"
  - "github.com/AARCSX/AARCSX_Forge/internal/health"
  - "github.com/AARCSX/AARCSX_Forge/internal/identity"
  - "github.com/AARCSX/AARCSX_Forge/internal/platform/middleware"
  - "github.com/AARCSX/AARCSX_Forge/internal/platform/contextx/middleware"
  - "github.com/AARCSX/AARCSX_Forge/internal/observability"
  - "github.com/AARCSX/AARCSX_Forge/internal/storage"
  - "github.com/AARCSX/AARCSX_Forge/internal/tenants"
  - "github.com/AARCSX/AARCSX_Forge/internal/platform/httpx"
- Declarations:
  - Line 23: `func main() {`
    - Kind: func
    - Symbol: `main`
    - Purpose: Implements behavior in this module; see body in file.
    - Called/Referenced From: No external references found by textual search.

### File: cmd/forge/create.go
- Package: package forge
- Imports:
  - "context"
  - "encoding/json"
  - "fmt"
  - "io"
  - "net/http"
  - "net/url"
  - "os"
  - "os/exec"
  - "path/filepath"
  - "regexp"
  - "strings"
  - "time"
  - "github.com/AARCSX/AARCSX_Forge/internal/logger"
  - "github.com/AARCSX/AARCSX_Forge/internal/config"
  - "github.com/AARCSX/AARCSX_Forge/pkg/version"
  - "github.com/AARCSX/AARCSX_Forge/internal/cli/scaffold"
  - "github.com/AARCSX/AARCSX_Forge/internal/cli/templates"
  - "github.com/spf13/cobra"
  - "github.com/spf13/viper"
  - "gopkg.in/yaml.v3"
- Declarations:
  - Line 28: `func NewCreateCommand() *cobra.Command {`
    - Kind: func
    - Symbol: `NewCreateCommand`
    - Purpose: Implements behavior in this module; see body in file.
    - Called/Referenced From: No external references found by textual search.

### File: cmd/forge/doctor.go
- Package: package forge
- Imports:
  - "context"
  - "fmt"
  - "github.com/AARCSX/AARCSX_Forge/internal/cli/doctor"
- Declarations:
  - Line 11: `func NewDoctorCommand() *cobra.Command {`
    - Kind: func
    - Symbol: `NewDoctorCommand`
    - Purpose: Implements behavior in this module; see body in file.
    - Called/Referenced From: No external references found by textual search.

### File: cmd/forge/root.go
- Package: package forge
- Imports:
  - "context"
  - "os"
  - "os/signal"
  - "syscall"
  - "time"
  - "github.com/AARCSX/AARCSX_Forge/internal/app"
  - "github.com/AARCSX/AARCSX_Forge/internal/config"
  - "github.com/AARCSX/AARCSX_Forge/internal/logger"
  - "github.com/spf13/cobra"
- Declarations:
  - Line 17: `func NewForgeCommand() *cobra.Command {`
    - Kind: func
    - Symbol: `NewForgeCommand`
    - Purpose: Implements behavior in this module; see body in file.
    - Called/Referenced From: No external references found by textual search.

### File: cmd/forge/version.go
- Package: package forge
- Imports:
  - "fmt"
  - "os"
  - "path/filepath"
  - "time"
  - "github.com/AARCSX/AARCSX_Forge/internal/logger"
  - "github.com/AARCSX/AARCSX_Forge/pkg/version"
  - "github.com/spf13/cobra"
  - "gopkg.in/yaml.v3"
- Declarations:
  - Line 16: `func NewVersionCommand() *cobra.Command {`
    - Kind: func
    - Symbol: `NewVersionCommand`
    - Purpose: Implements behavior in this module; see body in file.
    - Called/Referenced From: No external references found by textual search.

### File: cmd/worker/main.go
- Package: package main
- Imports:
  - "context"
  - "encoding/json"
  - "log"
  - "os"
  - "os/signal"
  - "strconv"
  - "time"
  - "github.com/hibiken/asynq"
  - "github.com/AARCSX/AARCSX_Forge/internal/app"
  - "github.com/AARCSX/AARCSX_Forge/internal/config"
  - "github.com/AARCSX/AARCSX_Forge/internal/notifications"
  - "github.com/AARCSX/AARCSX_Forge/internal/notifications/provider"
  - "github.com/AARCSX/AARCSX_Forge/internal/logger"
- Declarations:
  - Line 21: `type EmailTaskPayload struct {`
    - Kind: type
    - Symbol: `EmailTaskPayload`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From: No external references found by textual search.

### File: internal/app/bootstrap.go
- Package: package app
- Imports:
  - "context"
  - "log"
  - "github.com/AARCSX/AARCSX_Forge/internal/config"
  - "github.com/AARCSX/AARCSX_Forge/internal/database"
  - "github.com/AARCSX/AARCSX_Forge/internal/logger"
  - "github.com/AARCSX/AARCSX_Forge/internal/observability"
  - "github.com/AARCSX/AARCSX_Forge/internal/platform/events"
- Declarations:
  - Line 15: `type RuntimeDeps struct {`
    - Kind: type
    - Symbol: `RuntimeDeps`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From:
      - 1:27:func BuildRuntimeDeps(cfg config.Config) RuntimeDeps {

### File: internal/cli/doctor/checks.go
- Package: package doctor
- Imports:
  - "context"
  - "fmt"
  - "os"
  - "os/exec"
  - "path/filepath"
  - "strings"
  - "time"
  - "github.com/AARCSX/AARCSX_Forge/internal/config"
  - "github.com/AARCSX/AARCSX_Forge/internal/logger"
  - "gopkg.in/yaml.v3"
- Declarations:
  - Line 18: `type CheckResult struct {`
    - Kind: type
    - Symbol: `CheckResult`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From:
      - 6:47:func (c *GoCheck) Check(ctx context.Context) CheckResult {
      - 9:92:func (c *DockerCheck) Check(ctx context.Context) CheckResult {
      - 12:141:func (c *PostgreSQLCheck) Check(ctx context.Context) CheckResult {
      - 15:194:func (c *RedisCheck) Check(ctx context.Context) CheckResult {
      - 18:246:func (c *ForgeMetadataCheck) Check(ctx context.Context) CheckResult {

### File: internal/cli/doctor/runner.go
- Package: package doctor
- Imports:
  - "context"
  - "fmt"
- Declarations:
  - Line 9: `type Runner struct {`
    - Kind: type
    - Symbol: `Runner`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From:
      - 1:14:func NewRunner(checks ...Check) *Runner {
      - 2:21:func (r *Runner) Run(ctx context.Context) []CheckResult {

### File: internal/cli/runtime/version_manager.go
- Package: package runtime
- Imports:
  - "fmt"
  - "github.com/AARCSX/AARCSX_Forge/pkg/version"
- Declarations:
  - Line 10: `type VersionManager struct {`
    - Kind: type
    - Symbol: `VersionManager`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From:
      - 1:16:func NewVersionManager() *VersionManager {
      - 2:24:func (vm *VersionManager) CLIVersion() string {
      - 3:29:func (vm *VersionManager) RuntimeVersion() string {
      - 4:34:func (vm *VersionManager) String() string {
      - 5:40:func (vm *VersionManager) CheckCompatibility() error {

### File: internal/cli/scaffold/scaffold.go
- Package: package scaffold
- Imports:
  - "fmt"
  - "os"
  - "path/filepath"
  - "strings"
- Declarations:
  - Line 11: `func ReplacePlaceholders(dir string, replacements map[string]string) error {`
    - Kind: func
    - Symbol: `ReplacePlaceholders`
    - Purpose: Implements behavior in this module; see body in file.
    - Called/Referenced From: No external references found by textual search.

### File: internal/cli/templates/github_releases.go
- Package: package templates
- Imports:
  - "encoding/json"
  - "fmt"
  - "io"
  - "net/http"
  - "os"
  - "path/filepath"
  - "time"
  - "github.com/AARCSX/AARCSX_Forge/internal/logger"
- Declarations:
  - Line 16: `type GitHubReleasesProvider struct {`
    - Kind: type
    - Symbol: `GitHubReleasesProvider`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From:
      - 1:23:func NewGitHubReleasesProvider(owner, repo string) *GitHubReleasesProvider {
      - 2:34:func (p *GitHubReleasesProvider) GetTemplate(version string) ([]byte, error) {
      - 3:40:func (p *GitHubReleasesProvider) DownloadTemplate(releaseTag string) (string, error) {
      - 4:144:func (p *GitHubReleasesProvider) ExtractTemplate(archivePath, destDir string) error {

### File: internal/cli/templates/provider.go
- Package: package templates
- Imports:
  - "errors"
  - "fmt"
- Declarations:
  - Line 9: `type TemplateProvider interface {`
    - Kind: type
    - Symbol: `TemplateProvider`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From: No external references found by textual search.

### File: internal/config/config.go
- Package: package config
- Imports:
  - "errors"
  - "fmt"
  - "os"
  - "strings"
- Declarations:
  - Line 10: `type Environment string`
    - Kind: type
    - Symbol: `Environment`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From: No external references found by textual search.

### File: internal/database/postgres.go
- Package: package database
- Imports:
  - "context"
  - "database/sql"
  - "fmt"
  - "time"
  - _ "github.com/jackc/pgx/v5/stdlib"
  - "github.com/AARCSX/AARCSX_Forge/internal/config"
- Declarations:
  - Line 14: `type Postgres struct {`
    - Kind: type
    - Symbol: `Postgres`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From:
      - 1:19:func NewPostgres(cfg config.Config) (*Postgres, error) {
      - 2:42:func (p *Postgres) Close() error {

### File: internal/database/redis.go
- Package: package database
- Imports:
  - "context"
  - "fmt"
  - "time"
  - "github.com/redis/go-redis/v9"
  - "github.com/AARCSX/AARCSX_Forge/internal/config"
- Declarations:
  - Line 13: `type Redis struct {`
    - Kind: type
    - Symbol: `Redis`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From:
      - 1:18:func NewRedis(cfg config.Config) (*Redis, error) {
      - 2:37:func (r *Redis) Close() error {

### File: internal/health/handler.go
- Package: package health
- Imports:
  - "net/http"
  - "github.com/gin-gonic/gin"
  - "github.com/AARCSX/AARCSX_Forge/internal/app"
- Declarations:
  - Line 11: `type Handler struct {`
    - Kind: type
    - Symbol: `Handler`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From:
      - 1:16:func NewHandler(deps *app.RuntimeDeps) *Handler {
      - 3:32:func (h *Handler) Check(c *gin.Context) {
      - 4:59:func (h *Handler) Ready(c *gin.Context) {
      - 5:64:func (h *Handler) Live(c *gin.Context) {

### File: internal/identity/contracts.go
- Package: package identity
- Imports:
  - "context"
  - "github.com/google/uuid"
  - "errors"
- Declarations:
  - Line 11: `var (`
    - Kind: var
    - Symbol: `(`
    - Purpose: Variable definition for shared state/errors/config defaults.
    - Called/Referenced From: No external references found by textual search.
  - Line 18: `type Service interface {`
    - Kind: type
    - Symbol: `Service`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From: No external references found by textual search.

### File: internal/identity/dto.go
- Package: package identity
- Imports:
- Declarations:
  - Line 4: `type RegisterUserInput struct {`
    - Kind: type
    - Symbol: `RegisterUserInput`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From: No external references found by textual search.

### File: internal/identity/entity.go
- Package: package identity
- Imports:
  - "time"
  - "github.com/google/uuid"
  - "golang.org/x/crypto/bcrypt"
- Declarations:
  - Line 11: `type User struct {`
    - Kind: type
    - Symbol: `User`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From:
      - 1:26:func (u *User) SetPassword(password string) error {
      - 2:36:func (u *User) CheckPassword(password string) error {

### File: internal/identity/handler.go
- Package: package identity
- Imports:
  - "net/http"
  - "github.com/gin-gonic/gin"
  - "github.com/AARCSX/AARCSX_Forge/internal/logger"
  - "github.com/AARCSX/AARCSX_Forge/internal/platform/httpx"
- Declarations:
  - Line 12: `type AuthHandler struct {`
    - Kind: type
    - Symbol: `AuthHandler`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From:
      - 1:18:func NewAuthHandler(service Service, logger *logger.Logger) *AuthHandler {
      - 2:23:func (h *AuthHandler) Register(c *gin.Context) {
      - 3:55:func (h *AuthHandler) Login(c *gin.Context) {
      - 4:87:func (h *AuthHandler) Refresh(c *gin.Context) {

### File: internal/identity/permission.go
- Package: package identity
- Imports:
  - "context"
  - "github.com/AARCSX/AARCSX_Forge/internal/database"
  - "github.com/google/uuid"
- Declarations:
  - Line 11: `type RBACPermissionChecker struct {`
    - Kind: type
    - Symbol: `RBACPermissionChecker`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From:
      - 1:16:func NewRBACPermissionChecker(db *database.Postgres) *RBACPermissionChecker {
      - 2:21:func (c *RBACPermissionChecker) Can(ctx context.Context, actorID, tenantID string, action string) (bool, error) {

### File: internal/identity/repository.go
- Package: package identity
- Imports:
  - "context"
  - "database/sql"
  - "errors"
  - "time"
  - "github.com/google/uuid"
  - "github.com/AARCSX/AARCSX_Forge/internal/database"
- Declarations:
  - Line 14: `type IdentityRepository interface {`
    - Kind: type
    - Symbol: `IdentityRepository`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From:
      - 2:27:func NewIdentityRepositoryPostgres(db *database.Postgres) IdentityRepository {

### File: internal/identity/role_entity.go
- Package: package identity
- Imports:
  - "time"
  - "github.com/google/uuid"
- Declarations:
  - Line 10: `type Role struct {`
    - Kind: type
    - Symbol: `Role`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From: No external references found by textual search.

### File: internal/identity/route.go
- Package: package identity
- Imports:
  - "github.com/gin-gonic/gin"
  - "github.com/AARCSX/AARCSX_Forge/internal/app"
- Declarations:
  - Line 9: `func RegisterRoutes(router *gin.RouterGroup, deps *app.RuntimeDeps) {`
    - Kind: func
    - Symbol: `RegisterRoutes`
    - Purpose: Implements behavior in this module; see body in file.
    - Called/Referenced From: No external references found by textual search.

### File: internal/identity/service.go
- Package: package identity
- Imports:
  - "context"
  - "crypto/sha256"
  - "encoding/hex"
  - "errors"
  - "time"
  - "github.com/golang-jwt/jwt/v5"
  - "github.com/google/uuid"
  - "github.com/AARCSX/AARCSX_Forge/internal/config"
  - "github.com/AARCSX/AARCSX_Forge/internal/database"
  - "github.com/AARCSX/AARCSX_Forge/internal/logger"
  - "github.com/AARCSX/AARCSX_Forge/internal/platform/security"
- Declarations:
  - Line 19: `type IdentityService struct {`
    - Kind: type
    - Symbol: `IdentityService`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From:
      - 1:29:func NewIdentityService(repo Repository, cfg config.Config, logger *logger.Logger, permissionChecker PermissionChecker, db *database.Postgres) *IdentityService {
      - 2:49:func (s *IdentityService) RegisterUser(ctx context.Context, in *RegisterUserInput) (*UserResponse, error) {
      - 3:92:func (s *IdentityService) Authenticate(ctx context.Context, in *AuthenticateInput) (*TokenPairResponse, error) {
      - 4:139:func (s *IdentityService) RefreshToken(ctx context.Context, in *RefreshTokenInput) (*TokenPairResponse, error) {
      - 5:204:func (s *IdentityService) generateAccessToken(u *User) (string, error) {
      - 6:216:func (s *IdentityService) generateRefreshToken(u *User) (string, error) {

### File: internal/logger/logger.go
- Package: package logger
- Imports:
  - "go.uber.org/zap"
  - "go.uber.org/zap/zapcore"
- Declarations:
  - Line 9: `type Logger struct {`
    - Kind: type
    - Symbol: `Logger`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From:
      - 1:16:func New(level string, serviceName string) (*Logger, error) {
      - 2:54:func (l *Logger) Sugar() *zap.SugaredLogger {
      - 3:60:func (l *Logger) With(fields ...zap.Field) *Logger {
      - 4:65:func (l *Logger) SugarWith(fields ...zap.Field) *zap.SugaredLogger {

### File: internal/notifications/contracts.go
- Package: package notifications
- Imports:
  - "context"
  - "time"
- Declarations:
  - Line 8: `type Service interface {`
    - Kind: type
    - Symbol: `Service`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From: No external references found by textual search.

### File: internal/notifications/handler.go
- Package: package notifications
- Imports:
  - "net/http"
  - "github.com/gin-gonic/gin"
  - "github.com/AARCSX/AARCSX_Forge/internal/logger"
  - "github.com/AARCSX/AARCSX_Forge/internal/platform/httpx"
- Declarations:
  - Line 12: `type NotificationHandler struct {`
    - Kind: type
    - Symbol: `NotificationHandler`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From:
      - 1:18:func NewNotificationHandler(service Service, logger *logger.Logger) *NotificationHandler {
      - 2:23:func (h *NotificationHandler) SendEmail(c *gin.Context) {

### File: internal/notifications/provider/smtp.go
- Package: package provider
- Imports:
  - "context"
  - "fmt"
  - "net/smtp"
  - "github.com/AARCSX/AARCSX_Forge/internal/notifications"
- Declarations:
  - Line 12: `type SMTPProvider struct {`
    - Kind: type
    - Symbol: `SMTPProvider`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From:
      - 1:21:func NewSMTPProvider(host string, port int, username, password, from string) (*SMTPProvider, error) {
      - 2:42:func (p *SMTPProvider) SendEmail(ctx context.Context, in notifications.EmailMessage) error {

### File: internal/notifications/repository.go
- Package: package notifications
- Imports:
  - "context"
  - "time"
  - "github.com/google/uuid"
  - "github.com/AARCSX/AARCSX_Forge/internal/database"
- Declarations:
  - Line 12: `type DeliveryJobPostgres struct {`
    - Kind: type
    - Symbol: `DeliveryJobPostgres`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From:
      - 2:22:func (r *DeliveryJobPostgres) CreateDelivery(ctx context.Context, in DeliveryJob) (DeliveryJob, error) {
      - 3:44:func (r *DeliveryJobPostgres) MarkDelivered(ctx context.Context, jobID string) error {
      - 4:56:func (r *DeliveryJobPostgres) MarkFailed(ctx context.Context, jobID string, reason string) error {

### File: internal/notifications/route.go
- Package: package notifications
- Imports:
  - "fmt"
  - "github.com/gin-gonic/gin"
  - "github.com/AARCSX/AARCSX_Forge/internal/app"
- Declarations:
  - Line 12: `func RegisterRoutes(router *gin.RouterGroup, deps *app.RuntimeDeps) {`
    - Kind: func
    - Symbol: `RegisterRoutes`
    - Purpose: Implements behavior in this module; see body in file.
    - Called/Referenced From: No external references found by textual search.

### File: internal/notifications/service.go
- Package: package notifications
- Imports:
  - "context"
  - "errors"
  - "fmt"
  - "github.com/AARCSX/AARCSX_Forge/internal/logger"
- Declarations:
  - Line 12: `type NotificationService struct {`
    - Kind: type
    - Symbol: `NotificationService`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From:
      - 1:19:func NewNotificationService(logger *logger.Logger) (*NotificationService, error) {
      - 2:26:func (s *NotificationService) WithProvider(p Provider) *NotificationService {
      - 3:32:func (s *NotificationService) WithRepository(r Repository) *NotificationService {
      - 4:38:func (s *NotificationService) QueueEmail(ctx context.Context, in QueueEmailInput) (DeliveryJob, error) {

### File: internal/observability/contracts.go
- Package: package observability
- Imports:
  - "context"
- Declarations:
  - Line 5: `type Logger interface {`
    - Kind: type
    - Symbol: `Logger`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From: No external references found by textual search.

### File: internal/observability/dto.go
- Package: package observability
- Imports:
- Declarations:
  - Line 4: `type AuditLogCreateDTO struct {`
    - Kind: type
    - Symbol: `AuditLogCreateDTO`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From: No external references found by textual search.

### File: internal/observability/entity.go
- Package: package observability
- Imports:
- Declarations:
  - Line 4: `type AuditLog struct {`
    - Kind: type
    - Symbol: `AuditLog`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From: No external references found by textual search.

### File: internal/observability/handler.go
- Package: package observability
- Imports:
  - "context"
  - "net/http"
  - "strconv"
  - "github.com/gin-gonic/gin"
  - "github.com/AARCSX/AARCSX_Forge/internal/logger"
  - "github.com/AARCSX/AARCSX_Forge/internal/platform/httpx"
- Declarations:
  - Line 15: `type AuditLogHandler struct {`
    - Kind: type
    - Symbol: `AuditLogHandler`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From:
      - 2:30:func NewAuditLogHandler(service Service, logger *logger.Logger, metrics Metrics, tracer *Tracer) *AuditLogHandler {
      - 3:35:func (h *AuditLogHandler) CreateAuditLog(c *gin.Context) {
      - 4:87:func (h *AuditLogHandler) GetAuditLog(c *gin.Context) {
      - 5:152:func (h *AuditLogHandler) ListAuditLogs(c *gin.Context) {
      - 6:225:func (h *AuditLogHandler) handleServiceError(c *gin.Context, err error) {

### File: internal/observability/metrics.go
- Package: package observability
- Imports:
  - "sync"
  - "time"
- Declarations:
  - Line 10: `type MetricsCollector struct {`
    - Kind: type
    - Symbol: `MetricsCollector`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From:
      - 1:21:func NewMetricsCollector() *MetricsCollector {
      - 2:30:func (m *MetricsCollector) IncCounter(name string, labels map[string]string) {
      - 3:38:func (m *MetricsCollector) AddCounter(name string, value int64, labels map[string]string) {
      - 4:46:func (m *MetricsCollector) ObserveHistogram(name string, value float64, labels map[string]string) {
      - 5:54:func (m *MetricsCollector) SetGauge(name string, value float64, labels map[string]string) {
      - 6:62:func (m *MetricsCollector) GetCounter(name string, labels map[string]string) int64 {
      - 7:71:func (m *MetricsCollector) GetHistogram(name string, labels map[string]string) []float64 {
      - 8:81:func (m *MetricsCollector) GetGauge(name string, labels map[string]string) float64 {
      - ... plus 2 additional references

### File: internal/observability/repository.go
- Package: package observability
- Imports:
  - "context"
  - "database/sql"
  - "encoding/json"
  - "fmt"
  - "strings"
  - "time"
  - "github.com/AARCSX/AARCSX_Forge/internal/database"
  - "github.com/AARCSX/AARCSX_Forge/internal/logger"
  - "go.uber.org/zap"
- Declarations:
  - Line 18: `type AuditLogRepository struct {`
    - Kind: type
    - Symbol: `AuditLogRepository`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From:
      - 1:24:func NewAuditLogRepository(db *database.Postgres, logger logger.Logger) *AuditLogRepository {
      - 2:32:func (r *AuditLogRepository) Create(ctx context.Context, log AuditLog) (int64, error) {
      - 3:66:func (r *AuditLogRepository) GetByID(ctx context.Context, id int64) (*AuditLog, error) {
      - 4:100:func (r *AuditLogRepository) List(ctx context.Context, tenantID int64, userID int64,

### File: internal/observability/route.go
- Package: package observability
- Imports:
  - "github.com/gin-gonic/gin"
  - "github.com/AARCSX/AARCSX_Forge/internal/database"
  - "github.com/AARCSX/AARCSX_Forge/internal/logger"
- Declarations:
  - Line 11: `func RegisterRoutes(router *gin.RouterGroup, postgres *database.Postgres, logger *logger.Logger, metrics Metrics, tracer *Tracer) {`
    - Kind: func
    - Symbol: `RegisterRoutes`
    - Purpose: Implements behavior in this module; see body in file.
    - Called/Referenced From: No external references found by textual search.

### File: internal/observability/service.go
- Package: package observability
- Imports:
  - "context"
  - "encoding/json"
  - "time"
  - "github.com/AARCSX/AARCSX_Forge/internal/logger"
- Declarations:
  - Line 13: `type AuditLogService struct {`
    - Kind: type
    - Symbol: `AuditLogService`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From:
      - 2:29:func NewAuditLogService(repo Repository, logger logger.Logger, auditLogger Logger, metrics Metrics, tracer *Tracer) *AuditLogService {
      - 3:40:func (s *AuditLogService) CreateAuditLog(ctx context.Context, input AuditLogCreateDTO) (AuditLogResponse, error) {
      - 4:147:func (s *AuditLogService) GetAuditLog(ctx context.Context, id int64) (*AuditLogResponse, error) {
      - 5:203:func (s *AuditLogService) ListAuditLogs(ctx context.Context, tenantID int64, userID int64,

### File: internal/observability/tracing.go
- Package: package observability
- Imports:
  - "context"
  - "time"
  - "github.com/AARCSX/AARCSX_Forge/internal/platform/contextx"
- Declarations:
  - Line 11: `type TraceIDGenerator interface {`
    - Kind: type
    - Symbol: `TraceIDGenerator`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From:
      - 9:69:func NewTracer(traceIDGenerator TraceIDGenerator) *Tracer {

### File: internal/platform/contextx/context.go
- Package: package contextx
- Imports:
  - "context"
- Declarations:
  - Line 5: `type key string`
    - Kind: type
    - Symbol: `key`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From: No external references found by textual search.

### File: internal/platform/contextx/middleware/context.go
- Package: package middleware
- Imports:
  - "net/http"
  - "time"
  - "github.com/gin-gonic/gin"
  - "github.com/google/uuid"
  - "github.com/AARCSX/AARCSX_Forge/internal/app"
  - "github.com/AARCSX/AARCSX_Forge/internal/logger"
  - "github.com/AARCSX/AARCSX_Forge/internal/observability"
  - "github.com/AARCSX/AARCSX_Forge/internal/platform/contextx"
- Declarations:
  - Line 19: `func ContextMiddleware(deps *app.RuntimeDeps) gin.HandlerFunc {`
    - Kind: func
    - Symbol: `ContextMiddleware`
    - Purpose: Implements behavior in this module; see body in file.
    - Called/Referenced From: No external references found by textual search.

### File: internal/platform/events/bus.go
- Package: package events
- Imports:
  - "sync"
- Declarations:
  - Line 8: `type EventHandler func(event interface{})`
    - Kind: type
    - Symbol: `EventHandler`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From:
      - 4:47:func (b *Bus) Subscribe(eventType string, handler EventHandler) {
      - 5:54:func (b *Bus) SubscribeOnce(eventType string, handler EventHandler) {
      - 6:63:func (b *Bus) Unsubscribe(eventType string, handler EventHandler) {

### File: internal/platform/events/contracts.go
- Package: package events
- Imports:
  - "time"
- Declarations:
  - Line 5: `type BaseEvent struct {`
    - Kind: type
    - Symbol: `BaseEvent`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From:
      - 1:9:func (e BaseEvent) OccurredAt() time.Time { return e.Occurred }

### File: internal/platform/events/events.go
- Package: package events
- Imports:
  - "context"
  - "fmt"
  - "sync"
  - "time"
- Declarations:
  - Line 10: `type Event interface {`
    - Kind: type
    - Symbol: `Event`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From:
      - 1:16:type Handler func(context.Context, Event) error
      - 6:40:func (b *InMemoryBus) Publish(ctx context.Context, event Event) error {

### File: internal/platform/httpx/response.go
- Package: package httpx
- Imports:
  - "net/http"
  - "github.com/gin-gonic/gin"
  - "github.com/AARCSX/AARCSX_Forge/internal/platform/contextx"
- Declarations:
  - Line 11: `type ResponseEnvelope struct {`
    - Kind: type
    - Symbol: `ResponseEnvelope`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From: No external references found by textual search.

### File: internal/platform/middleware/auth.go
- Package: package middleware
- Imports:
  - "context"
  - "net/http"
  - "strings"
  - "github.com/gin-gonic/gin"
  - "github.com/golang-jwt/jwt/v5"
  - "github.com/google/uuid"
  - "github.com/AARCSX/AARCSX_Forge/internal/app"
  - "github.com/AARCSX/AARCSX_Forge/internal/database"
  - "github.com/AARCSX/AARCSX_Forge/internal/identity"
  - "github.com/AARCSX/AARCSX_Forge/internal/platform/contextx"
- Declarations:
  - Line 20: `func AuthMiddleware(deps *app.RuntimeDeps) gin.HandlerFunc {`
    - Kind: func
    - Symbol: `AuthMiddleware`
    - Purpose: Implements behavior in this module; see body in file.
    - Called/Referenced From: No external references found by textual search.

### File: internal/platform/middleware/rate_limit.go
- Package: package middleware
- Imports:
  - "net/http"
  - "sync"
  - "time"
  - "github.com/gin-gonic/gin"
  - "golang.org/x/time/rate"
- Declarations:
  - Line 13: `func RateLimitMiddleware(r rate.Limit, b int) gin.HandlerFunc {`
    - Kind: func
    - Symbol: `RateLimitMiddleware`
    - Purpose: Implements behavior in this module; see body in file.
    - Called/Referenced From: No external references found by textual search.

### File: internal/platform/middleware/secure_headers.go
- Package: package middleware
- Imports:
  - "net/http"
  - "github.com/gin-gonic/gin"
- Declarations:
  - Line 10: `func SecureHeadersMiddleware() gin.HandlerFunc {`
    - Kind: func
    - Symbol: `SecureHeadersMiddleware`
    - Purpose: Implements behavior in this module; see body in file.
    - Called/Referenced From: No external references found by textual search.

### File: internal/platform/middleware/tenant.go
- Package: package middleware
- Imports:
  - "net/http"
  - "github.com/gin-gonic/gin"
  - "github.com/google/uuid"
  - "github.com/AARCSX/AARCSX_Forge/internal/platform/contextx"
  - tenants "github.com/AARCSX/AARCSX_Forge/internal/tenants"
- Declarations:
  - Line 13: `func TenantMiddleware(ts *tenants.TenantService) gin.HandlerFunc {`
    - Kind: func
    - Symbol: `TenantMiddleware`
    - Purpose: Implements behavior in this module; see body in file.
    - Called/Referenced From: No external references found by textual search.

### File: internal/platform/security/argon2.go
- Package: package security
- Imports:
  - "encoding/base64"
  - "errors"
  - "fmt"
  - "math/rand"
  - "strconv"
  - "strings"
  - "sync"
  - "golang.org/x/crypto/argon2"
  - "golang.org/x/crypto/subtle"
- Declarations:
  - Line 17: `type Argon2idHasher struct {`
    - Kind: type
    - Symbol: `Argon2idHasher`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From:
      - 1:33:func NewArgon2idHasher(memory uint32, iterations uint32, parallelism uint8, saltLength int, keyLength uint32) *Argon2idHasher {
      - 2:47:func (h *Argon2idHasher) Hash(password string) (string, error) {
      - 3:70:func (h *Argon2idHasher) Compare(hashPassword, password string) error {

### File: internal/platform/security/contracts.go
- Package: package security
- Imports:
  - "context"
- Declarations:
  - Line 5: `type PasswordHasher interface {`
    - Kind: type
    - Symbol: `PasswordHasher`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From: No external references found by textual search.

### File: internal/platform/security/jwt.go
- Package: package security
- Imports:
  - "errors"
  - "time"
  - "github.com/golang-jwt/jwt/v5"
- Declarations:
  - Line 11: `type JWTKeyProvider interface {`
    - Kind: type
    - Symbol: `JWTKeyProvider`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From:
      - 5:41:func NewJWTTokenService(keyProvider JWTKeyProvider, cfg JWTConfig) *JWTTokenService {

### File: internal/platform/security/password.go
- Package: package security
- Imports:
  - "golang.org/x/crypto/bcrypt"
- Declarations:
  - Line 8: `type BcryptHasher struct {`
    - Kind: type
    - Symbol: `BcryptHasher`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From:
      - 1:14:func NewBcryptHasher(cost int) *BcryptHasher {
      - 2:22:func (h *BcryptHasher) Hash(password string) (string, error) {
      - 3:31:func (h *BcryptHasher) Compare(hash, password string) error {

### File: internal/storage/contracts.go
- Package: package storage
- Imports:
  - "context"
  - "errors"
  - "time"
- Declarations:
  - Line 10: `type Service interface {`
    - Kind: type
    - Symbol: `Service`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From: No external references found by textual search.

### File: internal/storage/handler.go
- Package: package storage
- Imports:
  - "net/http"
  - "github.com/gin-gonic/gin"
  - "github.com/AARCSX/AARCSX_Forge/internal/logger"
  - "github.com/AARCSX/AARCSX_Forge/internal/platform/httpx"
- Declarations:
  - Line 12: `type StorageHandler struct {`
    - Kind: type
    - Symbol: `StorageHandler`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From:
      - 1:18:func NewStorageHandler(service Service, logger *logger.Logger) *StorageHandler {
      - 2:23:func (h *StorageHandler) GetUploadURL(c *gin.Context) {
      - 3:48:func (h *StorageHandler) GetDownloadURL(c *gin.Context) {

### File: internal/storage/provider/minio.go
- Package: package provider
- Imports:
  - "context"
  - "fmt"
  - "time"
  - "github.com/minio/minio-go/v7"
  - "github.com/minio/minio-go/v7/pkg/credentials"
- Declarations:
  - Line 13: `type MinIOProvider struct {`
    - Kind: type
    - Symbol: `MinIOProvider`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From:
      - 1:20:func NewMinIOProvider(endpoint, accessKey, secretKey string, bucket string, useSSL bool) (*MinIOProvider, error) {
      - 2:50:func (p *MinIOProvider) PresignPut(ctx context.Context, objectKey string, ttlSeconds int) (string, error) {
      - 3:69:func (p *MinIOProvider) PresignGet(ctx context.Context, objectKey string, ttlSeconds int) (string, error) {

### File: internal/storage/provider/s3.go
- Package: package provider
- Imports:
  - "context"
  - "fmt"
- Declarations:
  - Line 10: `type S3Provider struct {`
    - Kind: type
    - Symbol: `S3Provider`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From:
      - 1:15:func NewS3Provider(endpoint, accessKey, secretKey string, bucket string, region string, useSSL bool) (*S3Provider, error) {
      - 2:21:func (p *S3Provider) PresignPut(ctx context.Context, objectKey string, ttlSeconds int) (string, error) {
      - 3:26:func (p *S3Provider) PresignGet(ctx context.Context, objectKey string, ttlSeconds int) (string, error) {

### File: internal/storage/repository.go
- Package: package storage
- Imports:
  - "context"
  - "database/sql"
  - "errors"
  - "time"
  - "github.com/google/uuid"
  - "github.com/AARCSX/AARCSX_Forge/internal/database"
- Declarations:
  - Line 14: `type ObjectRepositoryPostgres struct {`
    - Kind: type
    - Symbol: `ObjectRepositoryPostgres`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From:
      - 2:24:func (r *ObjectRepositoryPostgres) SaveObjectMeta(ctx context.Context, obj ObjectMeta) (ObjectMeta, error) {
      - 3:50:func (r *ObjectRepositoryPostgres) GetObjectMeta(ctx context.Context, tenantID, objectID string) (ObjectMeta, error) {

### File: internal/storage/route.go
- Package: package storage
- Imports:
  - "fmt"
  - "github.com/gin-gonic/gin"
  - "github.com/AARCSX/AARCSX_Forge/internal/app"
- Declarations:
  - Line 11: `func RegisterRoutes(router *gin.RouterGroup, deps *app.RuntimeDeps) {`
    - Kind: func
    - Symbol: `RegisterRoutes`
    - Purpose: Implements behavior in this module; see body in file.
    - Called/Referenced From: No external references found by textual search.

### File: internal/storage/service.go
- Package: package storage
- Imports:
  - "context"
  - "errors"
  - "fmt"
  - "github.com/AARCSX/AARCSX_Forge/internal/config"
  - "github.com/AARCSX/AARCSX_Forge/internal/logger"
  - "github.com/AARCSX/AARCSX_Forge/internal/storage/provider"
- Declarations:
  - Line 14: `type StorageService struct {`
    - Kind: type
    - Symbol: `StorageService`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From:
      - 1:22:func NewStorageService(cfg *config.Config, logger *logger.Logger) (*StorageService, error) {
      - 2:53:func (s *StorageService) WithRepository(repo Repository) *StorageService {
      - 3:59:func (s *StorageService) CreateUploadURL(ctx context.Context, in CreateUploadURLInput) (SignedURL, error) {
      - 4:103:func (s *StorageService) CreateDownloadURL(ctx context.Context, in CreateDownloadURLInput) (SignedURL, error) {

### File: internal/storage/service_test.go
- Package: package storage
- Imports:
  - "context"
  - "testing"
  - "github.com/AARCSX/AARCSX_Forge/internal/logger"
  - "github.com/stretchr/testify/require"
- Declarations:
  - Line 12: `type MockProvider struct {`
    - Kind: type
    - Symbol: `MockProvider`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From:
      - 1:17:func (m *MockProvider) PresignPut(ctx context.Context, objectKey string, ttlSeconds int) (string, error) {
      - 2:24:func (m *MockProvider) PresignGet(ctx context.Context, objectKey string, ttlSeconds int) (string, error) {

### File: internal/tenants/contracts.go
- Package: package tenants
- Imports:
  - "context"
  - "errors"
  - "github.com/google/uuid"
- Declarations:
  - Line 11: `type Repository interface {`
    - Kind: type
    - Symbol: `Repository`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From: No external references found by textual search.

### File: internal/tenants/dto.go
- Package: package tenants
- Imports:
- Declarations:
  - Line 4: `type CreateTenantRequest struct {`
    - Kind: type
    - Symbol: `CreateTenantRequest`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From: No external references found by textual search.

### File: internal/tenants/entity.go
- Package: package tenants
- Imports:
  - "time"
  - "github.com/google/uuid"
- Declarations:
  - Line 10: `type Tenant struct {`
    - Kind: type
    - Symbol: `Tenant`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From: No external references found by textual search.

### File: internal/tenants/handler.go
- Package: package tenants
- Imports:
  - "net/http"
  - "github.com/gin-gonic/gin"
  - "github.com/google/uuid"
  - "github.com/AARCSX/AARCSX_Forge/internal/logger"
  - "github.com/AARCSX/AARCSX_Forge/internal/platform/httpx"
- Declarations:
  - Line 13: `type TenantHandler struct {`
    - Kind: type
    - Symbol: `TenantHandler`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From:
      - 1:19:func NewTenantHandler(service *TenantService, logger *logger.Logger) *TenantHandler {
      - 2:24:func (h *TenantHandler) CreateTenant(c *gin.Context) {
      - 3:49:func (h *TenantHandler) GetTenant(c *gin.Context) {
      - 4:82:func (h *TenantHandler) UpdateTenant(c *gin.Context) {

### File: internal/tenants/repository.go
- Package: package tenants
- Imports:
  - "context"
  - "database/sql"
  - "errors"
  - "time"
  - "github.com/google/uuid"
  - "github.com/AARCSX/AARCSX_Forge/internal/database"
- Declarations:
  - Line 15: `type TenantRepositoryPostgres struct {`
    - Kind: type
    - Symbol: `TenantRepositoryPostgres`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From:
      - 2:25:func (r *TenantRepositoryPostgres) Create(ctx context.Context, t *Tenant) error {
      - 3:38:func (r *TenantRepositoryPostgres) GetByID(ctx context.Context, id uuid.UUID) (*Tenant, error) {
      - 4:58:func (r *TenantRepositoryPostgres) GetBySlug(ctx context.Context, slug string) (*Tenant, error) {
      - 5:78:func (r *TenantRepositoryPostgres) Update(ctx context.Context, t *Tenant) error {

### File: internal/tenants/route.go
- Package: package tenants
- Imports:
  - "github.com/gin-gonic/gin"
  - "github.com/AARCSX/AARCSX_Forge/internal/app"
- Declarations:
  - Line 9: `func RegisterRoutes(router *gin.RouterGroup, deps *app.RuntimeDeps) {`
    - Kind: func
    - Symbol: `RegisterRoutes`
    - Purpose: Implements behavior in this module; see body in file.
    - Called/Referenced From: No external references found by textual search.

### File: internal/tenants/service.go
- Package: package tenants
- Imports:
  - "context"
  - "time"
  - "github.com/google/uuid"
  - "github.com/AARCSX/AARCSX_Forge/internal/logger"
- Declarations:
  - Line 12: `type TenantService struct {`
    - Kind: type
    - Symbol: `TenantService`
    - Purpose: Declares structural contract/data model used by module logic.
    - Called/Referenced From:
      - 1:18:func NewTenantService(repo Repository, logger *logger.Logger) *TenantService {
      - 2:23:func (s *TenantService) CreateTenant(ctx context.Context, req *CreateTenantRequest) (*TenantResponse, error) {
      - 3:49:func (s *TenantService) GetTenantByID(ctx context.Context, id uuid.UUID) (*TenantResponse, error) {
      - 4:65:func (s *TenantService) GetTenantBySlug(ctx context.Context, slug string) (*TenantResponse, error) {
      - 5:81:func (s *TenantService) UpdateTenant(ctx context.Context, id uuid.UUID, req *UpdateTenantRequest) (*TenantResponse, error) {

### File: pkg/version/version.go
- Package: package version
- Imports:
- Declarations:
  - Line 4: `const CLIVersion = "1.0.0"`
    - Kind: const
    - Symbol: `CLIVersion`
    - Purpose: Constant definition for stable configuration/value semantics.
    - Called/Referenced From: No external references found by textual search.

## 6. Non-Go Assets and Operational Features

### File: README.md
- Role: Repository-level documentation/configuration/dependency metadata.
- Notes: Reviewed as part of root-scope audit.

### File: docs/ARCHITECTURE.md
- Role: Repository-level documentation/configuration/dependency metadata.
- Notes: Reviewed as part of root-scope audit.

### File: OSS_VS_ENTERPRISE.md
- Role: Repository-level documentation/configuration/dependency metadata.
- Notes: Reviewed as part of root-scope audit.

### File: VERSIONING.md
- Role: Repository-level documentation/configuration/dependency metadata.
- Notes: Reviewed as part of root-scope audit.

### File: go.mod
- Role: Repository-level documentation/configuration/dependency metadata.
- Notes: Reviewed as part of root-scope audit.

### File: go.sum
- Role: Repository-level documentation/configuration/dependency metadata.
- Notes: Reviewed as part of root-scope audit.

### File: .env
- Role: Repository-level documentation/configuration/dependency metadata.
- Notes: Reviewed as part of root-scope audit.

### SQL Migrations
- migrations/000001_init_schema.down.sql: schema evolution step for platform domains and infra.
- migrations/000001_init_schema.up.sql: schema evolution step for platform domains and infra.
- migrations/000002_2.down.sql: schema evolution step for platform domains and infra.
- migrations/000002_2.up.sql: schema evolution step for platform domains and infra.
- migrations/000003_000003_notification_deliveries.down.sql: schema evolution step for platform domains and infra.
- migrations/000003_000003_notification_deliveries.up.sql: schema evolution step for platform domains and infra.
- migrations/000004_audit_logs.down.sql: schema evolution step for platform domains and infra.
- migrations/000004_audit_logs.up.sql: schema evolution step for platform domains and infra.

## 7. Architectural Observations

- The codebase follows layered composition: handler -> service -> repository/provider.
- Platform package centralizes cross-cutting concerns (security, context, middleware, response standards, eventing).
- Strong modular separation by business domain under `internal/`.
- CLI and server runtimes share internal modules and configuration conventions.
- Migrations and repository code indicate PostgreSQL-backed persistence as system of record.
- Redis appears as cache/queue support dependency.
- Observability is both first-class (audit logs, metrics, tracing) and integrated into handlers/services.

## 8. Coverage Statement

This report covers the whole repository root with detailed symbol extraction for all Go source files and functional summaries for non-Go artifacts.

## 9. End of Report
