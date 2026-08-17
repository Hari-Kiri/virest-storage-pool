# Feature: Testserver

## Purpose

Optional HTTP harness for manual/e2e exercise of `storagepool`. **Not** the primary product.

## Responsibilities

- Load `.env` + users YAML
- JWT auth store
- Register libvirt event loop
- Route HTTP → handlers → library
- Serve Swagger UI

## Entry

`cmd/testserver/main.go` — `main` loads env, builds auth via `mustAuthStore`, starts libvirt events (`startLibvirtEvents`), registers routes (`registerRoutes`), then listens (`runServer`).

Handler shared code: `cmd/testserver/handlers/common.go` (`authorize` → hypervisor URI; `uuidData` for UUID JSON payloads; `httpStatusFor` via `utilities.Code`).

## Used By

Developers / CI-like local testing. Consumers of the **library** should not require this binary.

## Dependencies

`storagepool`, `auth`, swaggo, libvirt event APIs

## Configuration

[configs/environment.md](../configs/environment.md)

## Related Docs

- [api/endpoints.md](../api/endpoints.md)
- [auth/authentication.md](../auth/authentication.md)
- [deployment/testserver.md](../deployment/testserver.md)
- [INDEX.md](../INDEX.md)
