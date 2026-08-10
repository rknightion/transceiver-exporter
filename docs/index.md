---
title: transceiver-exporter — Prometheus exporter for pluggable optical transceivers
description: Prometheus exporter that reads SFP, SFP+, QSFP and other pluggable transceiver diagnostics from Linux NICs over ethtool.
image: assets/social-card.png
---

# transceiver-exporter

**Read the diagnostics off your pluggable optics and expose them as Prometheus metrics.**
`transceiver-exporter` is a small Go binary that walks every network interface on a Linux host,
asks the kernel driver for its transceiver EEPROM over `ethtool`, and serves the result — vendor
identity, wavelength, power class, and live digital optical monitoring (DOM) data where the module
supports it — on a `/metrics` endpoint.

It works with any pluggable form factor the kernel's `ethtool` ioctls can see EEPROM for: SFP,
SFP+, SFP28, QSFP, QSFP+, QSFP28 and similar modules. There is no hardcoded form-factor or vendor
list — whatever the NIC driver's `get_module_info`/`get_module_eeprom` implementation reports is
what gets exported. This also means driver support is the exporter's actual dependency: a NIC
driver that does not implement the ethtool module EEPROM ops (or a virtual/software interface with
no transceiver at all) simply yields no EEPROM metrics for that interface — see
[Troubleshooting](troubleshooting.md).

The source, releases and issue tracker live on
**[GitHub](https://github.com/rknightion/transceiver-exporter)**.

## Quickstart

Linux only — the exporter reads transceiver EEPROM through the kernel's `ethtool`
ioctls:

```bash
docker run -d \
  --name transceiver-exporter \
  --network host \
  --cap-drop ALL \
  --cap-add NET_ADMIN \
  --read-only \
  ghcr.io/rknightion/transceiver-exporter:latest
```

Metrics are then at `http://localhost:9458/metrics`. `--network host` and
`--cap-add NET_ADMIN` are both load-bearing: in bridge mode the exporter only
sees container-local interfaces, and the `SIOCETHTOOL` ioctl needs the
capability.

## Start here

<div class="grid cards" markdown>

- **[Getting started](getting-started.md)** — run the binary, scrape it, see a metric.
- **[Installation](installation.md)** — Docker, Docker Compose, a release binary, or building
  from source.
- **[Configuration](configuration.md)** — every command-line flag and its default.
- **[Metrics catalog](metrics.md)** — every metric this exporter emits, with labels and units.
- **[Permissions](permissions.md)** — the `CAP_NET_ADMIN` requirement, and what happens without it.

</div>

## What it collects

| Area | What you get |
|---|---|
| **Driver info** | Driver name, driver version, firmware version, bus info, expansion ROM version — from the interface's `ethtool` driver info, independent of any transceiver being present. |
| **Interface features** | The `ethtool` feature set (offloads, etc.) per interface, as available/active gauges. Optional — can produce a lot of series on a many-port switch. |
| **Transceiver identity** | Identifier type, encoding, power class and its wattage, signaling rate, supported link lengths per media, vendor name/part number/revision/serial/OUI, manufacture date code, wavelength. |
| **Digital optical monitoring (DOM)** | Module temperature and supply voltage, per-laser bias current, and per-laser TX/RX optical power — each with its supported alarm/warning thresholds, when the module reports them. |

Everything is a per-interface (and, for laser measurements, per-interface-per-laser-index) gauge.
See the [metrics catalog](metrics.md) for the exhaustive list.

## How it works

One HTTP scrape of `/metrics` triggers one collection pass: enumerate interfaces, open an
`ethtool` handle, and query each interface's driver info and EEPROM in turn. There is no
background polling loop and no caching between scrapes — every scrape is a fresh read. See
[Architecture](architecture.md).

## Project

transceiver-exporter is an **independently maintained continuation** of
[wobcom/transceiver-exporter](https://github.com/wobcom/transceiver-exporter) and the underlying
[wobcom/go-ethtool](https://github.com/wobcom/go-ethtool) library. All original credit for the
project's design belongs to the wobcom authors — [@fluepke](https://github.com/fluepke),
[@BarbarossaTM](https://github.com/BarbarossaTM) and [@vidister](https://github.com/vidister).
This repository is not in GitHub's fork network; it is a from-scratch continuation carrying the
same license.

**Licensed under AGPL-3.0**, unlike most exporters in this fleet, which use a permissive license.
If you run a modified version of this exporter as a network service, the AGPL-3.0 requires you to
make your modified source available to users of that service — see
[`LICENSE.md`](https://github.com/rknightion/transceiver-exporter/blob/main/LICENSE.md) in the
repository for the full terms.

Bug reports, feature requests and pull requests are welcome — see the
[open issues](https://github.com/rknightion/transceiver-exporter/issues) or the
[latest release](https://github.com/rknightion/transceiver-exporter/releases/latest).
