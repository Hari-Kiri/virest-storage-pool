# Authentication

Applies to **testserver only**. Library has no app-level auth.

## Login Flow

1. `GET /storage-pool/authenticate` + HTTP Basic
2. Lookup user in YAML store; bcrypt compare
3. Issue JWT HS512: `sub`=username, `role`, `iss`=app name, `iat`/`exp`
4. Client sends `Authorization: Bearer <token>` + `Hypervisor-Uri`

## Session / Cookie / OAuth

**N/A** — stateless JWT only. No cookies, no OAuth.

## Users file

`users.yaml` (from `users.example.yaml`):

| Field | Rule |
|-------|------|
| username | required |
| password_hash | preferred bcrypt |
| password | lab-only; hashed at load, cleared in memory |
| role | string claim (e.g. operator/reader) — **not enforced** on routes today |
| hypervisor_allowlist | if empty → all URIs allowed; else exact match |

## Middleware

No separate middleware package. Handlers call `Deps.authorize` / `connect` ([middleware/overview.md](../middleware/overview.md)).

## Security Constraints

- Signing key required (`VIREST_STORAGE_POOL_APPLICATION_JWT_SIGNATURE_KEY`)
- Issuer must match app name
- Unexpected JWT alg rejected (HS512 only)
- Prefer TLS via `VIREST_TLS_*` beyond localhost
- Do not commit real `users.yaml` / `.env` (gitignored)

## Related Docs

- [configs/environment.md](../configs/environment.md)
- [api/endpoints.md](../api/endpoints.md)
- [INDEX.md](../INDEX.md)
