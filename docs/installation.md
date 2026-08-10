---
title: Installation
description: Install transceiver-exporter from a release binary, container image, Docker Compose, or build it from source.
---

# Installation

## Docker (recommended)

Multi-arch images (`linux/amd64`, `linux/arm64`) are published to GHCR on every release, plus an
`:main` edge tag built from the default branch:

```bash
docker run -d \
  --name transceiver-exporter \
  --network host \
  --cap-drop ALL \
  --cap-add NET_ADMIN \
  --read-only \
  ghcr.io/rknightion/transceiver-exporter:latest
```

`--network host` is required — see [Getting Started](getting-started.md). `--cap-add NET_ADMIN`
is required for the `ethtool` ioctl; see [Permissions](permissions.md) for exactly why, and what
happens if you drop it. The image (`gcr.io/distroless/static-debian13` base) runs as **root**, not
a distroless non-root user — the container-level capability drop (`--cap-drop ALL` plus the single
`NET_ADMIN` re-add) is what keeps the container's privileges bounded, not the process UID.

## Docker Compose

A ready-to-use [`compose.yml`](https://github.com/rknightion/transceiver-exporter/blob/main/compose.yml)
ships in the repository, with every command-line option documented inline as commented-out
`command:` entries:

```bash
docker compose up -d
```

It sets `network_mode: host`, `cap_drop: [ALL]`, `cap_add: [NET_ADMIN]`,
`security_opt: [no-new-privileges:true]` and `read_only: true` — the same posture as the `docker
run` example above. Uncomment and edit the `command:` block to change any flag; see
[Configuration](configuration.md) for what each one does.

## Release binaries

Cross-compiled, checksummed, SBOM-attached and cosign-signed **Linux** binaries for `amd64` and
`arm64` are published to every [GitHub release](https://github.com/rknightion/transceiver-exporter/releases/latest)
via GoReleaser. Only `linux` targets are built — this exporter is Linux-only by design (it reads
transceiver EEPROM through Linux-specific `ethtool` ioctls), so there is no `darwin` or `windows`
build, and `386`/`arm`/`armv7` were dropped to match the published container image's architecture
matrix.

```bash
curl -LO https://github.com/rknightion/transceiver-exporter/releases/latest/download/transceiver-exporter_Linux_x86_64.tar.gz
tar xzf transceiver-exporter_Linux_x86_64.tar.gz
sudo setcap cap_net_admin=+ep ./transceiver-exporter   # or run as root / under sudo
./transceiver-exporter -web.listen-address="[::]:9458"
```

Verify the download against `checksums.txt` and its cosign bundle
(`checksums.txt.sigstore.json`), both published alongside the archives on the release page.

## From source

Requires a Go toolchain matching `go.mod` (currently Go 1.26):

```bash
git clone https://github.com/rknightion/transceiver-exporter.git
cd transceiver-exporter
go build -o transceiver-exporter .
./transceiver-exporter -web.listen-address="[::]:9458"
```

A `go build`/`go install` binary reports its version as `dev` — official release and container
builds inject the real version via `-ldflags -X main.version=...`.

Metrics are served at `http://<host>:9458/metrics` by default in every installation method above.
