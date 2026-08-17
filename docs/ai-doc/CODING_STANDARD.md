# Coding Standard

Rules: `.cursor/rules/golang-standards.mdc`, `.cursor/rules/mission-critical-library.mdc`  
Skill: `.cursor/skills/go-libvirt-library/SKILL.md`

## Role

Senior Golang + go-libvirt expert. Production library for **HA hypervisor cluster** nodes; same style bar for `cmd/testserver/` and `examples/`. Mirror `storagePool/` structure/style; keep handlers/examples thin.

## Scope

| Path | Applies |
|------|---------|
| `storagePool/`, `utilities/` | Full bar (style, deps, lint, ≥70% coverage, hot-path benchmarks) |
| `cmd/testserver/` | Style, lint, testserver deps; ≥70% on touched `auth/` / `handlers/`; benchmarks if hot path |
| `examples/` | Style, lint, library+stdlib deps, Close; no coverage/benchmark mandate |

## Spec

| Rule | Expectation |
|------|-------------|
| Complexity | Prefer O(1)/O(n); avoid O(n²)+ on hot paths |
| N+1 | No per-item lookup/reconnect when batch/list suffices |
| Grade | Production: explicit errors, Free/Close, no silent fail |
| Nesting | No nested `if` / nested `for`; extract methods |
| Loops | Prefer `for i := 0; i < n; i++` over `for range` |
| Efficiency | Minimal allocs + libvirt round-trips under HA load |
| Deps (library) | go-libvirt (+ `libvirtxml`), each other, stdlib only |
| Deps (testserver) | + swaggo, golang-jwt, yaml, bcrypt |
| Deps (examples) | library + stdlib only |
| Lint | After every Go task: run [`.golangci.yaml`](../../.golangci.yaml) on touched packages; **fix all findings**. Zero errors and zero warnings |
| Unit coverage | **≥70%** on touched library + testable cmd packages; positive **and** negative cases |
| Benchmarks | Always add/update for touched library hot paths; `-bench=. -benchmem`; HA-safe |

## Structure to copy

- API on `*Connection` in `storagePool`
- Ports: `hypervisor`, `poolHandle`
- Adapter: `libvirt_adapter.go`; tests: fakes + `newConnectionForTest`
- Errors/connect helpers: `utilities`
- Types: `types.go`
- `cmd/testserver`: auth → library connect → envelope; no duplicated pool logic
- `examples`: Connect, defer Close, call library, print

## Tests + benchmarks (mandatory where scoped)

```bash
go test ./storagePool/ ./utilities/ ./cmd/testserver/... -cover
go test ./storagePool/ ./utilities/ -bench=. -benchmem
```

Details: [testing/overview.md](testing/overview.md).

## Lint (mandatory after every Go task)

Project config: [`.golangci.yaml`](../../.golangci.yaml) (`run.tests: false`).

After any change under `storagePool/`, `utilities/`, `cmd/`, or `examples/`:

1. `golangci-lint run ./path/to/touched/package/`
2. If there are **any** issues (errors or warnings), fix the code and re-run until clean.
3. Do not finish the task with dirty lint. Do not `//nolint` what a rewrite can fix.

```bash
golangci-lint run ./path/to/touched/package/
```

## Docs sync

After `storagePool/`, `utilities/`, `cmd/`, `examples/` changes → one primary page via [INDEX.md](INDEX.md) / [CATALOG.md](CATALOG.md).

## Related Docs

- [ARCHITECTURE.md](ARCHITECTURE.md)
- [modules/storage-pool.md](modules/storage-pool.md)
- [testing/overview.md](testing/overview.md)
- [INDEX.md](INDEX.md)
