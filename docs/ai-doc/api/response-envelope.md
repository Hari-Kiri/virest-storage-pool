# Response Envelope & Errors

## Envelope

```json
{"response": true|false, "code": <httpStatus>, "data": ..., "error": "..."}
```

Success: `response=true`, optional `data`. Error: `response=false`, `error` message.

## HTTP status mapping (`httpStatusFor`)

| Condition | Status |
|-----------|--------|
| libvirt `ERR_AUTH_FAILED` | 401 |
| `ERR_NO_STORAGE_POOL`, `ERR_NO_SUPPORT` | 404 |
| `ERR_INVALID_ARG`, `ERR_INVALID_CONN` | 400 |
| Other libvirt codes (via `viresterrors.Code`) | 500 |
| Message contains bearer/credential/token/not allowed | 401 |
| Message contains hypervisor uri | 400 |
| Else | 500 |

Libvirt codes extracted via ``utilities`/utils/errors.Code`.

## SSE errors

`event: error` + `data: <message>` when stream fails and ctx still active.

## Related Docs

- [api/endpoints.md](endpoints.md)
- [DATA_FLOW.md](../DATA_FLOW.md)
- [INDEX.md](../INDEX.md)
