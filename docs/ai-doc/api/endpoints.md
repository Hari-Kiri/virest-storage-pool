# API Endpoints

Base: testserver (default `:8000`). Prefix `/storage-pool`.

## Auth

| Route | Method | Security | Notes |
|-------|--------|----------|-------|
| `/authenticate` | GET | Basic | Returns `{data.token}` |

All other routes: **Bearer** + header **`Hypervisor-Uri`** (allowlisted). Meta routes still require both for authorize check.

## Pools

| Route | Method | Query / Body | Library |
|-------|--------|--------------|---------|
| `/list` | GET | `option`, `inactive` | `List` |
| `/info` | GET | `uuid` | `Info` |
| `/detail` | GET | `uuid`, `option` | `Detail` |
| `/capabilities` | GET | — | `Capabilities` |
| `/define` | POST | `{option, storagePool}` | `Define` → 201 `{uuid}` |
| `/build` | PATCH | `{uuid, option}` | `Build` |
| `/create` | PATCH | `{uuid, option}` | `Create` |
| `/destroy` | PATCH | `{uuid}` | `Destroy` |
| `/autostart` | PATCH | `{uuid, autostart}` | `Autostart` |
| `/undefine` | DELETE | `{uuid}` | `Undefine` |
| `/delete` | DELETE | `{uuid, option}` | `Delete` |
| `/refresh` | POST | query `uuid` or body `{uuid}` | `Refresh` |
| `/find-storage-pool-sources` | POST | `{type, srcSpec}` | `FindSources` |

## Events

| Route | Method | Query | Behavior |
|-------|--------|-------|----------|
| `/event` | GET | `uuid`, `types`, `stream`, `timeout` | Wait or SSE |

## Meta (process, not hypervisor)

| Route | Method | Data |
|-------|--------|------|
| `/get-uid` | GET | `{uid}` |
| `/get-gid` | GET | `{gid}` |

## Swagger

`/swagger/` → OpenAPI from annotations; regenerate: `make swagger`.

## Validation

- JSON decode `DisallowUnknownFields`
- Method checks → 405
- Auth/allowlist failures → see [response-envelope.md](response-envelope.md)

## Related Services

[services/storagepool.md](../services/storagepool.md), [auth/authentication.md](../auth/authentication.md)

## Related Docs

- [api/response-envelope.md](response-envelope.md)
- [features/testserver.md](../features/testserver.md)
- [INDEX.md](../INDEX.md)
