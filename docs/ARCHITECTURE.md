# AARCSX Forge Architecture Skeleton

This repository starts as a **modular monolith** with extraction-ready seams.

## Runtime split

- `cmd/api`: HTTP/API runtime
- `cmd/worker`: background jobs runtime

## Foundation packages

- `internal/config`: typed configuration and startup validation
- `internal/platform/events`: internal domain event bus + versioned contracts
- `internal/platform/contextx`: typed request context propagation helpers
- `internal/platform/httpx`: standardized response envelope contract
- `internal/platform/security`: auth/security interfaces

## Domain contracts

- `internal/identity`
- `internal/tenants`
- `internal/storage`
- `internal/notifications`
- `internal/observability`

## Governance highlights

- No `os.Getenv` outside `internal/config`
- Repositories contain persistence only; no business or authorization logic
- Shared package cannot contain domain workflows
- API versioning begins under `/api/v1`
