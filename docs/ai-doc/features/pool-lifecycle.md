# Feature: Pool Lifecycle

## Purpose

Define → build → start → stop → delete resources → undefine.

## Library API

| Op | File | Signature notes |
|----|------|-----------------|
| Define | lifecycle.go | Marshal XML → `StoragePoolDefineXML` → UUID |
| Build | lifecycle.go | `pool.Build(flags)` |
| Create | lifecycle.go | `pool.Create(flags)` |
| Destroy | lifecycle.go | `pool.Destroy()` |
| Delete | lifecycle.go | `pool.Delete(flags)` |
| Undefine | lifecycle.go | `pool.Undefine()` |
| Refresh | lifecycle.go | `Refresh(0)` fixed flags |
| Autostart | lifecycle.go | `SetAutostart(bool)` |

## HTTP (testserver)

See [api/endpoints.md](../api/endpoints.md) — define POST; build/create/destroy/autostart PATCH; undefine/delete DELETE; refresh POST.

## Business Rules

- Typical teardown order: Destroy → Delete → Undefine
- Build may fail harmlessly for some `dir` pools (integration test logs and continues)
- Destructive disk tests gated by `VIREST_TEST_DISK`

## Related Docs

- [BUSINESS_PROCESS.md](../BUSINESS_PROCESS.md)
- [modules/storage-pool.md](../modules/storage-pool.md)
- [INDEX.md](../INDEX.md)
