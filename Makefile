.PHONY: build lint test deadcode all

BIN := bin/custom-gcl$(if $(filter Windows_NT,$(OS)),.exe,)

# build lazily rebuilds custom-gcl only when .custom-gcl.yml or go.sum
# changed since the last build (mirrors CI's selflint cache key).
build:
	@if [ ! -f "$(BIN)" ] || [ .custom-gcl.yml -nt "$(BIN)" ] || [ go.sum -nt "$(BIN)" ]; then \
		echo "building custom-gcl..."; \
		golangci-lint custom; \
	else \
		echo "custom-gcl is up to date, skipping rebuild"; \
	fi

lint: build
	./$(BIN) run ./...

test:
	go test ./...

# deadcode mirrors CI's exception list for the plugin registration surface
# (see docs/spikes.md): those functions are invoked by golangci-lint's
# external custom-gcl framework, not by cmd/dlinter, and the helpers they
# alone reach (rolegraph.New/classify/Graph.Resolve/Graph.Allowed,
# maydependon's runWithGraph/relativize) share the same false-positive
# class.
deadcode:
	@out="$$(go run golang.org/x/tools/cmd/deadcode@latest ./... 2>&1 | grep -v -E '(plugin\.go:.*unreachable func: (init#1|New|plugin\.BuildAnalyzers|plugin\.GetLoadMode)|analyzer\.go:.*unreachable func: (NewAnalyzer|runWithGraph|relativize)|rolegraph\.go:.*unreachable func: (New|classify|Graph\.Resolve|Graph\.Allowed))')"; \
	if [ -n "$$out" ]; then \
		echo "$$out"; \
		echo "deadcode found unreachable code (see output above)"; \
		exit 1; \
	fi; \
	echo "no dead code found (plugin-entrypoint exceptions documented in docs/spikes.md)"

all: test lint deadcode
