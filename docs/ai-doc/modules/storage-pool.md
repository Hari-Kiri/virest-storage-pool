# Module: storage-pool

## Purpose

Primary product module: typed Go API over libvirt storage pools.

## Responsibilities

- Connect/close hypervisor
- Lifecycle: Define, Build, Create, Destroy, Delete, Undefine, Refresh, Autostart
- Query: List, Info, Detail
- Discovery: Capabilities, FindSources
- Events: WaitEvent, StreamEvents

## Public Interface

Package: `github.com/Hari-Kiri/virest/storagePool`

| Method | Inputs | Outputs |
|--------|--------|---------|
| `Connect(uri)` | URI string | `*Connection`, error |
| `Close()` | — | error |
| `Define(pool, flags)` | `libvirtxml.StoragePool`, flags | UUID string, error |
| `Build(uuid, flags)` | UUID, build flags | error |
| `Create(uuid, flags)` | UUID, create flags | error |
| `Destroy(uuid)` | UUID | error |
| `Delete(uuid, flags)` | UUID, delete flags | error |
| `Undefine(uuid)` | UUID | error |
| `Refresh(uuid)` | UUID | error |
| `Autostart(uuid, bool)` | UUID, flag | error |
| `List(listFlags, xmlFlags)` | uint, uint | `[]PoolSummary`, error |
| `Info(uuid)` | UUID | `Info`, error |
| `Detail(uuid, xmlFlags)` | UUID, uint | `Detail`, error |
| `Capabilities()` | — | `Capabilities`, error |
| `FindSources(type, spec)` | string, `SourceSpec` | `Sources`, error |
| `WaitEvent(uuid, kind)` | UUID, `EventKind` | `Event`, error |
| `StreamEvents(ctx, uuid, kind, emit)` | ctx, UUID, kind, callback | error |

Types: [backend/types.md](../backend/types.md)

## Dependencies

| Dep | Role |
|-----|------|
| `libvirt` / `libvirtxml` | Native pool ops / XML models |
| `utilities` | `Wrap`, `ConnectWithAuth`, `Code` |
| stdlib `sync` | `runParallel` for List/Info attribute reads |

## Database / Frontend / Queue

N/A

## API (optional harness)

[api/endpoints.md](../api/endpoints.md)

## Events

[features/pool-events.md](../features/pool-events.md)

## Configuration

Library: URI only. Runtime env for tests/examples: [configs/environment.md](../configs/environment.md)

## Related Modules / Docs

- [services/storagepool.md](../services/storagepool.md)
- [features/pool-lifecycle.md](../features/pool-lifecycle.md)
- [features/pool-query.md](../features/pool-query.md)
- [features/pool-discovery.md](../features/pool-discovery.md)
- [INDEX.md](../INDEX.md)
