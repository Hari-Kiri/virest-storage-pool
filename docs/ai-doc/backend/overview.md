# Backend Overview

## Purpose

All server-side logic lives in Go packages; there is no separate app server for production.

## Map

| Area | Path | Doc |
|------|------|-----|
| Library core | `storagepool/` | [modules/storage-pool.md](../modules/storage-pool.md) |
| Connection | `connect.go` | [library-connection.md](library-connection.md) |
| Types | `types.go` | [types.md](types.md) |
| Adapters | `hypervisor.go`, `libvirt_adapter.go`, `fake_test.go` | [adapters.md](adapters.md) |
| Lifecycle | `lifecycle.go` | [../features/pool-lifecycle.md](../features/pool-lifecycle.md) |
| Query | `list_info.go`, `detail.go` | [../features/pool-query.md](../features/pool-query.md) |
| Discovery | `detail.go` (Capabilities/FindSources) | [../features/pool-discovery.md](../features/pool-discovery.md) |
| Events | `event.go` | [../features/pool-events.md](../features/pool-events.md) |
| Errors | `errors.go` | wrap → `utilities` |
| HTTP | `cmd/testserver` | [../features/testserver.md](../features/testserver.md) |

## Controllers / Repositories / Cron / Queue

**N/A** as separate layers. Testserver handlers ≈ thin controllers. No repos/cron/queues.

## Related Docs

- [ARCHITECTURE.md](../ARCHITECTURE.md)
- [INDEX.md](../INDEX.md)
