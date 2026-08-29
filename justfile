set shell := ["bash", "-euo", "pipefail", "-c"]

version := env('VERSION', `git describe --tags --always --dirty 2>/dev/null || echo dev`)
commit := env('COMMIT', `git rev-parse HEAD 2>/dev/null || echo unknown`)

# renovate: datasource=go depName=github.com/golangci/golangci-lint/v2
golangci_lint_version := "v2.13.2"
# renovate: datasource=go depName=golang.org/x/vuln
govulncheck_version := "v1.3.0"
# renovate: datasource=go depName=github.com/goreleaser/goreleaser/v2
goreleaser_version := "v2.18.0"

# show the task surface
default:
    @just --list </dev/null

# download Go module dependencies
setup:
    go mod download

# format Go sources and this justfile in place
[group('dev')]
[script('bash')]
fmt:
    gofmt -s -w $(git ls-files '*.go')
    just --fmt

# verify Go and justfile formatting without mutating files
[group('check')]
[no-exit-message]
[script('bash')]
fmt-check:
    files="$(gofmt -l -s $(git ls-files '*.go'))"
    if [ -n "$files" ]; then
        echo "gofmt: files need formatting; run 'just fmt'" >&2
        echo "$files" >&2
        exit 1
    fi
    just --fmt --check </dev/null

# run Go vet and the pinned static-analysis suite
[group('check')]
[no-exit-message]
lint:
    go vet ./...
    go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@{{ golangci_lint_version }} run ./...

# build every package and run the race-enabled test suite; filter narrows by -run
[group('check')]
[no-exit-message]
[script('bash')]
test filter="":
    go build ./...
    if [ -n "{{ filter }}" ]; then
        go test -race -v -run "{{ filter }}" ./...
    else
        go test -race -v ./...
    fi

# scan dependencies for known vulnerabilities
[group('check')]
[no-exit-message]
vuln:
    go run golang.org/x/vuln/cmd/govulncheck@{{ govulncheck_version }} ./...

# write an atomic coverage profile for the informational Codacy upload
[group('check')]
coverage:
    go test -covermode=atomic -coverprofile=coverage.out ./...

# compile the exporter binary into bin/
[group('build')]
build:
    mkdir -p bin
    go build -trimpath -ldflags "-X main.version={{ version }}" -o bin/transceiver-exporter .

# cross-compile a non-publishing GoReleaser snapshot; requires cross-compilation
[group('build')]
snapshot:
    go run github.com/goreleaser/goreleaser/v2@{{ goreleaser_version }} release --snapshot --clean --skip=publish,sign,sbom

# build a local container image; requires a Docker daemon
[group('build')]
image tag="transceiver-exporter:dev":
    docker build --build-arg VERSION="{{ version }}" --build-arg COMMIT="{{ commit }}" -t "{{ tag }}" .

# run the exporter locally; CAP_NET_ADMIN is required to read real EEPROM data
[group('dev')]
run *args:
    go run . {{ args }}

# run every pre-commit check that needs only the Go toolchain
[group('check')]
check: fmt-check lint test vuln

# run check plus snapshot (cross-compilation) and image (Docker daemon) CI legs
[group('check')]
ci: check snapshot image
