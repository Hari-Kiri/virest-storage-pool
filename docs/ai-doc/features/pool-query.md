# Feature: Pool Query

## Purpose

Inspect pools without mutating.

## API

| Op | Behavior |
|----|----------|
| `Info(uuid)` | Parallel GetName/GetInfo/GetAutostart/IsPersistent via `runParallel` + Ref/Free (`loadInfoAttrs`) |
| `List(listFlags, xmlFlags)` | ListAllStoragePools → indexed loop → `summarizePool` (XML sizes + info) |
| `Detail(uuid, xmlFlags)` | GetXMLDesc → unmarshal `libvirtxml.StoragePool` |

## Warnings

- List: Free each pool; on mid-loop error call `freePoolsFrom` for the remainder
- Free/Ref errors are wrapped and returned when the primary op succeeds
- Info/List hold multiple refs concurrently — Ref failure fails whole call
- Info: fixed 4-way `runParallel` for one pool
- List/`summarizePool`: **sequential** attribute reads (avoids N×4 goroutines on HA orchestrator nodes)

## HTTP

GET `/storage-pool/list|info|detail` — [api/endpoints.md](../api/endpoints.md)

## Related Docs

- [backend/types.md](../backend/types.md)
- [INDEX.md](../INDEX.md)
