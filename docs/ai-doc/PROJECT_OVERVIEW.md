# Project Overview

Go library to manage **libvirt storage pools** from application code. Import `storagePool`, connect via hypervisor URI, call typed methods. **No HTTP server required** for production use.

Module: `github.com/Hari-Kiri/virest`

## Stack

| Layer | Choice |
|-------|--------|
| Language | Go |
| Hypervisor API | go-libvirt / libvirtxml |
| Helpers | in-repo `virestUtilities` (connect + error wrap) |
| Optional HTTP | `cmd/testserver` (JWT, swaggo) |

## Primary entry

| Surface | Symbol | Role |
|---------|--------|------|
| Library | `storagePool.Connect` | Primary public API |

## Layout

- `storagePool/` — Connection methods, types, adapters
- `virestUtilities/` — Wrap / ConnectWithAuth / Code
- `cmd/testserver/handlers` — HTTP → `storagePool`
- `examples/basic` — sample consumer

## Related Docs

- [ARCHITECTURE.md](ARCHITECTURE.md)
- [modules/storage-pool.md](modules/storage-pool.md)
- [INDEX.md](INDEX.md)
