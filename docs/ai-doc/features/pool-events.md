# Feature: Pool Events

## Purpose

Observe lifecycle or refresh events for a pool.

## Prerequisites

Process must call:

1. `libvirt.EventRegisterDefaultImpl()`
2. Loop `libvirt.EventRunDefaultImpl()` (testserver does this in `main`)

## API

| Op | Behavior |
|----|----------|
| `WaitEvent(uuid, kind)` | Register → block on first event → deregister |
| `StreamEvents(ctx, uuid, kind, emit)` | Register → loop emit until ctx done or emit error |

`EventKind`: 0 lifecycle, 1 refresh. Invalid kind → error. `emit == nil` → error.

## Constraints / Warnings

- Needs native libvirt pool (`underlying()`); fakes return `errNoNativePool`
- Stream buffer size 8; full buffer **drops** events (`select default`)
- ctx cancel / deadline → StreamEvents returns nil

## HTTP

GET `/storage-pool/event?uuid=&types=&stream=1&timeout=`

- Default: WaitEvent JSON envelope
- `stream=1`: SSE `text/event-stream`

## Related Docs

- [DATA_FLOW.md](../DATA_FLOW.md)
- [integrations/libvirt.md](../integrations/libvirt.md)
- [features/testserver.md](testserver.md)
- [INDEX.md](../INDEX.md)
