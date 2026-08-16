# Environment & Config

## Library / examples / tests

| Variable | Used by | Purpose | Secret? |
|----------|---------|---------|---------|
| `VIREST_LIBVIRT_URI` | examples, integration tests | Hypervisor URI (default `qemu:///system`) | No |
| `VIREST_TEST_POOL_DIR` | integration | Pool target dir | No |
| `VIREST_TEST_DISK` | integration | Gate destructive disk tests | No |
| `CGO_ENABLED` | build/test | Must be `1` | No |

## Testserver (`.env` from `.env_dummy`)

| Variable | Purpose | Secret? |
|----------|---------|---------|
| `VIREST_STORAGE_POOL_APPLICATION_NAME` | JWT issuer / app name | No |
| `VIREST_STORAGE_POOL_APPLICATION_PORT` | Listen port (default 8000) | No |
| `VIREST_STORAGE_POOL_APPLICATION_JWT_LIFETIME_SECONDS` | Token TTL | No |
| `VIREST_STORAGE_POOL_APPLICATION_JWT_SIGNATURE_KEY` | HS512 key | **Yes** |
| `VIREST_USERS_FILE` | Path to users YAML | No (file may hold hashes) |
| `VIREST_TLS_CERT_FILE` | TLS cert path | No |
| `VIREST_TLS_KEY_FILE` | TLS key path | **Yes** (file) |

## Files

| File | Role |
|------|------|
| `cmd/testserver/.env_dummy` | Template (safe defaults) |
| `cmd/testserver/.env` | Local runtime (gitignored) |
| `cmd/testserver/users.example.yaml` | User template |
| `cmd/testserver/users.yaml` | Runtime users (gitignored) |
| `go.mod` / `Makefile` | Module + swagger regen |

## Feature flags

**None.**

## Dotenv loader

`main.loadDotEnv` — simple KEY=VALUE; no export/quotes handling. Missing `.env` = silent skip.

## Related Docs

- [deployment/testserver.md](../deployment/testserver.md)
- [auth/authentication.md](../auth/authentication.md)
- [INDEX.md](../INDEX.md)
