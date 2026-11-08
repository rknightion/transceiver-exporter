# transceiver-exporter

Prometheus exporter that reads pluggable-optic EEPROM through `ethtool`. Forked from
`wobcom/transceiver-exporter`: the module path is `github.com/rknightion/transceiver-exporter` but
the decoding still comes from the upstream `github.com/wobcom/go-ethtool` dependency.

## Task interface

`just check` is the pre-commit gate. CI runs its components as separate jobs and additionally gates
the `just ci` legs, `snapshot` (cross-compilation) and `image` (Docker daemon). `just coverage`
feeds an informational Codacy upload and is deliberately not a required check.

**The gate cannot see the hardware.** The exporter's real work is decoding EEPROM from physical
optics on a Linux host, and there is no such host in CI, no simulator and no captured-EEPROM
fixtures. A green gate on a decoding change proves compilation, static analysis and pure helpers,
not correct decoding. Say which claims are gate-backed and which are reasoned; never collapse the
two. `just run` needs `CAP_NET_ADMIN` to read real EEPROM.

Release and RC binary workflows read the Go version from `go.mod` through the shared binaries reusable; avoid a separate `go-version` override. The runtime image digest includes Debian ca-certificates `20250419+deb13u1`.

## Commits and releases

Conventional Commit subjects are load-bearing: release-please derives the version bump from them,
and a `!` or `BREAKING CHANGE:` produces a major. Never hand-edit `CHANGELOG.md`,
`.release-please-manifest.json` or a version literal in Go source. release-please owns all three, so
a hand edit either loses to it or corrupts its state.

The docs hub in `m7kni/m7kni-net-site` generates several paths into a clone at build time, and they
are gitignored here on purpose: `zensical.toml`, `docs/overrides/`, `docs/stylesheets/brand.css`,
`docs/assets/project-icon*.svg`, `docs/assets/social-card.*` and `docs/fonts/`. Do not create them
locally to "fix" a missing theme.

## Tracking work

Work is tracked in this repository's GitHub issues; they are the only queue. Do not add a checklist
file such as `ROADMAP.md` beside them.

- A pull request body says how it relates to an issue, on a line of its own: `Closes #N` with the
  `Acceptance evidence` table from `.github/pull_request_template.md` (one row for each acceptance
  criterion of that issue), `Refs #N` when it only relates to one, or `Task: none, <reason>`.
- A pull request that changes `.github/workflows/` or `.github/actions/` carries no closing keyword:
  write `Task: none, <reason>`, with `Refs #N` beside it if you like.
- Contributors never need the project board or the tracker tooling; a maintainer handles both.
- Issues and pull requests are public, so they carry no real identifiers: no email addresses,
  handles, account ids, host or device names, addresses or coordinates. Write the shape, not the
  instance.
- Older commits, docs and tests cite legacy board ids such as `TXE-0001`. `just issue <ID>` prints
  the issue each one became; `.github/tracker/id-map.json` is the map.

## Deeper references

- `docs/permissions.md` - read before changing how the exporter obtains EEPROM access or what
  capabilities the container needs.
- `docs/metrics.md` - read before adding, renaming or relabelling a metric.
