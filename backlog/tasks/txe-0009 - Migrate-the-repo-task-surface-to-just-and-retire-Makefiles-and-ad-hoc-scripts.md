---
id: TXE-0009
title: Migrate the repo task surface to just and retire Makefiles and ad-hoc scripts
status: To Do
assignee: []
created_date: '2026-08-28 19:27'
updated_date: '2026-08-29 10:57'
labels:
  - 'wave:2-fleet'
dependencies: []
priority: medium
type: chore
ordinal: 9000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
# Migrate task surface to just

## 1. Outcome

`transceiver-exporter` gains a top-level `justfile` implementing the fleet-mandatory recipe
vocabulary (`default`, `setup`, `fmt`, `fmt-check`, `lint`, `test`, `check`) plus `build` and `run`.
`just --list` is the one true answer to "what can I do in this repo". CI's `tests` job in
`.github/workflows/ci.yml` calls `just check` instead of inlining `go test -v ./...` and the
golangci-lint-action. `AGENTS.md`'s "The gate" section and `backlog/config.yml`'s
`definition_of_done` both name `just` recipes instead of raw `go`/`golangci-lint` invocations.

This repo has **no Makefile and no tracked shell/helper scripts** (verified: `find . -iname
Makefile -o -iname GNUmakefile` and `git ls-files | grep -E '\.(sh|bash|zsh|ps1)$'` both return
empty, excluding `vendor/`). There is nothing to delete and nothing to classify ABSORB/KEEP. This
is a pure "author the justfile from scratch and wire it into CI/docs" task — do not go looking for
a Makefile or scripts that do not exist.

## 2. The complete justfile

Create `justfile` at repo root with exactly this content (adjust only if `go.mod`'s Go version or
`golangci-lint`'s pinned version has since changed):

```just
set shell := ["bash", "-euo", "pipefail", "-c"]

# show the task surface
default:
    @just --list

# install toolchain + deps into the repo-local environment
[group('dev')]
setup:
    go mod download

# format source in place
[group('check')]
fmt:
    gofmt -l -s -w .
    just --fmt

# verify formatting without mutating
[group('check')]
fmt-check:
    @test -z "$(gofmt -l -s .)" || (gofmt -l -s . && echo 'gofmt: files need formatting, run `just fmt`' >&2 && exit 1)
    just --fmt --check

# static analysis
[group('check')]
[no-exit-message]
lint:
    go vet ./...
    golangci-lint run

# run the full test suite (race + verbose); optional substring filter
[group('check')]
[no-exit-message]
test filter="":
    #!/usr/bin/env bash
    set -euo pipefail
    if [ -n "{{filter}}" ]; then
        go test -race -v -run "{{filter}}" ./...
    else
        go test -race -v ./...
    fi

# build the binary into bin/
[group('build')]
build:
    go build -trimpath -o bin/transceiver-exporter .

# run the exporter locally (needs CAP_NET_ADMIN to read real EEPROM data)
[group('dev')]
run *args:
    go run . {{args}}

# the full PR gate — exactly what CI enforces
[group('check')]
check: fmt-check lint test
```

Notes for the implementing agent:
- `fmt-check`'s `gofmt -l -s .` line lists unformatted files to stderr before failing, so the
  agent sees exactly which files are dirty — matches `[no-exit-message]` convention used elsewhere
  in the fleet standard for tool output that's already good; `fmt-check` itself doesn't carry
  `[no-exit-message]` because the custom message names the fix (`just fmt`).
  `just --fmt --check` is included per §5.10 — mandatory.
- `test` uses `[script('bash')]`-style shebang recipe (the `#!/usr/bin/env bash` line) rather than
  a plain line body, because it needs `if`/`else` control flow, which the §10 "extra leading
  whitespace" gotcha rules out for line-based recipes.
- No `typecheck` recipe — Go's compiler IS the typecheck, already covered by `go vet` (lint) and
  `go build`/`go test` (build/test). Do not invent a redundant `typecheck` recipe.
- No `gen`/`gen-check` — nothing in this repo is code-generated (verified: no `//go:generate`
  directives, `grep -rn "go:generate" *.go transceiver-collector/*.go` returns empty; confirm this
  still holds before skipping).
- `docs/` is a mkdocs/zensical site built by the external `m7kni-net-site` hub, not by this repo
  (see `docs.toml` header comment) — no `docs`/`docs-serve` recipe belongs here.
- `ci` recipe is deliberately omitted: CI's `tests` job (build/lint/test) is a full subset of
  `check`; the `docker-build-verify` and `coverage` jobs are GitHub-native concerns (buildx,
  Codacy upload) that do not belong in `just` per §8 — they stay as separate workflow steps, not
  folded into a `ci` superset recipe.
- `setup` is intentionally minimal (`go mod download`) — this repo has no other dev-time
  dependency (no linter installer needed locally; golangci-lint is expected preinstalled or via
  `go run` — see Traps §9 below for the CI-vs-local asymmetry this creates).

## 3. Makefile disposition

None. No `Makefile` or `GNUmakefile` exists anywhere in this repo (verified, excluding `vendor/`).
No `git rm` step needed for this section.

## 4. Script disposition

None. `git ls-files | grep -E '\.(sh|bash|zsh|ps1)$'` returns empty. No `scripts/` directory
exists. No ABSORB or KEEP classification needed.

## 5. CI changes

### `.github/workflows/ci.yml`

Only the `tests` job changes. `docker-build-verify`, `coverage`, and `ci-success` are untouched —
do not touch the `needs: [tests, docker-build-verify]` list on `ci-success`, `permissions:`,
`concurrency:`, `persist-credentials: false`, or any SHA-pinned `uses:`.

**Before** (`tests` job steps, lines 19–35 of the current file):
```yaml
      - uses: actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1 # v7.0.1
        with:
          persist-credentials: false

      - uses: actions/setup-go@b7ad1dad31e06c5925ef5d2fc7ad053ef454303e # v7.0.0
        with:
          go-version-file: go.mod

      - name: Run linters
        uses: golangci/golangci-lint-action@ba0d7d2ec06a0ea1cb5fa41b2e4a3ab91d21278a # v9.3.0
        with:
          version: latest
          args: --verbose

      - name: Run tests
        run: go test -v ./...
```

**After:**
```yaml
      - uses: actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1 # v7.0.1
        with:
          persist-credentials: false

      - uses: actions/setup-go@b7ad1dad31e06c5925ef5d2fc7ad053ef454303e # v7.0.0
        with:
          go-version-file: go.mod

      - name: Install golangci-lint
        uses: golangci/golangci-lint-action@ba0d7d2ec06a0ea1cb5fa41b2e4a3ab91d21278a # v9.3.0
        with:
          version: latest
          install-mode: binary

      - uses: extractions/setup-just@<pin-current-sha-here> # v4
        with:
          just-version: '1.58.0'

      - name: Run gate
        run: just check
```

**Do not** delete the `golangci-lint-action` step outright — it is the fleet's existing mechanism
for pinning/caching the golangci-lint binary in CI. Switch it to `install-mode: binary` (installs
the tool onto `PATH` without running the lint itself) so `just check`'s `lint` recipe can then
invoke the already-installed `golangci-lint run`. Resolve `extractions/setup-just`'s pinned SHA at
implementation time (`gh api repos/extractions/setup-just/git/refs/tags/v4 --jq .object.sha` or
check what SHA other already-migrated fleet repos are using) and add the matching `# v4` comment
per the fleet's SHA-pin convention.

Everything else in `ci.yml` — `docker-build-verify`, `coverage`, `ci-success` — is unchanged.

### Other workflow files — no changes

`actionlint.yml`, `arm-automerge.yml`, `auto-rc.yml`, `codeql.yml`, `dependency-review.yml`,
`docker-security.yml`, `ghcr-cleanup.yml`, `publish.yml`, `release-please.yml`, `scorecard.yml`,
`trigger-docs-sync.yml`, `zizmor.yml` — none contain build/test/lint/format/generate `run:` bodies
that belong in `just`. Verified: `grep -rn "make \|\.sh" .github/workflows/` returns empty across
the whole directory. Do not touch these files.

## 6. Docs and agent-contract changes

### `AGENTS.md` — "The gate" section

Current (lines 9–20 approx):
```markdown
## The gate

\`\`\`bash
go build ./...
go vet ./...
go test -race ./...
golangci-lint run
\`\`\`

These are the four items every new task inherits as its definition of done. CI additionally
verifies the Docker build; `coverage` is deliberately not a required check.
```

Replace with:
```markdown
## The gate

This repo's task surface is a `justfile`. Discover it, don't guess it:

\`\`\`bash
just --list                        # human-readable
just --dump --dump-format json     # machine-readable
just --show <recipe>               # what a recipe actually runs
\`\`\`

- `just check` is the full gate and is exactly what CI enforces. It must pass before you commit.
- Prefer `just <recipe>` over the underlying tool. If you are typing `go test`, you want
  `just test`.
- Run `just` with stdin from /dev/null. This repo has no `[confirm]` recipes today, but if one is
  added later, stop and ask before running it — never pass `--yes` or `JUST_YES=1`.
- If a task you need does not exist, add a recipe with a `#` doc comment and a `[group(...)]`
  rather than running a bare command.

CI additionally verifies the Docker build; `coverage` is deliberately not a required check.
```

Do **not** paste the actual recipe list (`build`, `test`, `lint`, etc.) into `AGENTS.md` — only
the contract paragraph above, per §9 of the fleet standard. `CLAUDE.md` needs no edit — it's a
one-line `@AGENTS.md` import and picks the change up automatically.

### `CONTRIBUTING.md`

`grep -n "make \|\.sh" CONTRIBUTING.md` returns only "Make your change, adding or updating tests
where it makes sense." (line 33) — plain English "make", not a `make` command reference. No edit
needed there. Re-check at implementation time in case this has drifted.

### `README.md`

No `make` or script-path references found (`grep -n "make\|\.sh\b" README.md` returned empty). No
edit needed. Re-check at implementation time.

## 7. `backlog/config.yml`

Current:
```yaml
definition_of_done:
  - "go build ./..."
  - "go vet ./..."
  - "go test -race ./..."
  - "golangci-lint run"
```

New:
```yaml
definition_of_done:
  - "just check"
```

This file is the one Backlog.md markdown-adjacent file the fleet standard permits hand-editing
(list-valued keys can't be set through `backlog config set`) — see `AGENTS.md`'s own note on this.
Edit it directly with a text editor, not the `backlog` CLI.

## 8. Order of work

1. Add `justfile` at repo root (§2). Run `just --fmt --check`, then `just check` locally end to
   end. Fix anything that doesn't pass before touching CI.
2. Update `.github/workflows/ci.yml`'s `tests` job (§5). Push and confirm the `tests` job goes
   green and `ci-success` still gates on the same two job names.
3. Update `AGENTS.md`'s "The gate" section (§6).
4. Update `backlog/config.yml`'s `definition_of_done` (§7).
5. There is nothing to delete (no Makefile, no scripts) — skip any "delete last" step.

## 9. Traps specific to this repo

- **`golangci-lint-action`'s dual role.** The action both installs the binary AND runs the lint by
  default. Switching CI to `install-mode: binary` stops it from running the lint itself — if you
  only add the `setup-just` step and leave the action in its default mode, CI lints **twice**
  (once via the action, once via `just check`), which isn't wrong but wastes CI minutes and can
  produce two different diagnostic formats in the log. Set `install-mode: binary` as shown in §5.
- **`lint` requires golangci-lint on `PATH` locally too**, and this repo has no `.golangci.yml` —
  it runs with golangci-lint's built-in defaults. A contributor without golangci-lint installed
  gets a `command not found` from `just lint` / `just check`. This is a pre-existing condition
  (the old `AGENTS.md` gate already assumed `golangci-lint` was on PATH) — the justfile does not
  need to fix it, but do not add a `require('golangci-lint')` guard that changes today's behavior
  without discussing it; a bare command-not-found is the same failure mode as before.
- **`test`'s `-race` flag needs CGO** (`CGO_ENABLED=1`, the Go default on most dev machines) —
  unlike the Dockerfile's `CGO_ENABLED=0` static build. If `just test` ever fails with a "race
  detector requires cgo" style error on a from-scratch container, that's an environment gap, not a
  justfile bug — do not silently drop `-race` to work around it, since `AGENTS.md`'s existing gate
  already mandates `go test -race`.
- **`run` needs `CAP_NET_ADMIN`** to read real EEPROM data (see `compose.yml`'s comment block) —
  `just run` without elevated privileges/capabilities will start but report no transceivers, not
  error. Don't read that as a bug when smoke-testing the recipe.
- **No generated-file drift gate applies here** — confirmed no `//go:generate` directives exist,
  so don't add a `gen`/`gen-check` pair speculatively.
- **`bin/` from `just build`** is not currently gitignored-checked in this task — verify
  `bin/` (or the exact output path) is in `.gitignore` before merging; if it's missing, add it as
  a small adjacent fix (out of scope to invent a different build output path).

## 10. Out of scope

- **Every other workflow file**: `actionlint.yml`, `arm-automerge.yml`, `auto-rc.yml`,
  `codeql.yml`, `dependency-review.yml`, `docker-security.yml`, `ghcr-cleanup.yml`, `publish.yml`,
  `release-please.yml`, `scorecard.yml`, `trigger-docs-sync.yml`, `zizmor.yml` — GitHub-native,
  untouched, per §8 of the fleet standard.
- **`docker-build-verify` and `coverage` jobs in `ci.yml`** — stay as-is; buildx and the Codacy
  coverage upload are not `just` concerns.
- **`ci-success`'s `needs:` list, `permissions:`, `concurrency:`, SHA pins,
  `persist-credentials: false`** — do not touch.
- **`Dockerfile` and `compose.yml`** — the Dockerfile's own multi-stage `RUN go build …` line is a
  container build step baked into an image layer, not a developer/CI task; it is not converted to
  `just build` (it runs inside `docker build`, which has no `just` available at that stage).
- **`docs/` content and `docs.toml`** — owned by the external `m7kni-net-site` hub build, not this
  repo's task surface.
- **`release-please-config.json`, `.release-please-manifest.json`, `CHANGELOG.md`** —
  release-please owns these; do not add a `release` recipe that touches them (`[[release-please-pat]]`
  fleet convention: releases are automated, not developer-triggered).
- **`renovate.json`** — untouched.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Top-level justfile exists implementing default, setup, fmt, fmt-check, lint, test, check plus build and run, matching the fleet vocabulary in the task body
- [ ] #2 just check passes locally and runs exactly fmt-check, lint, test — the same checks the CI tests job in .github/workflows/ci.yml enforces
- [ ] #3 just --fmt --check passes on the justfile
- [ ] #4 just --list shows a # doc comment and a [group(...)] for every public recipe
- [ ] #5 No Makefile or GNUmakefile exists in the repo (none existed pre-migration; this confirms none was introduced)
- [ ] #6 No tracked shell/helper script exists that duplicates a just recipe (none existed pre-migration; this confirms none was introduced)
- [ ] #7 .github/workflows/ci.yml's tests job calls just check via a setup-just step instead of inlining go test and golangci-lint-action's default lint mode, while ci-success's needs list, permissions, concurrency and persist-credentials stay unchanged
- [ ] #8 AGENTS.md's 'The gate' section documents just --list / just --dump / just --show discovery and states just check is exactly what CI enforces, with no raw go/golangci-lint commands and no pasted recipe list
- [ ] #9 backlog/config.yml's definition_of_done is exactly ["just check"]
- [ ] #10 CONTRIBUTING.md and README.md contain no make or ./scripts/ references (re-verified at implementation time)
<!-- AC:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [ ] #1 go build ./...
- [ ] #2 go vet ./...
- [ ] #3 go test -race ./...
- [ ] #4 golangci-lint run
<!-- DOD:END -->

## Comments

<!-- COMMENTS:BEGIN -->
author: campaign-ordering
created: 2026-08-29 09:18
---
## Fleet ordering — WAVE 2. Starts after the Wave 0 pilot (`sf2loki` / SFL-0073) and the Wave 1 hubs land.

Within Wave 2 the order is free — these repos do not depend on each other. Batching by language is worthwhile so one lane reuses its Makefile-to-recipe mapping across similar repos.

Do not start before the pilot reports. The standard may be amended off the back of it, and picking this up early risks coding against a superseded seam.

**Provisioning `just` in CI.** Which mechanism depends on the runner, and the two must not be mixed:

| Runner | Mechanism |
| --- | --- |
| `arc-arm64` (m7kni self-hosted) | `just` is **baked into the runner image** by `m7kni/ci-tools` (`runner-image/Dockerfile`, `ARG JUST_VERSION`). Do **not** add `extractions/setup-just`, and delete the step if this repo already has one — it installs a second `just` earlier on `PATH` and turns the image pin into a lie. |
| GitHub-hosted (all `rknightion` repos) | `extractions/setup-just`, SHA-pinned, with an explicit `just-version:`. |

Both sides currently sit on **1.58.0** and are Renovate-managed. `ci-tools`' `Tool version drift` workflow fails if the Dockerfile `ARG` and the published image ever disagree, and lists any repo still carrying a second pin.

**While you are in the workflow files, check the hub pin.** On 2026-08-29 Renovate was unfrozen for `rknightion/.github` in `m7kni/renovate-config` — it had been `enabled: false` on the mistaken belief that callers tracked `@main`, which froze the fleet across 19 different hub SHAs (v1.3.1 June → v1.9.7 August) so that no hub fix ever propagated. Bumps now arrive as one grouped, CI-gated, automerged PR per repo. **A `uses:` whose comment is not a real `# vX.Y.Z` still cannot be bumped** (it resolves to a digest-only update, which the fleet rules disable) — if you find one, repair the comment as part of this task.
---

author: campaign-ordering
created: 2026-08-29 10:43
---
## Standard amendment — `ci` is the sanctioned superset of `check` (RATIFIED)

This supersedes the frozen wording *"`check` is the complete local gate and reproduces every CI job that can run off a GitHub runner"*, which several lanes could not honour without making the pre-commit gate depend on a Docker daemon.

**The definitions now are:**

- **`check`** — everything that runs with **only the language toolchain installed**. This is the pre-commit gate. A leg that runs on a bare toolchain belongs here *however long it takes*.
- **`ci`** — `check` plus the legs CI gates that need a **Docker daemon, a service container, or cross-compilation**, and nothing else. Written as `ci: check <heavy legs>`.

**Every leg you put in `ci` must carry a comment naming which of those three it needs.** That comment is the guard: without it `ci` becomes the bin for anything slow or awkward, `check` quietly stops meaning much, and the fleet is back to a per-repo gate.

Eleven of the 42 lanes arrived at this shape independently before it was ratified, which is why it won.

**If this repo has no such legs, it has no `ci` recipe at all** and `check` is the whole gate. Do not add an empty one.
---

author: campaign-ordering
created: 2026-08-29 10:57
---
## Fleet alignment — the 2otel family converges on one CI shape

These seven Go repos are near-identical applications and had drifted into **two naming dialects and materially different coverage**. The migration rewrites every `run:` block anyway, so converge them in the same change rather than preserving the drift in new clothes.

**Canonical job names** — used by tailscale2otel, graph2otel, polylens2otel and rfc6035-2otel, so this is the majority convention, not an invention:

`build-test` · `lint` · `govulncheck` · `goreleaser-snapshot` · `docker-build` · `coverage` · `ci-success`

`opnsense2otel` and `transceiver-exporter` currently use a second dialect — `tests`, `race`, `docker-build-verify`. Rename to the canonical set as part of this task.

**`ci-success` is the only check the branch ruleset gates**, so jobs can be renamed or merged freely *provided* `ci-success`'s `needs:` list is updated in the same commit. Never rename `ci-success` itself.

**Required gates, and where each lives after the migration:**

| Gate | Recipe | Note |
| --- | --- | --- |
| build + test + `-race` | `just test` | `-race` belongs in the standard test run |
| golangci-lint | `just lint` | needs a `.golangci.yml`, schema v2 |
| **gosec** | `just lint` | **a golangci-lint linter, NOT a separate job** — enable it in `.golangci.yml`. Four of the seven already do it this way; a standalone gosec job would be a third dialect |
| govulncheck | `just vuln` | pinned `golang.org/x/vuln/cmd/govulncheck@v1.3.0`, matching the family |
| goreleaser snapshot | `just snapshot` | cross-compile ⇒ belongs in `ci`, not `check` |
| container build | `just image` | needs a Docker daemon ⇒ belongs in `ci`, not `check` |

**Already done for you (2026-08-29):** `govulncheck` was added to `opnsense2otel`, `transceiver-exporter` and `codexlb2otel` ahead of the migration, because those three had no dependency vulnerability scanning at all. Convert those jobs to `just vuln` like any other; do not re-add them.

**Still missing, fix as part of this task:**

- `opnsense2otel` — has `.golangci.yml` but **`gosec` is not enabled** in it.
- `transceiver-exporter` — **no `.golangci.yml` at all**, and no `-race` in its test job.
- `codexlb2otel` — no `.golangci.yml`, no `-race`, no container build, and **no `ci-success` job and no branch ruleset**, so nothing gates its CI. Adding an aggregator is the right fix but is a separate decision; raise it rather than assuming.

**One known trap:** the `govulncheck@v1.3.0` pins are invisible to Renovate — `go install pkg@version` inside a `run:` block matches no manager. All five are four minor versions behind (current is v1.7.0). Once the version moves into the justfile as a `# renovate:`-annotated `:=` assignment, it becomes managed. That is a real benefit of this migration, not incidental.
---
<!-- COMMENTS:END -->
