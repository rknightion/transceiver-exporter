# v1 Roadmap

Tracking document for getting transceiver-exporter to a 1.0 release. (GitHub
Issues are currently disabled on this repository, so this file stands in for an
issue tracker. Enable **Settings → Features → Issues** to move these to real
issues.)

## Blockers — must fix before v1

- [x] **Metric prefix mismatch.** Emitted metrics used the `transceiver_` prefix
  while the docs used `transceiver_exporter_`. Fixed in code by switching
  `const prefix` to `transceiver_exporter_`. ⚠️ Breaking change for existing
  dashboards/alerts.
- [x] **Stale hardcoded version.** `main.go` fallback was `1.5.1`; changed to
  `dev` (official builds inject the real version via ldflags).
- [x] **Data race in `NewCollector`.** Per-scrape mutation of package-global
  `*prometheus.Desc` values was a data race under concurrent scrapes. All
  descriptors are now built once in `init()`.

## Quality — targeted for v1

- [x] **Unit tests.** Added coverage for `util.go`, `compileRegexFlags`, and
  collector `Describe` (incl. a concurrency test run under `-race`).
- [x] **Documentation.** Corrected metric names, documented the `*_dbm` family,
  `bus_info`, and `laser_supports_monitoring_bool`; added an install / Docker /
  compose section; refreshed maintainer/author attribution.
- [x] **Community-health files.** Added `SECURITY.md`, `CODE_OF_CONDUCT.md`,
  expanded `CONTRIBUTING.md`, plus issue/PR templates.

## Open decisions / follow-ups

- [x] **Architecture support.** Standardised on `linux/amd64` + `linux/arm64`
  for both container images and release binaries. `386` and `arm`/`armv7` were
  dropped from `.goreleaser.yaml` to match the container-publish matrix. This is
  a breaking change for anyone consuming those binaries/images.
- [ ] **Enable GitHub Issues** and close the redundant Mend Bolt onboarding
  PR (#1); decide on the Renovate onboarding PR (#2).
- [ ] **Tag v1.0.0** once the above are resolved (release-please will pick up the
  breaking-change commit and propose the major bump).
