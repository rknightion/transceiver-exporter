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

## Commits and releases

Conventional Commit subjects are load-bearing: release-please derives the version bump from them,
and a `!` or `BREAKING CHANGE:` produces a major. Never hand-edit `CHANGELOG.md`,
`.release-please-manifest.json` or a version literal in Go source. release-please owns all three, so
a hand edit either loses to it or corrupts its state.

The docs hub in `m7kni/m7kni-net-site` generates several paths into a clone at build time, and they
are gitignored here on purpose: `zensical.toml`, `docs/overrides/`, `docs/stylesheets/brand.css`,
`docs/assets/project-icon*.svg`, `docs/assets/social-card.*` and `docs/fonts/`. Do not create them
locally to "fix" a missing theme.

## Task tracking

The Backlog board is the only queue. Do not reintroduce a checklist file such as `ROADMAP.md`
alongside it.

- `backlog/` is committed, so tasks, docs and decisions carry no real identifiers: no email
  addresses, handles, account IDs, host or device names, addresses or coordinates. Write the shape,
  not the instance; aggregate counts, timings and structural findings are fine. Sweep before
  committing:

      grep -rniE "rob-knight\.|@gmail|[0-9]{1,3}(\.[0-9]{1,3}){3}" backlog/ && echo "PII FOUND"

- `backlog/config.yml` is the one file exempt from driving the tracker through its CLI, because
  list-valued keys cannot be set through `backlog config set`.
- Finalize in one call, so an interrupted session cannot leave finished work looking unfinished:
  `backlog task edit TXE-0007 --check-ac 1 --check-ac 2 -s Done`.
- Never let two agents edit the same task. The upstream concurrent-edit fix covers the edit funnel
  but not reorder, draft saves, the TUI path, `doc update` or decision updates.

Read the `Agent fan-out protocol (canonical)` doc before designing a wave, and `Wave operating model`
for this project's lane conventions, ownership, escape hatch and the recurring defects the hardware
boundary above produces. Both are in `backlog doc list --plain`.

## Deeper references

- `docs/permissions.md` - read before changing how the exporter obtains EEPROM access or what
  capabilities the container needs.
- `docs/metrics.md` - read before adding, renaming or relabelling a metric.

<!-- BACKLOG.MD GUIDELINES START -->
<!-- backlog.md-instructions-version: 1.50.1 -->
<CRITICAL_INSTRUCTION>

## Backlog.md Workflow

This project uses Backlog.md for task and project management.

**For every user request in this project, run `backlog instructions overview` before answering or taking action.**

Use the overview to decide whether to search, read, create, or update Backlog tasks.

Before task lifecycle actions, read the matching detailed guide:
- `backlog instructions task-creation` before creating or splitting tasks
- `backlog instructions task-execution` before planning, changing status or assignee, adding a plan or implementation notes, or implementing task work
- `backlog instructions task-finalization` before checking acceptance criteria, writing final summaries, or moving tasks to terminal statuses

Use `backlog <command> --help` before running unfamiliar commands. Help shows options, fields, and examples.

Do not edit Backlog task, draft, document, decision, or milestone markdown files directly. Use the `backlog` CLI so metadata, relationships, and history stay consistent.

</CRITICAL_INSTRUCTION>
<!-- BACKLOG.MD GUIDELINES END -->
