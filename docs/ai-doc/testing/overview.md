# Testing

## Make targets (preferred)

```bash
make test        # unit tests + -cover on library and cmd/testserver
make test-cover  # library cover profiles; fail if storagePool or virestUtilities < 70%
make bench       # -bench=. -benchmem (ns/op, B/op, allocs/op) for HA review
make check       # test-cover + bench
```

`CGO_ENABLED=1` by default. Override: `COVER_MIN=70 make test-cover`.

## Unit tests

```bash
go test ./storagePool/ ./virestUtilities/ ./cmd/testserver/... -cover
```

Coverage bar (library): **≥70%** statement coverage on `storagePool/` / `virestUtilities/`, with **positive and negative** cases. `cmd/testserver` reports cover; no hard floor in Make.

```bash
go test ./storagePool/ ./virestUtilities/ -cover
go test ./storagePool/ ./virestUtilities/ -coverprofile=cover.out
go tool cover -func=cover.out | grep total
```

Thin `libvirt_adapter.go` passthrough may be excluded from the bar (covered by integration); all other library logic counts toward ≥70%.

## Benchmarks (HA-safe)

Always add or update `Benchmark*` tests for touched hot paths (HA orchestrator machines).

```bash
make bench
# or:
go test ./storagePool/ ./virestUtilities/ -run='^$' -bench=. -benchmem
```

Review **ns/op**, **B/op**, and **allocs/op**; fix regressions (goroutine storms, excess allocs, O(n²)). List path is sequential per pool (HA-safe). Info keeps fixed 4-way parallel for a single pool.

## Integration

```bash
CGO_ENABLED=1 go test -tags=integration ./storagePool/
```

## Related Docs

- [CODING_STANDARD.md](../CODING_STANDARD.md)
- [deployment/ci.md](../deployment/ci.md)
- [INDEX.md](../INDEX.md)
