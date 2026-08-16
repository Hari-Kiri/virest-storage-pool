# Middleware Overview

## Purpose

Document auth gating pattern. **No** dedicated middleware package or chain.

## Pattern

| Helper | Role |
|--------|------|
| `Deps.authorize` | Parse Bearer → require `Hypervisor-Uri` → allowlist |
| `Deps.connect` | authorize + `storagepool.Connect(uri)` |
| Method checks | Per-handler |
| `decodeJSON` | Body validate |

## Permissions / Roles

JWT carries `role` but handlers do **not** currently gate by role. Authorization = valid token + URI allowlist.

## Related Docs

- [auth/authentication.md](../auth/authentication.md)
- [api/endpoints.md](../api/endpoints.md)
- [INDEX.md](../INDEX.md)
