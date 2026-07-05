# Contributing

Thanks for your interest in improving transceiver-exporter!

## Reporting bugs & requesting features

Open an issue on GitHub. For larger changes, please open an issue to discuss the
approach before sending a pull request.

## Development

The project is a small, dependency-light Go program.

```bash
# Build
go build ./...

# Run the test suite (with the race detector — CI does this too)
go test -race ./...

# Vet and lint
go vet ./...
golangci-lint run
```

Note: the exporter reads real transceiver EEPROM data via `ethtool`, so a full
end-to-end run requires a Linux host with pluggable optics. The unit tests do
not require any hardware.

## Pull requests

1. Fork the repository and create a branch off `main`.
2. Make your change, adding or updating tests where it makes sense.
3. Ensure `go build ./...`, `go vet ./...`, and `go test -race ./...` all pass.
4. Use [Conventional Commits](https://www.conventionalcommits.org/) for commit
   messages (`feat:`, `fix:`, `docs:`, `chore:`, …). Releases and the changelog
   are generated automatically from commit messages via release-please, so this
   matters.
5. Open a pull request against `main` and fill in the description.

## Code of Conduct

This project follows the [Contributor Covenant](./CODE_OF_CONDUCT.md). By
participating you are expected to uphold it.
