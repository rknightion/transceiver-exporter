---
title: Security
description: Reporting vulnerabilities, and the privilege and network exposure this exporter carries.
---

# Security

## Reporting a vulnerability

Do **not** open a public GitHub issue for a security vulnerability. Report it privately via
GitHub's [private vulnerability reporting](https://github.com/rknightion/transceiver-exporter/security/advisories/new)
("Report a vulnerability" under the repository's **Security** tab), including a description of the
impact, reproduction steps or a proof of concept, and the affected version(s). Expect an initial
acknowledgement within a few days; a fix ships as a new release, with a GitHub Security Advisory
where appropriate.

Security fixes are applied to the **latest released version only** — make sure you're running the
current release before reporting.

## Scope

Per [`SECURITY.md`](https://github.com/rknightion/transceiver-exporter/blob/main/SECURITY.md), the
following are all in scope: the exporter's handling of untrusted input, its privilege usage, and
dependency vulnerabilities.

## What this exporter is trusted with

- **`CAP_NET_ADMIN`** (or root) — required to read transceiver EEPROM over `ethtool` ioctls. See
  [Permissions](permissions.md) for exactly what this grants and how to scope it as tightly as
  possible (capability-only, not full root, wherever your deployment mechanism supports it).
- **Read-only in practice.** The exporter only issues `ethtool` operations to *read* driver info
  and EEPROM; it does not configure interfaces or write module state. `CAP_NET_ADMIN` is a broader
  grant than that narrow need — Linux does not expose a finer-grained capability for these
  particular ioctls — so treat the grant as "this process could, in principle, reconfigure
  networking on this host," not as scoped strictly to reads, when reasoning about blast radius.

## What this exporter exposes

- **`/metrics` is unauthenticated**, by design — the exporter has no built-in auth, TLS, or access
  control of any kind. [`SECURITY.md`](https://github.com/rknightion/transceiver-exporter/blob/main/SECURITY.md)
  states this explicitly: operators are responsible for network-level access control (bind to a
  management interface, firewall port `9458`, or restrict access with a reverse proxy or
  service-mesh policy).
- **The data exposed is hardware/vendor metadata and optical telemetry** — vendor names, serial
  numbers, part numbers, temperatures, voltages, and optical power levels for the host's
  transceivers. This is operationally sensitive (it reveals your physical network topology and
  hardware inventory) but is not credential material.
- **Read-only host filesystem is safe.** Both the shipped container image and Compose file run
  with `read_only: true` — the exporter writes nothing to disk at runtime.

## Supply chain

Release binaries are cross-compiled and cosign-signed (keyless, GitHub OIDC) with SBOMs attached
to every archive — verify `checksums.txt` and its `.sigstore.json` bundle before trusting a
downloaded binary. Container images go through the shared `container-publish` reusable workflow
(SBOM, provenance attestation, Trivy scanning). OpenSSF Scorecard results are published for the
repository — see the badge on the [project README](https://github.com/rknightion/transceiver-exporter).

## Licensing note

This project is AGPL-3.0. If you deploy a modified copy of this exporter as a network service, the
AGPL-3.0 obliges you to make your modified source available to users interacting with that service
over a network — a different obligation than most permissively-licensed exporters carry. See
[the licence](https://github.com/rknightion/transceiver-exporter/blob/main/LICENSE.md) for the
exact terms.
