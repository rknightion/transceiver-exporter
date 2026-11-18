# Transceiver exporter review guide

## What this is

This Go service exports Prometheus metrics for pluggable optics on Linux. `main.go` handles HTTP, scrape requests and interface-selection flags. `transceiver-collector/` reads optical diagnostics and driver data through ethtool.

## Invariants

- Preserve metric names, labels, units and help text. Apply a deliberate metric contract change to `docs/metrics.md` in the same change. Use the single `transceiver_exporter_` prefix from `transceiver-collector/collector.go`.
- Initialise shared Prometheus descriptors once. Preserve request-local collectors and registries so concurrent scrapes and differing power-unit configurations cannot mutate shared descriptors.
- Keep mW and dBm descriptors paired and high/low alarm and warning thresholds complete. `util_test.go` covers the conversion, including zero power.

## Where to look hardest

- Trace interface inclusion, exclusion and administrative-down filters at the boundary between `main.go` flags and `getMonitoredInterfaces`. Preserve partial scrapes when a selected interface cannot be read, and close each ethtool handle after collection.
- Review changes to the HTTP server and exporter handler for timeout and concurrent scrape behaviour. `main_test.go` covers flag compilation and connection timeout configuration. `collector_test.go` exercises concurrent descriptor collection and unit selection.
- Treat changes to EEPROM decoding, host networking, runtime identity or `NET_ADMIN` as hardware-sensitive. Review the `Dockerfile`, `compose.yml` and permission documentation together, with real optics or a captured EEPROM fixture and expected decode for hardware claims.

## Checks already covered

Use the existing Go formatting, vet, golangci-lint, race, govulncheck and dependency review results for their covered findings. Actions workflows also check workflow syntax, Zizmor, CodeQL and container security.

## Verification

Run `just setup` and `just check </dev/null>`. CI also builds a release snapshot and the container. Its software checks cover the pure helpers, compilation, descriptor contract and races. Validate a decoding claim with matching hardware or a captured EEPROM fixture and expected results.
