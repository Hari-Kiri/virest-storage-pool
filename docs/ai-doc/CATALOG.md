# Catalog (fallback)

Use when [INDEX.md](INDEX.md) has no row. Prefer INDEX.

## All pages

| Doc | When |
|-----|------|
| PROJECT_OVERVIEW.md | Stack, purpose, constraints |
| ARCHITECTURE.md | Layers, boundaries |
| DIRECTORY_STRUCTURE.md | Folders |
| BUSINESS_PROCESS.md | Pool lifecycle business flow |
| DATA_FLOW.md | Call / HTTP / auth / event flows |
| DEPENDENCY_GRAPH.md | Package deps |
| CODING_STANDARD.md | Go + docs sync rules |
| GLOSSARY.md | Terms |
| modules/storage-pool.md | Full library surface |
| services/storagepool.md | Service view of library |
| backend/overview.md | Backend map |
| backend/library-connection.md | Connect/Close/lookup |
| backend/types.md | Exported types |
| backend/adapters.md | libvirt + fake adapters |
| features/pool-lifecycle.md | Mutating lifecycle ops |
| features/pool-query.md | List/Info/Detail |
| features/pool-discovery.md | Capabilities/FindSources |
| features/pool-events.md | Events |
| features/testserver.md | HTTP harness feature |
| features/examples.md | examples/basic |
| api/endpoints.md | Routes |
| api/response-envelope.md | JSON/SSE errors |
| auth/authentication.md | JWT/users |
| middleware/overview.md | Auth-in-handler pattern |
| configs/environment.md | Env vars |
| integrations/libvirt.md | Libvirt/CGO |
| deployment/testserver.md | Run/Swagger/TLS |
| deployment/ci.md | GitHub Actions |
| testing/overview.md | Unit/integration |
| frontend/overview.md | N/A |
| database/overview.md | N/A |
| AI_TOOLING.md | Cursor rules/skills/hooks; multi-agent sync protocol |

## Update map (code → one primary doc)

| Changed path | Primary update |
|--------------|----------------|
| `storagePool/lifecycle.go` | features/pool-lifecycle.md |
| `storagePool/list_info.go`, `detail.go` (query) | features/pool-query.md or pool-discovery.md |
| `storagePool/event.go` | features/pool-events.md |
| `storagePool/types.go` | backend/types.md |
| `storagePool/connect.go`, `hypervisor.go`, `libvirt_adapter.go` | backend/library-connection.md or adapters.md |
| `utilities/*` | backend/library-connection.md or CODING_STANDARD |
| `storagePool` public method set | modules/storage-pool.md (if surface changes) |
| `cmd/testserver/handlers/*`, `main.go` routes | api/endpoints.md (+ features/testserver.md if behavior) |
| `cmd/testserver/auth/*` | auth/authentication.md |
| `cmd/testserver/.env_dummy`, users example | configs/environment.md |
| `examples/*` | features/examples.md |
| New topic | new page + INDEX row + this catalog row |

Skip overview/architecture/DATA_FLOW unless structural.

## Related Docs

- [INDEX.md](INDEX.md)
