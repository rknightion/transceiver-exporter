# Loop: transceiver-exporter
tier: guarded
gate: just check
ci-required: ci-success
grafana-stack: none

`backlog/` is committed: no hostname, address or email in any task or doc (sweep in `AGENTS.md`). Run
`just` with stdin from `/dev/null`.

## Credentials

- None held by the loop. The dev machine is macOS and CI has no optics, so no lane can read a real
  EEPROM.

## Traps

- The gate proves compilation, typing, lint and the pure helpers, not decoded optics data. A green
  `go test -race ./...` on an EEPROM decode change is not evidence it is correct. A task whose
  acceptance needs real optics is Parked with what was implemented, what could not be observed and the
  exact command a human with hardware would run; never redefine acceptance downward. A captured EEPROM
  byte-dump with its expected decode is a real contribution. A return states which claims are
  gate-backed and which are reasoned.
- Descriptors in `transceiver-collector/collector.go` are package globals built once in `init()`: a
  new metric is a new global, never assigned from `Collect` or `NewCollector`, never dependent on
  per-scrape state. The concurrency test in `collector_test.go` keeps running under `-race`.
- A change to a metric name, label set or help text is a two-file change with `docs/metrics.md`, in
  the same commit. `const prefix = "transceiver_exporter_"` is the only source of the prefix.
- Unit pairs are added whole: laser power in both milliwatts and dBm (`milliwattsToDbm` in `util.go`),
  and alarm/warning high/low thresholds in fours.
- No release number in Go source, `CHANGELOG.md` or `.release-please-manifest.json`: release-please owns
  all three and ldflags inject the version (`dev` fallback).
- The image runs as root on purpose (EEPROM access through `ethtool` needs it). Revisit only via
  capabilities tested on hardware.
- Do not create the hub-injected docs paths (`zensical.toml`, `docs/overrides/`, brand CSS, icons, social
  card, `docs/fonts/`): the `m7kni-net-site` hub generates them and they are gitignored on purpose.

## Mutexes

- `transceiver-collector/collector.go` has one owning lane per run; it cannot be split. The Go lane
  (`main.go`, `transceiver-collector/`) also owns `docs/metrics.md`. Disjoint lanes: docs
  (`docs/`, `README.md`, `docs.toml`, except `docs/metrics.md`) and CI or supply chain (workflows,
  `Dockerfile`, `.goreleaser.yaml`, `renovate.json`).
