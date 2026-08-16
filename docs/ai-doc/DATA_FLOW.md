# Data Flow

## Library call flow

```text
Caller → Connection method → lookupPool / hv.* → poolHandle / libvirt
       → XML marshal/unmarshal as needed → typed result or wrap(err)
```

## Request flow (testserver)

```text
HTTP → mux route → handler
  → authorize: Bearer JWT + Hypervisor-Uri allowlist
  → storagepool.Connect(uri)
  → library method
  → Envelope JSON (or SSE for stream events)
  → conn.Close()
```

## Response flow

| Path | Shape |
|------|-------|
| Success | `{response:true, code:<http>, data:...}` |
| Error | `{response:false, code:<http>, error:"..."}` |
| SSE | `data: <Event JSON>\n\n`; errors as `event: error` |

Status mapping: [api/response-envelope.md](api/response-envelope.md).

## Database / cache / queue / notification

**N/A** — no persistence layer, cache, or message queue in this repo.

## Authentication flow

```text
GET /storage-pool/authenticate + Basic
  → bcrypt verify → JWT HS512 (sub=user, role, iss=app, exp)
Subsequent:
  Authorization: Bearer <jwt>
  Hypervisor-Uri: <allowlisted URI>
```

Details: [auth/authentication.md](auth/authentication.md).

## Event flow

```text
Process: EventRegisterDefaultImpl + loop EventRunDefaultImpl
Register callback on pool → libvirt fires → channel/callback → Event struct
WaitEvent: block until one
StreamEvents: emit until ctx cancel / emit error (overflow drops with default select)
```

## Related Docs

- [ARCHITECTURE.md](ARCHITECTURE.md)
- [features/pool-events.md](features/pool-events.md)
- [api/endpoints.md](api/endpoints.md)
- [INDEX.md](INDEX.md)
