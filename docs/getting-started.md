---
title: Getting Started
description: Run transceiver-exporter, scrape it, and see your first transceiver metric.
---

# Getting Started

This walks you through running `transceiver-exporter` on a Linux host with pluggable optics and
confirming metrics land on `/metrics`.

## Prerequisites

- A **Linux host** (this is a Linux-only tool — it reads transceiver EEPROM via `ethtool` ioctls,
  which do not exist on other platforms).
- At least one NIC with a **pluggable transceiver** whose driver implements the `ethtool`
  `get_module_info`/`get_module_eeprom` operations. Most mainstream NIC drivers do; see
  [Troubleshooting](troubleshooting.md) if a given interface yields no transceiver metrics.
- **`CAP_NET_ADMIN`**, or root — required for the `SIOCETHTOOL` ioctl the exporter uses to read
  driver info and EEPROM. See [Permissions](permissions.md) before deploying anywhere that isn't a
  quick local test.

## Run it

Docker is the fastest path — see [Installation](installation.md) for the full set of options
(binary, Docker Compose, building from source). The short version:

```bash
docker run -d \
  --name transceiver-exporter \
  --network host \
  --cap-drop ALL \
  --cap-add NET_ADMIN \
  --read-only \
  ghcr.io/rknightion/transceiver-exporter:latest
```

`--network host` is required: the exporter enumerates interfaces with `net.Interfaces()`, which in
a bridged container network namespace would only see the container's own virtual interfaces
(`eth0`, `lo`), never the host's physical NICs.

## Confirm it's serving metrics

```bash
curl http://localhost:9458/metrics
```

You should see driver-info metrics for every non-loopback interface immediately, for example:

```text
transceiver_exporter_driver_name_info{driver_name="ixgbe",interface="eth0"} 1
```

If a plugged-in, DOM-capable transceiver was found on that interface, you'll also see readings
such as:

```text
transceiver_exporter_module_temperature_degrees_celsius{interface="eth0"} 34.5
transceiver_exporter_laser_rx_power_milliwatts{interface="eth0",laser_index="0"} 0.512
```

An interface with no transceiver plugged in, or whose driver doesn't expose EEPROM, will still
show up with driver-info metrics but no `identifier_info`/module/laser metrics at all — see
[Troubleshooting](troubleshooting.md#a-transceiver-that-reports-no-dom-data).

## Point Prometheus at it

Add a scrape job:

```yaml
scrape_configs:
  - job_name: transceiver-exporter
    static_configs:
      - targets: ["<host>:9458"]
```

Since the exporter runs with `--network host`, scrape it on the host's own address, not a
container-network address.

## Next steps

- [Installation](installation.md) — Docker Compose, a release binary, or building from source.
- [Configuration](configuration.md) — every command-line flag, including interface filtering and
  the dBm power-unit switch.
- [Metrics catalog](metrics.md) — the full metric reference.
- [Permissions](permissions.md) — the `CAP_NET_ADMIN` requirement in detail, and running under
  systemd with the minimum grant.
