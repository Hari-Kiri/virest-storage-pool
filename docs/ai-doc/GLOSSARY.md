# Glossary

| Term | Meaning |
|------|---------|
| Storage pool | Libvirt object managing a storage backend (dir, netfs, logical, …) |
| Volume | Storage object inside a pool (library focuses on **pools**, not volume CRUD) |
| Define | Create inactive pool definition from XML; returns UUID |
| Build | Prepare underlying storage for a defined pool |
| Create | Start / activate a pool |
| Destroy | Stop an active pool (definition may remain) |
| Delete | Remove pool backing resources |
| Undefine | Remove pool definition |
| Refresh | Rescan pool volumes |
| Autostart | Start pool when host/libvirt starts |
| UUID | Stable pool id used by almost all methods |
| Hypervisor URI | Libvirt connect string, e.g. `qemu:///system` |
| Connection | `storagepool.Connection` handle |
| hypervisor / poolHandle | Internal interfaces for adapter + tests |
| EventLifecycle / EventRefresh | `EventKind` 0 / 1 |
| Allowlist | Per-user permitted hypervisor URIs (testserver) |
| Envelope | Standard testserver JSON response wrapper |
| Testserver | Optional HTTP lab harness under `cmd/testserver` |
| CGO | Required to link libvirt |

## Related Docs

- [BUSINESS_PROCESS.md](BUSINESS_PROCESS.md)
- [modules/storage-pool.md](modules/storage-pool.md)
- [INDEX.md](INDEX.md)
