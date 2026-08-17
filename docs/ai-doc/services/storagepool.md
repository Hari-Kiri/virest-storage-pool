# Service: storagepool

## Responsibilities

Single domain service = the `storagepool` package (no separate service layer).

## Inputs

- Hypervisor URI
- Pool UUID (post-define)
- XML models / flags / `SourceSpec` / `EventKind` / `context`

## Outputs

- Typed structs (`Info`, `Detail`, `PoolSummary`, `Capabilities`, `Sources`, `Event`)
- UUID on Define
- Wrapped errors

## Who calls it

- Application code (import)
- `examples/basic`
- `cmd/testserver/handlers`

## What it calls

- ``utilities`/utils/hypervisor.ConnectWithAuth`
- libvirt via `libvirtHypervisor` / `libvirtPool`
- ``utilities`/utils/errors.Wrap`

## Dependencies

See [DEPENDENCY_GRAPH.md](../DEPENDENCY_GRAPH.md)

## Related Docs

- [modules/storage-pool.md](../modules/storage-pool.md)
- [backend/library-connection.md](../backend/library-connection.md)
- [INDEX.md](../INDEX.md)
