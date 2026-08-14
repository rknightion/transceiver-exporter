---
id: doc-0002
title: Wave operating model
type: guide
created_date: '2026-08-14 16:38'
updated_date: '2026-08-14 16:38'
---
This document carries **only what is specific to `transceiver-exporter`**. The campaign model itself —
run contract and run modes, the routing contract, authority and the thread pool, child lane briefs,
external-contract freezing, the structural failure patterns, the unattended blocker contract, the
goal-file template and the pre-flight checklist — lives in the **Agent fan-out protocol (canonical)**
document. Read that first. Nothing here restates it; if a section below could be pasted into another
project unchanged, it is in the wrong document.

## The exclusive resource is hardware, and no agent has it

The exporter reads transceiver EEPROM through `ethtool` ioctls against real pluggable optics on a
Linux host. **No lane can reach that.** The dev machine is macOS, CI is a GitHub `ubuntu-24.04`
runner with no optics, and there is no simulator, no capture corpus and no fixture set of real EEPROM
dumps in the repo. `CONTRIBUTING.md` says as much: the unit tests require no hardware *because they
cannot exercise the part that matters*.

Consequences, and they are binding:

- **The gate proves compilation, typing, lint and the pure helpers. It proves nothing about decoded
  optics data.** A green `go test -race ./...` on a change to EEPROM decoding is not evidence the
  change is correct, and a lane must not report it as though it were.
- **Any task whose acceptance genuinely requires real optics is `Parked`, not `Done`** — parked with
  the concrete boundary: what was implemented, what could not be observed, and the exact command a
  human with hardware would run to settle it. A lane that quietly redefines acceptance downward to
  something the runner can check has produced the worst outcome available here: a wrong change
  wearing a green tick.
- The honest escape from this is **fixtures**. A task that adds a captured EEPROM byte-dump plus its
  expected decode converts an unverifiable area into a verifiable one permanently, and is worth more
  than the feature that motivated it. There are none today; the first one is a real contribution.

## Recurring defects in this codebase

Each of these has actually happened here, and each is shaped to happen again.

**Per-scrape mutation of package-global descriptors.** `transceiver-collector/collector.go` declares
~60 package-level `*prometheus.Desc` variables and builds every one of them in `init()`. They were
previously rebuilt or mutated during a scrape, which is a genuine data race under concurrent
`/metrics` requests (fixed in `36800a4`). **A new metric is a new package-global built once in
`init()`.** Never assign to a descriptor from `Collect` or from `NewCollector`, and never make a
descriptor depend on per-scrape state. The concurrency test in `collector_test.go` exists to catch
this and must keep running under `-race`.

**Metric names drift between the code and the published docs.** `const prefix = "transceiver_exporter_"`
at `transceiver-collector/collector.go:16` is the only source of the prefix, but `docs/metrics.md`
restates every full metric name by hand. They have already disagreed once — emitted metrics used
`transceiver_` while the docs promised `transceiver_exporter_`, and fixing it was a **breaking change
for every existing dashboard and alert**. So: **any change to a metric name, label set or help text is
a two-file change**, code and `docs/metrics.md`, in the same commit. A lane that touches one and not
the other has not finished.

**Unit-pair metrics get added by half.** Laser power is exposed in both milliwatts and dBm, converted
by `milliwattsToDbm` in `transceiver-collector/util.go`. The descriptor names carry the unit as a
suffix (`…DescMw`) and the `_dbm` family mirrors it. Adding a power metric in one unit and not the
other compiles, lints and tests clean, and is wrong. The same applies to the alarm/warning
high/low **threshold quartets** — they come in fours, and three of four is the failure shape.

**Version strings hardcoded in source.** `main.go` once carried a literal `1.5.1` fallback that
outlived the release it named. It is now `dev`; real builds inject the version through ldflags. **Do
not put a release number in Go source, `CHANGELOG.md`, or `.release-please-manifest.json` by hand** —
release-please owns all three, derives the bump from Conventional Commit subjects, and a hand edit
either loses to it or corrupts its state. A `!` or `BREAKING CHANGE:` in a commit subject is what
produces a major bump, and one is already pending as an open release PR.

**Container-runs-as-root looks like a defect and is not.** The image deliberately runs as root
because EEPROM access through `ethtool` needs it. A lane doing security hardening will find this and
try to "fix" it; that change breaks the exporter's only function. If it is to be revisited, the route
is capabilities (`CAP_NET_ADMIN`/`CAP_NET_RAW`) actually tested against hardware — which, per the
section above, no lane can do.

## Lane conventions

**`transceiver-collector/collector.go` is a single-owner file and cannot be split.** It is one 544-line
file holding every descriptor declaration, its `init()`, `Describe` and `Collect`. Two lanes editing
it in parallel is not a merge problem, it is a semantic one: both will add descriptors to the same
`var` block and the same `init()`. One lane owns the collector per wave, full stop.

The three genuinely disjoint lanes here:

- **Go** — `main.go`, `transceiver-collector/`. One owner while a wave is running.
- **Docs** — `docs/`, `README.md`, `docs.toml`. Disjoint from Go *except* for `docs/metrics.md`,
  which is co-owned with the collector by the two-file rule above; when a wave changes metrics, the
  Go lane owns `docs/metrics.md` and the docs lane does not touch it.
- **CI / supply chain** — `.github/workflows/`, `Dockerfile`, `.goreleaser.yaml`, `renovate.json`.

**Do not create the hub-injected paths.** The published site is built by the `rknightion/m7kni-net-site`
hub from `docs/` plus `docs.toml`; the hub generates `zensical.toml` and copies in
`docs/overrides/`, `docs/stylesheets/brand.css`, the project icons, the social card and `docs/fonts/`.
All of those are gitignored on purpose. A lane that "fixes the missing theme files" by creating them
commits drift that the manifest model exists to prevent.

## Ownership and the escape hatch

A lane owns the files its brief names and nothing else. When a lane finds that finishing its task
requires a change outside its boundary:

**It stops and returns the question — it does not reach across.** The reach-across is what produces
two lanes silently editing one file. The return names the file, the change needed, and whether the
lane is blocked or merely degraded without it.

The root then either widens that lane's boundary explicitly, files the out-of-boundary change as its
own task, or serialises it into the wiring pass. **A boundary with no escape hatch is a stop
condition wearing a safety label** — if a lane cannot ask, it will either guess or stall, and both are
worse than one round-trip.

The root owns every commit. Lanes never commit, never push, and never touch `git` state.

## Run-end against this tracker

Task state is the record, so a run ends when the board is true:

- landed work is `Done`, with the commit SHA in its final summary and its definition-of-done
  checklist checked in the **same** `backlog task edit` call that sets the status;
- attempted-and-blocked work is `Parked` with a concrete resume boundary — for this repo that most
  often means "implemented, unverifiable without optics", and it must say which command settles it;
- untouched work is self-evidently still `To Do` and needs no ceremony;
- work discovered mid-run is a new task labelled `needs-triage`, never a note buried in another
  task's notes;
- anything learned that no single task owns goes in this document or the domain reference, not in
  the closing terminal message. The closing message is a covering note, and nothing durable may live
  only there.

Because this repo's gate cannot see the hardware, the run-end report must state **which claims are
gate-backed and which are reasoned**. "Tests pass" and "the decode is correct" are different
statements here, and collapsing them is the specific dishonesty this repo is prone to.
