---
title: Troubleshooting
description: Common transceiver-exporter failure modes — unsupported drivers, permission errors, missing transceivers, and missing DOM data.
---

# Troubleshooting

## Permission denied reading transceiver data

**Symptom.** `/metrics` returns `200 OK`, but there are no `transceiver_exporter_identifier_info`
(or any EEPROM-derived) series, and the process log shows one of:

```text
could not instanciate ethtool: <error>
error fetching information for interface <name>: <error>
```

**Cause.** The process lacks `CAP_NET_ADMIN`. This is not surfaced as an HTTP error — the scrape
still succeeds, it simply carries fewer or zero transceiver series. See
[Permissions](permissions.md) for exactly what's required and why a Prometheus target showing
"up" can still mean this exporter has nothing useful to report.

**Fix.** Grant `CAP_NET_ADMIN` per [Permissions](permissions.md) — `--cap-add NET_ADMIN` in
Docker/Compose, `AmbientCapabilities=CAP_NET_ADMIN` under systemd, or `setcap` on the binary.

## An interface's driver doesn't support ethtool module EEPROM

**Symptom.** An interface shows up with `driver_name_info`, `driver_version_info` etc., but never
gets `identifier_info`, `vendor_name_info`, or any other EEPROM/DOM metric — even with the correct
capability and a transceiver physically present.

**Cause.** `transceiver-exporter` has no hardcoded list of supported NIC drivers or vendors; it
calls whatever the kernel's `ethtool` `get_module_info`/`get_module_eeprom` operations return for
that interface. A driver that does not implement those operations (common for some virtual NICs,
some older or minimal drivers, and interfaces that are not pluggable-optics ports at all — a
bonded/virtual interface, a management NIC with a fixed copper PHY) simply has no EEPROM data to
read, and the exporter emits driver-info metrics only. This is expected behaviour, not a bug in
the exporter.

**Fix.** Confirm the interface actually carries a pluggable transceiver and that
`ethtool -m <interface>` on the host succeeds outside the exporter — if the host's own `ethtool`
CLI cannot read the module EEPROM either, the exporter cannot either; this is a driver/kernel
capability, not something the exporter's flags can work around. If the interface is a virtual or
non-optic interface, exclude it with `-exclude.interfaces` or `-exclude.interfaces-regex` (see
[Configuration](configuration.md)) to keep it out of your dashboards' interface lists.

## An interface with no transceiver plugged in

**Symptom.** Same as above — driver-info metrics only, no identity or DOM metrics — but on an
interface you know supports pluggable optics.

**Cause.** There is genuinely no module in the cage. The `ethtool` EEPROM read has nothing to
return.

**Fix.** Nothing to fix in the exporter; this is a correct, empty result. If you want to alert on
"port that should have optics has none", build that from the *absence* of
`transceiver_exporter_identifier_info` for an interface name you expect to always carry one,
rather than from a metric value — there is no explicit "transceiver present" boolean, presence is
implied by the identity metrics existing at all.

## A transceiver that reports no DOM data

**Symptom.** `identifier_info`, `vendor_name_info` and the other identity metrics appear for the
interface, but `module_temperature_degrees_celsius`, `module_voltage_volts`, and all
`laser_*` metrics are absent.

**Cause.** `module_supports_monitoring_bool` is `0` for this module — some transceivers,
especially older or budget optics, expose vendor/identity EEPROM fields but implement no digital
optical monitoring (DOM) page at all. This is a property of the specific transceiver hardware, not
a configuration issue.

**Fix.** Check `transceiver_exporter_module_supports_monitoring_bool` for that interface. A value
of `0` means the module itself has no DOM data to give — there is nothing further to configure.
Also check per-laser: even on a DOM-capable module, an individual laser can report
`laser_supports_monitoring_bool=0` (seen on some multi-lane QSFP optics with partial monitoring
support per lane) while the module-level temperature/voltage metrics are still present.

## Threshold metrics missing for a value that IS present

**Symptom.** `laser_bias_current_milliamperes` (or another measurement) has a value, but none of
its four threshold series (`high_alarm`, `high_warning`, `low_alarm`, `low_warning`) exist.

**Cause.** `*_supports_thresholds_bool` for that measurement is `0` — the module reports live
readings but not alarm/warning threshold pages. Check the `_supports_thresholds_bool` metric for
that same measurement before assuming a missing threshold series is a bug.

## dBm conversion produces `-Inf`

**Symptom.** A laser power reading of exactly `0` mW, converted with
`-collector.optical-power-in-dbm`, becomes `-Inf` rather than a finite number.

**Cause.** `10 * log10(0)` is mathematically `-Inf`; the exporter does not special-case a zero
reading. This is documented, expected math, not a bug — see [Metrics](metrics.md).

**Fix.** If your scrape or storage pipeline can't handle `-Inf` (some JSON encoders that don't
fully implement IEEE 754 choke on it — Prometheus's own text exposition format handles `-Inf`
natively), leave `-collector.optical-power-in-dbm` off and consume the milliwatt series instead,
converting downstream where you control the numeric handling.

## Interface filtering flags conflict at startup

**Symptom.** The process refuses to start, or the compiled binary logs a fatal error naming a
regex flag.

**Cause.** Either `-exclude.interfaces` and `-include.interfaces` were both set (mutually
exclusive — see [Configuration](configuration.md#interface-filtering-rules)), or one of
`-exclude.interfaces-regex`/`-include.interfaces-regex` is not a valid regular expression.

**Fix.** Use only one of the list-based include/exclude flags at a time, and validate any regex
flag value independently (e.g. with `python3 -c "import re; re.compile('...')"`) before passing it.

## Still stuck?

- **[Search existing issues](https://github.com/rknightion/transceiver-exporter/issues)**.
- **[Open a new issue](https://github.com/rknightion/transceiver-exporter/issues/new)** — include
  the exporter version (`transceiver-exporter -version`), the flags in use, the NIC driver name
  (`ethtool -i <interface>` on the host), and whether `ethtool -m <interface>` succeeds outside the
  exporter.
