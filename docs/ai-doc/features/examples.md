# Feature: Examples

## Purpose

Minimal consumer of the library. Follows the same Go style bar as library/cmd (see [CODING_STANDARD.md](../CODING_STANDARD.md)): indexed loops, `defer Close`, library + stdlib only.

## Path

`examples/basic/main.go`

## Behavior

- URI from `VIREST_LIBVIRT_URI` or default `qemu:///system`
- `Connect` → `defer Close` → `List(0,0)` → print name/uuid/state (indexed loop)
- No coverage/benchmark mandate for example mains

## Related Docs

- [modules/storage-pool.md](../modules/storage-pool.md)
- [configs/environment.md](../configs/environment.md)
- [CODING_STANDARD.md](../CODING_STANDARD.md)
- [INDEX.md](../INDEX.md)
