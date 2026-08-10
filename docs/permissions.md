---
title: Permissions
description: The CAP_NET_ADMIN requirement for reading transceiver EEPROM, how to grant it minimally, and what happens without it.
---

# Permissions

Reading transceiver diagnostics is not a plain unprivileged network read. `transceiver-exporter`
issues `SIOCETHTOOL` ioctls against each interface to fetch driver info and EEPROM contents, and
**`CAP_NET_ADMIN` is required for those ioctls to succeed**. This page covers what to grant, how to
grant only that, and what the exporter does when it doesn't have it.

## What's required

| Deployment | Minimum grant |
|---|---|
| Docker / Compose | `--cap-drop ALL --cap-add NET_ADMIN` (drop everything, re-add only `NET_ADMIN`) |
| Bare binary | `CAP_NET_ADMIN` via `setcap` on the binary, or run as root |
| systemd unit | `AmbientCapabilities=CAP_NET_ADMIN` + `CapabilityBoundingSet=CAP_NET_ADMIN` |

There is no partial-functionality mode that needs less than this — driver info alone does not
require it, but EEPROM/DOM reads do, and both go through the same ioctl path in practice.

## Docker: runs as root, not a distroless non-root user

The shipped [`Dockerfile`](https://github.com/rknightion/transceiver-exporter/blob/main/Dockerfile)
is explicit about this: it uses `gcr.io/distroless/static-debian13` (the root-default variant, not
`:nonroot`) and sets no `USER` directive, so the container process runs as **root (UID 0)**. The
image comment states this outright: "Runs as root (required for ethtool module EEPROM access)".
This was a deliberate fix — an earlier revision tried running as a non-root user and had to be
reverted (`run container as root for ethtool EEPROM access`, see
[`CHANGELOG.md`](https://github.com/rknightion/transceiver-exporter/blob/main/CHANGELOG.md)).

Running as root inside the container is bounded by capability dropping, not by UID: both the
`docker run` example and the shipped [`compose.yml`](https://github.com/rknightion/transceiver-exporter/blob/main/compose.yml)
pair `cap_drop: [ALL]` with `cap_add: [NET_ADMIN]`, so the process is root but holds exactly one
capability beyond what an unprivileged process gets. Compose additionally sets
`security_opt: [no-new-privileges:true]` and `read_only: true`.

**Do not drop `NET_ADMIN` to "harden" the container further** — without it, `ethtool.NewEthtool()`
or the per-interface `ethtool` calls fail (see below), and you get a container that starts and
serves `/metrics` but reports nothing useful.

`--network host` (required for interface visibility — see [Getting Started](getting-started.md))
means the container also shares the host's network namespace, which is a broader grant than
`NET_ADMIN` alone; there is no way to give this exporter host-NIC visibility without it.

## systemd: the minimum grant without running as root

For a bare binary under systemd, prefer ambient capabilities over running the whole unit as root:

```ini
[Service]
ExecStart=/usr/local/bin/transceiver-exporter -web.listen-address=[::]:9458
User=transceiver-exporter
AmbientCapabilities=CAP_NET_ADMIN
CapabilityBoundingSet=CAP_NET_ADMIN
NoNewPrivileges=true
ProtectSystem=strict
ProtectHome=true
PrivateTmp=true
```

`AmbientCapabilities` grants `CAP_NET_ADMIN` to the process without making it root, and
`CapabilityBoundingSet` caps what it could ever acquire even if compromised, to that one
capability. `NoNewPrivileges=true` is safe here because the capability comes from systemd's
ambient set, not from a setuid/setcap escalation at exec time — it does not need `no_new_privs` to
be false the way a `setcap` binary invoked via `sudo` sometimes does.

The equivalent for a manually-run binary (no systemd) is:

```bash
sudo setcap cap_net_admin=+ep ./transceiver-exporter
./transceiver-exporter -web.listen-address=[::]:9458
```

This is a standard Linux capability pattern, not something this repository ships a unit file
for — adapt paths and the user/group to your environment.

## What happens without it

There is no separate "permission denied" error surface — a missing capability shows up as an
`ethtool` call failing, and the exporter handles that as an ordinary per-interface collection
error:

- If opening the `ethtool` handle itself fails, the whole collection pass for that scrape aborts
  with `could not instanciate ethtool: <error>`, logged at ERROR, and the scrape returns **no
  transceiver metrics at all** for any interface.
- If the handle opens but reading a specific interface fails, that one interface is skipped with
  `error fetching information for interface <name>: <error>`, logged at ERROR, and collection
  continues for the remaining interfaces.

In both cases **the HTTP response is still `200 OK`** — `/metrics` does not fail or return a
non-2xx status because a capability is missing, it just serves fewer (or zero) transceiver series
than expected. A Prometheus target that shows "up" with a plausible scrape duration but no
`transceiver_exporter_identifier_info` series at all is the signature of this failure mode; check
the exporter's own log output, not just scrape health. See
[Troubleshooting](troubleshooting.md#permission-denied-reading-transceiver-data).

## Scope note

The `/metrics` endpoint itself is unauthenticated — see [Security](security.md) for network-level
access control guidance. The capability discussion on this page is about what the *process* needs
from the kernel, which is separate from who can reach the HTTP endpoint.
