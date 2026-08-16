# AI Doc Router

**Protocol:** Read this table only → open **≤3** linked docs → verify critical facts in source. Do not scan the repo first. Do not read all of `docs/ai-doc/`.

**Update after code change:** map change → one primary page (see [CATALOG.md](CATALOG.md) update map). Touch INDEX/CATALOG only for new topics.

---

## Topic → Doc

| Topic | Doc |
|-------|-----|
| Project purpose / stack | [PROJECT_OVERVIEW.md](PROJECT_OVERVIEW.md) |
| Architecture / layers | [ARCHITECTURE.md](ARCHITECTURE.md) |
| Directory map | [DIRECTORY_STRUCTURE.md](DIRECTORY_STRUCTURE.md) |
| Business flows (pool lifecycle) | [BUSINESS_PROCESS.md](BUSINESS_PROCESS.md) |
| Request / data / event flow | [DATA_FLOW.md](DATA_FLOW.md) |
| Package dependencies | [DEPENDENCY_GRAPH.md](DEPENDENCY_GRAPH.md) |
| Coding rules | [CODING_STANDARD.md](CODING_STANDARD.md) |
| Terms | [GLOSSARY.md](GLOSSARY.md) |
| Fallback catalog | [CATALOG.md](CATALOG.md) |
| Library API (`storagePool`) | [modules/storage-pool.md](modules/storage-pool.md) |
| Connection / Connect / Close | [backend/library-connection.md](backend/library-connection.md) |
| Define / Build / Create / Destroy / Delete / Undefine / Refresh / Autostart | [features/pool-lifecycle.md](features/pool-lifecycle.md) |
| List / Info / Detail | [features/pool-query.md](features/pool-query.md) |
| Capabilities / FindSources | [features/pool-discovery.md](features/pool-discovery.md) |
| WaitEvent / StreamEvents | [features/pool-events.md](features/pool-events.md) |
| Types / structs | [backend/types.md](backend/types.md) |
| Hypervisor adapters / fakes | [backend/adapters.md](backend/adapters.md) |
| Library service surface | [services/storagepool.md](services/storagepool.md) |
| Backend overview | [backend/overview.md](backend/overview.md) |
| HTTP testserver | [features/testserver.md](features/testserver.md) |
| HTTP endpoints | [api/endpoints.md](api/endpoints.md) |
| Response envelope / errors | [api/response-envelope.md](api/response-envelope.md) |
| Auth / JWT / allowlist | [auth/authentication.md](auth/authentication.md) |
| Auth middleware pattern | [middleware/overview.md](middleware/overview.md) |
| Env / config / secrets | [configs/environment.md](configs/environment.md) |
| Libvirt integration | [integrations/libvirt.md](integrations/libvirt.md) |
| Testserver run / TLS / Swagger | [deployment/testserver.md](deployment/testserver.md) |
| CI | [deployment/ci.md](deployment/ci.md) |
| Tests | [testing/overview.md](testing/overview.md) |
| Examples | [features/examples.md](features/examples.md) |
| Frontend | [frontend/overview.md](frontend/overview.md) |
| Database | [database/overview.md](database/overview.md) |
| AI tooling / doc-sync config | [AI_TOOLING.md](AI_TOOLING.md) |

---

## Task shortcuts

| Task shape | Read (≤3) |
|------------|-----------|
| Add/change library method | modules/storage-pool → features/\<area\> → backend/types |
| Change pool lifecycle | features/pool-lifecycle → modules/storage-pool → BUSINESS_PROCESS |
| Change events | features/pool-events → integrations/libvirt → DATA_FLOW |
| Change HTTP route/handler | api/endpoints → features/testserver → auth/authentication |
| Change auth | auth/authentication → configs/environment → api/endpoints |
| Change env/config | configs/environment → deployment/testserver |
| Add tests | testing/overview → features/\<area\> |
| Understand architecture | PROJECT_OVERVIEW → ARCHITECTURE → DEPENDENCY_GRAPH |
