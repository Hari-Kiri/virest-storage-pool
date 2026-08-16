# CI / CD

## GitHub Actions

`.github/workflows/enforcer.yaml`

- Trigger: `pull_request`
- Rule: PRs into `main` must come from `development` only
- **No** build/test/lint job in this workflow

## Local quality gates

Prefer Make (see [testing/overview.md](../testing/overview.md)):

```bash
make check       # library ≥70% cover + HA benches (-benchmem)
make test        # unit tests with -cover
make test-cover  # cover profiles; fail below COVER_MIN (default 70)
make bench       # ns/op, B/op, allocs/op for HA review
```

`CGO_ENABLED=1` by default. Integration tag optional (`go test -tags=integration`).

## Docker / cloud deploy

**N/A** in repository.

## Related Docs

- [testing/overview.md](../testing/overview.md)
- [deployment/testserver.md](testserver.md)
- [INDEX.md](../INDEX.md)
