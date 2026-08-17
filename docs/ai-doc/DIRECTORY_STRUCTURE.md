# Directory Structure

| Path | Role | Kind | Notes |
|------|------|------|-------|
| `storagePool/` | Core library | library | Connection API, types, adapters, tests | Primary product | No HTTP/auth |
| `utilities/` | Shared helpers | library | Wrap, Connect*, Code | Used by storagePool + testserver |
| `cmd/testserver/` | HTTP harness | tooling | JWT auth, Swagger | Optional |
| `cmd/testserver/handlers/` | HTTP handlers | tooling | routes via `main` | Logic stays in `storagePool` |
| `cmd/testserver/auth/` | Users/JWT | tooling | yaml users file | |
| `examples/basic/` | Sample app | example | | |
| `docs/ai-doc/` | Agent docs | docs | INDEX router | |

## Related Docs

- [modules/storage-pool.md](modules/storage-pool.md)
- [INDEX.md](INDEX.md)
