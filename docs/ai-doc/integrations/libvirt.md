# Integration: Libvirt

## Purpose

Hypervisor backend for all pool operations.

## Requirements

- CGO on
- `libvirt` / `libvirt-dev`, gcc
- Reachable daemon (e.g. `qemu:///system` or session)
- User in `libvirt`/`kvm` groups as needed

## Go modules

- `libvirt.org/go/libvirt`
- `libvirt.org/go/libvirtxml`
- Connect helper: `virest-utilities/utils/hypervisor`

## URI

See https://libvirt.org/uri.html — pass as `Connect(uri)` or testserver `Hypervisor-Uri`.

## Events

Must run default event impl in-process or WaitEvent/StreamEvents hang/fail.

## Related Docs

- [backend/adapters.md](../backend/adapters.md)
- [features/pool-events.md](../features/pool-events.md)
- [INDEX.md](../INDEX.md)
