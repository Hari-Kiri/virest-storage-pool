# Types

Source: `storagepool/types.go`

| Type | Purpose | Key fields |
|------|---------|------------|
| `Info` | Volatile pool stats | name, uuid, state, persistent, autostart, capacity/allocation/available |
| `Detail` | Full XML model | embeds `libvirtxml.StoragePool` |
| `PoolSummary` | List row | uuid, name, state, autostart, persistent, size structs |
| `EventKind` | 0 lifecycle, 1 refresh | constants `EventLifecycle`, `EventRefresh` |
| `Event` | Notification | EventLifecycle, EventRefresh, timestamps, EventId |
| `SourceSpec` | FindSources input XML | Host, Initiator/IQN |
| `Sources` | FindSources result | `[]libvirtxml.StoragePoolSource` |
| `Capabilities` / `CapabilityPool` | Capability XML | pool types, vol/pool format enums |

## Constraints

JSON tags on exported fields for testserver serialization. XML tags where marshaled to libvirt.

## Related Docs

- [modules/storage-pool.md](../modules/storage-pool.md)
- [features/pool-query.md](../features/pool-query.md)
- [INDEX.md](../INDEX.md)
