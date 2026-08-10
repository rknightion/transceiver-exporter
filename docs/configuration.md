---
title: Configuration
description: Every transceiver-exporter command-line flag, its default, and what it controls.
---

# Configuration Reference

`transceiver-exporter` is configured entirely by **command-line flags** — there is no config file
and no environment-variable layer. Every flag below is read directly from
[`main.go`](https://github.com/rknightion/transceiver-exporter/blob/main/main.go).

| Flag | Type | Default | Description |
|---|---|---|---|
| `-version` | bool | `false` | Print the version, maintainer and original-author credit, then exit. |
| `-web.listen-address` | string | `[::]:9458` | Address the HTTP server listens on. |
| `-web.telemetry-path` | string | `/metrics` | Path under which metrics are exposed. |
| `-collector.interface-features.enable` | bool | `true` | Collect `ethtool` interface features (offloads etc.) as `transceiver_exporter_interface_feature_active`/`_available`. Consider disabling on a many-port switch — this can produce a large number of series per interface. |
| `-collector.optical-power-in-dbm` | bool | `false` | Report laser TX/RX power (and their thresholds) in dBm instead of milliwatts. When enabled, the `*_milliwatts` power metrics are **replaced** by `*_dbm` counterparts, not emitted alongside them — see [Metrics](metrics.md). |
| `-exclude.interfaces` | string | `""` | Comma-separated list of interface names to exclude from collection. |
| `-exclude.interfaces-regex` | string | `""` | Regular expression of interface names to exclude. |
| `-exclude.interfaces-down` | bool | `false` | Don't report on interfaces that are administratively down (`net.FlagUp` not set). |
| `-include.interfaces` | string | `""` | Comma-separated list of interface names to include. When set, only these interfaces are collected. |
| `-include.interfaces-regex` | string | `""` | Regular expression of interface names to include. |

## Interface filtering rules

Filtering is applied once per scrape, in `getMonitoredInterfaces`:

1. **Loopback is always excluded**, unconditionally, regardless of any flag.
2. Setting **both** `-exclude.interfaces` and `-include.interfaces` (the plain, non-regex list
   forms) at the same time is a **startup error**: "Cannot include and exclude interfaces at the
   same time". The regex forms are independent of this check and may be combined with each other
   or with the list forms.
3. `-exclude.interfaces-down` drops interfaces without `IFF_UP` before any include/exclude list or
   regex is evaluated.
4. Comma-separated values in `-exclude.interfaces` / `-include.interfaces` are trimmed of
   surrounding whitespace.
5. An `-exclude.interfaces-regex` / `-include.interfaces-regex` value of the empty string (the
   default) is treated as "no regex filter" rather than a regex that matches everything — the
   default flag value can never accidentally exclude or restrict every interface.

An invalid regex passed to `-exclude.interfaces-regex` or `-include.interfaces-regex` fails at
startup (`error compiling exclude.interfaces-regex: ...` / `error compiling include.interfaces-regex: ...`)
before the HTTP server starts.

## Power unit

`-collector.optical-power-in-dbm` is evaluated once per scrape at request time (it is passed into
the collector constructor built on every request), so it takes effect immediately — it is not a
startup-only flag baked into the process. The conversion is `10 * log10(milliwatts)`; a reading of
exactly `0` mW converts to `-Inf`, which some consumers (notably strict JSON encoders) do not
represent — see [Troubleshooting](troubleshooting.md).

## Example

```bash
./transceiver-exporter \
  -web.listen-address="[::]:9458" \
  -collector.interface-features.enable=false \
  -exclude.interfaces-regex="^(lo|docker|veth|br-).*" \
  -collector.optical-power-in-dbm=true
```
