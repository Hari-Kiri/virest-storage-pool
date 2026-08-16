# Adapters

## Purpose

Isolate libvirt CGO behind interfaces for unit tests.

## Interfaces (`hypervisor.go`)

- `hypervisor` — connection ops (list, lookup, define XML, capabilities, find sources, event register/deregister)
- `poolHandle` — pool ops + `underlying() *libvirt.StoragePool` for events

## Production (`libvirt_adapter.go`)

- `libvirtHypervisor` / `libvirtPool`
- Event register requires native pool; else `errNoNativePool`

## Tests (`fake_test.go`)

- Fakes implement interfaces; used by `storagepool_test.go` without live libvirt

## Related Docs

- [library-connection.md](library-connection.md)
- [testing/overview.md](../testing/overview.md)
- [INDEX.md](../INDEX.md)
