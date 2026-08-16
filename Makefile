# Local quality gates and codegen (run from repo root).
#
# HA bar: library packages ≥ COVER_MIN% statement coverage;
# benchmarks print ns/op, B/op, allocs/op for regression review.

CGO_ENABLED ?= 1
export CGO_ENABLED

COVER_MIN ?= 70
TEST_PKGS := ./storagePool/ ./virestUtilities/ ./cmd/testserver/...
COVER_PKGS := ./storagePool/ ./virestUtilities/
BENCH_PKGS := ./storagePool/ ./virestUtilities/

.PHONY: swagger test test-cover bench check

# Unit tests with per-package coverage summary (pos+neg cases in packages).
test:
	go test $(TEST_PKGS) -count=1 -cover

# Library coverage profile + fail if either package is below COVER_MIN%.
test-cover:
	@fail=0; \
	for pkg in $(COVER_PKGS); do \
		out=$$(echo "$$pkg" | tr -d './' | tr '/' '_').cover.out; \
		echo "==> cover $$pkg"; \
		go test "$$pkg" -count=1 -covermode=atomic -coverprofile="$$out" || fail=1; \
		go tool cover -func="$$out" | tail -n 1; \
		pct=$$(go tool cover -func="$$out" | awk '/^total:/ { gsub(/%/,"",$$NF); print $$NF }'); \
		awk -v p="$$pct" -v m="$(COVER_MIN)" -v pkg="$$pkg" 'BEGIN { \
			if ((p+0) < (m+0)) { \
				printf "%s coverage %.1f%% < required %s%% (HA library bar)\n", pkg, p, m; \
				exit 1 \
			} \
			printf "%s coverage %.1f%% OK (HA library bar ≥%s%%)\n", pkg, p, m \
		}' || fail=1; \
	done; \
	echo "==> cmd/testserver (coverage reported; no hard floor)"; \
	go test ./cmd/testserver/... -count=1 -cover || fail=1; \
	exit $$fail

# Hot-path benchmarks with allocs (HA-safe: watch B/op + allocs/op; no storms).
bench:
	@echo "HA-safe benchmarks: review ns/op, B/op, allocs/op (fix regressions / alloc storms)"
	go test $(BENCH_PKGS) -run='^$$' -bench=. -benchmem -count=1

# Full local HA gate: coverage floor + benches with alloc metrics.
check: test-cover bench

# Regenerate OpenAPI/Swagger from Go annotations.
swagger:
	go run github.com/swaggo/swag/cmd/swag@v1.16.4 init \
		-g main.go \
		-d cmd/testserver,cmd/testserver/handlers \
		-o cmd/testserver/docs \
		--parseDependency=false \
		--parseInternal
