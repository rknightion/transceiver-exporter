---
title: Architecture
description: How transceiver-exporter talks to ethtool and what happens during one scrape.
---

# Architecture

`transceiver-exporter` is a single Go binary with no background workers, no caching layer, and no
persistent state between scrapes. Everything happens synchronously, once, per HTTP request to the
metrics path.

## Request path

1. **HTTP request arrives** at `-web.telemetry-path` (default `/metrics`).
2. `handleMetricsRequest` (in [`main.go`](https://github.com/rknightion/transceiver-exporter/blob/main/main.go))
   builds a **fresh `prometheus.Registry`** and a fresh `TransceiverCollector` for this request
   only, parameterized by the current flag values (include/exclude lists and regexes, the
   interface-features toggle, and the dBm power-unit toggle). Nothing is reused or cached across
   requests.
3. The registry's `Collect` invokes the collector, which:
   - **Enumerates interfaces** via `net.Interfaces()`, applying the loopback-always-excluded rule
     and the configured include/exclude filters (see
     [Configuration](configuration.md#interface-filtering-rules)).
   - **Opens one `ethtool` handle** (`ethtool.NewEthtool()`, from
     [`wobcom/go-ethtool`](https://github.com/wobcom/go-ethtool)) for the whole scrape.
   - **For each surviving interface**, calls `tool.NewInterface(name, true)`, which pulls both the
     driver info (`ETHTOOL_GDRVINFO`-equivalent) and the module EEPROM
     (`get_module_info`/`get_module_eeprom`) over `SIOCETHTOOL` ioctls.
   - Converts what comes back into Prometheus gauges and writes them to the metrics channel.
4. `promhttp.HandlerFor` serializes the registry's collected metrics into the Prometheus text
   exposition format and writes the HTTP response.

Because everything is per-request, a flag like `-collector.optical-power-in-dbm` or an interface
filter takes effect on the very next scrape — there is nothing to restart or invalidate.

## Failure isolation

Collection failures are handled at two levels, both non-fatal to the HTTP response:

- **Handle-level failure** (opening `ethtool` itself fails — see
  [Permissions](permissions.md)) aborts the whole pass for that scrape; the response is still
  `200 OK`, just empty of transceiver metrics.
- **Interface-level failure** (one interface's `NewInterface` call errors) is logged and that
  interface is skipped; every other interface in the same scrape is unaffected.

Errors are carried on a dedicated Go channel (`errs`) back to the collector wrapper, which logs
each one at ERROR level while metrics collection for the rest of the scrape continues on a
separate goroutine — see the `Collect` method on `transceiverCollectorWrapper` in `main.go`.

## Why per-scrape, not a polling loop

There is no interval, no ticker, and no in-memory metric cache: each scrape does the ioctl work
fresh. This keeps the implementation simple and the data always current as of the scrape instant,
at the cost of the ioctl work happening on the Prometheus scrape's own request latency rather than
being amortized in the background. In practice this is cheap — the ioctl calls are local kernel
calls with no network round trip — but it does mean a very short scrape timeout on a host with many
interfaces could, in principle, be tighter than on an exporter that serves from a pre-populated
cache.

## Package layout

| Package | Role |
|---|---|
| `main` (repo root) | Flag parsing, HTTP server, wiring the collector into a `promhttp` handler. |
| `transceiver-collector` | The `prometheus.Collector` implementation: interface enumeration/filtering, EEPROM-to-metric translation, and the mW→dBm conversion. |
| [`wobcom/go-ethtool`](https://github.com/wobcom/go-ethtool) (external dependency) | The actual `ethtool` ioctl calls and EEPROM parsing this exporter is built on. |
