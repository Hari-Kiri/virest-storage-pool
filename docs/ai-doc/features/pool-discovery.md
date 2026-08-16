# Feature: Pool Discovery

## Purpose

Discover supported pool types and remote/local sources before Define.

## API

| Op | Behavior |
|----|----------|
| `Capabilities()` | `GetStoragePoolCapabilities` → unmarshal `Capabilities` |
| `FindSources(poolType, src)` | Marshal `SourceSpec` XML → `FindStoragePoolSources` → `Sources` |

## HTTP

- GET `/storage-pool/capabilities`
- POST `/storage-pool/find-storage-pool-sources` body `{type, srcSpec}`

## Related Docs

- [backend/types.md](../backend/types.md)
- [integrations/libvirt.md](../integrations/libvirt.md)
- [INDEX.md](../INDEX.md)
