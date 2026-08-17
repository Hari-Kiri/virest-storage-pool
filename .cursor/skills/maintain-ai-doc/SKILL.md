---
name: maintain-ai-doc
description: >-
  Token-budgeted docs/ai-doc: read INDEX routing (≤3 pages), update one primary
  doc after storagePool/utilities/cmd/examples changes. Use for codebase
  Q&A, plans, library/handlers/auth, coding-standard sync, or doc sync.
---

# Maintain AI docs

## Read (≤3 pages)

1. `docs/ai-doc/INDEX.md` — routing table only.
2. Open at most **3** linked docs from that row.
3. Open `docs/ai-doc/CATALOG.md` only if no routing row matches.
4. Verify critical behavior in source (`storagePool/`, `utilities/`, `cmd/testserver/`).

Do not re-read INDEX/CATALOG if already loaded this turn.

## Update (minimal)

Map change → **one** primary page (prefer CATALOG update map):

| Changed | Primary update |
|---------|----------------|
| `storagePool/lifecycle.go` | `features/pool-lifecycle.md` |
| `storagePool/list_info.go` / query Detail | `features/pool-query.md` |
| Capabilities / FindSources | `features/pool-discovery.md` |
| `storagePool/event.go` | `features/pool-events.md` |
| `storagePool/types.go` | `backend/types.md` |
| connect / hypervisor / adapter | `backend/library-connection.md` or `backend/adapters.md` |
| `utilities/*` | `backend/library-connection.md` or CODING_STANDARD |
| Public method surface | `modules/storage-pool.md` |
| `cmd/testserver/handlers/*` or routes in `main.go` | `api/endpoints.md` |
| `cmd/testserver/auth/*` | `auth/authentication.md` |
| env / users example | `configs/environment.md` |
| `examples/*` | `features/examples.md` |
| new topic | new page + INDEX row and/or CATALOG entry |

Skip overview/architecture/DATA_FLOW unless the change is structural.

## Style

Bullets/tables only. Behavior and contracts, not tutorials. No app source edits for docs-only tasks.

## Agent protocol (every coding prompt)

- Main docs entry: `docs/ai-doc/INDEX.md`.
- If implementation changes anything under `storagePool/`, `utilities/`, `cmd/`, or `examples/`, update the matching primary doc before finishing.
- Do not read the entire `docs/ai-doc/` tree — only INDEX + ≤3 linked pages.
