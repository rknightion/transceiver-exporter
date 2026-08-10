---
title: Metrics
description: Every metric transceiver-exporter emits, with its labels, type, unit and meaning.
---

# Metrics Reference

All metrics are **Gauges**, prefixed `transceiver_exporter_`, and derived from
[`transceiver-collector/collector.go`](https://github.com/rknightion/transceiver-exporter/blob/main/transceiver-collector/collector.go).
Every metric carries an `interface` label naming the NIC it was read from; laser measurements
additionally carry `laser_index` (a transceiver can have more than one laser, e.g. parallel-fiber
QSFP optics).

There is no background collection loop: everything below is read fresh, synchronously, on every
scrape of `/metrics`.

## Emission conditions

Not every metric appears for every interface, by design:

- **Driver-info metrics** (`driver_name_info`, `driver_version_info`, `firmware_version_info`,
  `bus_info`, `expansion_rom_version_info`) appear for any interface the driver reports info for,
  whether or not it has a transceiver plugged in.
- **Interface-feature metrics** (`interface_feature_active`/`_available`) appear only when
  `-collector.interface-features.enable` is `true` (the default) and the driver's feature query
  succeeds.
- **Transceiver-identity and EEPROM-derived metrics** (everything from `identifier_info` onward)
  appear only when the driver returns EEPROM data at all — no transceiver plugged in, or a driver
  that doesn't implement the `ethtool` module-EEPROM operations, means none of these are emitted
  for that interface.
- **DOM metrics** (module temperature/voltage, laser bias/power) appear only when the module
  reports `SupportsMonitoring() == true`. A transceiver without digital diagnostics still yields
  the identity metrics above, just none of these.
- **Per-laser metrics** appear only for lasers where `laser.SupportsMonitoring()` is true.
- **Threshold metrics** (`*_high_alarm_threshold_*`, `*_high_warning_threshold_*`,
  `*_low_alarm_threshold_*`, `*_low_warning_threshold_*`) for a given measurement appear only when
  that measurement's `SupportsThresholds()` is true; the matching `*_supports_thresholds_bool`
  metric is always emitted alongside the value once monitoring is supported, so it tells you
  whether to expect the four threshold series.

A failure reading one interface (permission denied, a driver error) is logged and skipped — it
does not stop metrics for other interfaces. See [Troubleshooting](troubleshooting.md).

## Driver info

| Metric | Labels | Value | Meaning |
|---|---|---|---|
| `transceiver_exporter_driver_name_info` | `interface`, `driver_name` | always `1` | Info metric; the driver name is carried in the label. |
| `transceiver_exporter_driver_version_info` | `interface`, `driver_version` | always `1` | Driver version, as a label. |
| `transceiver_exporter_firmware_version_info` | `interface`, `firmware_version` | always `1` | NIC firmware version, as a label. |
| `transceiver_exporter_bus_info` | `interface`, `bus_information` | always `1` | Bus address (e.g. PCI slot) reported by the driver, as a label. |
| `transceiver_exporter_expansion_rom_version_info` | `interface`, `expansion_rom_version` | always `1` | Expansion ROM version, as a label. |

## Interface features

| Metric | Labels | Unit | Meaning |
|---|---|---|---|
| `transceiver_exporter_interface_feature_available` | `interface`, `feature_name` | boolean (`1`/`0`) | `1` if the named `ethtool` feature (offload, etc.) is available on this interface. |
| `transceiver_exporter_interface_feature_active` | `interface`, `feature_name` | boolean (`1`/`0`) | `1` if the named feature is currently active/enabled. |

## Transceiver identity

| Metric | Labels | Unit | Meaning |
|---|---|---|---|
| `transceiver_exporter_identifier_info` | `interface`, `identifier` | info | Type of transceiver (e.g. SFP, QSFP), as a label. |
| `transceiver_exporter_encoding_info` | `interface`, `encoding` | info | Transceiver line-encoding scheme, as a label. |
| `transceiver_exporter_powerclass_info` | `interface` | integer | Highest power class supported by the transceiver (the raw class byte). |
| `transceiver_exporter_powerclass_watts` | `interface` | watts | Maximum wattage allowed by that power class. |
| `transceiver_exporter_signalingrate_bauds_per_second` | `interface` | baud/s | Signaling rate supported by the transceiver. |
| `transceiver_exporter_supported_link_length_meter` | `interface`, `media` | meters | Maximum supported link length, one series per supported media type. |
| `transceiver_exporter_vendor_name_info` | `interface`, `vendor_name` | info | Vendor name, as a label. |
| `transceiver_exporter_vendor_part_number_info` | `interface`, `vendor_part_number` | info | Vendor part number, as a label. |
| `transceiver_exporter_vendor_revision_info` | `interface`, `vendor_revision` | info | Vendor hardware revision, as a label. |
| `transceiver_exporter_vendor_serial_number_info` | `interface`, `vendor_serial_number` | info | Vendor serial number, as a label. |
| `transceiver_exporter_vendor_oui_info` | `interface`, `vendor_oui` | info | Vendor IEEE company (OUI) ID, as a label. |
| `transceiver_exporter_date_code_unix_time` | `interface` | Unix seconds | Vendor-supplied manufacture date code. |
| `transceiver_exporter_wavelength_nanometer` | `interface` | nanometers | Optical wavelength. |
| `transceiver_exporter_module_supports_monitoring_bool` | `interface` | boolean (`1`/`0`) | `1` if the module supports digital optical monitoring (DOM). Gates all metrics in the next two sections. |

## Module temperature and voltage

Emitted only when `module_supports_monitoring_bool` is `1`.

| Metric | Labels | Unit | Meaning |
|---|---|---|---|
| `transceiver_exporter_module_temperature_degrees_celsius` | `interface` | °C | Module temperature. |
| `transceiver_exporter_module_temperature_supports_thresholds_bool` | `interface` | boolean | `1` if temperature alarm/warning thresholds are supported. |
| `transceiver_exporter_module_temperature_high_alarm_threshold_degrees_celsius` | `interface` | °C | High-alarm threshold. |
| `transceiver_exporter_module_temperature_high_warning_threshold_degrees_celsius` | `interface` | °C | High-warning threshold. |
| `transceiver_exporter_module_temperature_low_alarm_threshold_degrees_celsius` | `interface` | °C | Low-alarm threshold. |
| `transceiver_exporter_module_temperature_low_warning_threshold_degrees_celsius` | `interface` | °C | Low-warning threshold. |
| `transceiver_exporter_module_voltage_volts` | `interface` | V | Module supply voltage. |
| `transceiver_exporter_module_voltage_supports_thresholds_bool` | `interface` | boolean | `1` if voltage alarm/warning thresholds are supported. |
| `transceiver_exporter_module_voltage_high_alarm_threshold_voltage` | `interface` | V | High-alarm threshold. |
| `transceiver_exporter_module_voltage_high_warning_threshold_voltage` | `interface` | V | High-warning threshold. |
| `transceiver_exporter_module_voltage_low_alarm_threshold_voltage` | `interface` | V | Low-alarm threshold. |
| `transceiver_exporter_module_voltage_low_warning_threshold_voltage` | `interface` | V | Low-warning threshold. |

## Per-laser metrics

Emitted per `laser_index` (`interface` + `laser_index` labels), only for lasers where
`laser.SupportsMonitoring()` is true.

| Metric | Unit | Meaning |
|---|---|---|
| `transceiver_exporter_laser_supports_monitoring_bool` | boolean | `1` if this laser supports real-time monitoring. Gates the rest of this table. |
| `transceiver_exporter_laser_bias_current_milliamperes` | mA | Laser bias current. |
| `transceiver_exporter_laser_bias_current_supports_thresholds_bool` | boolean | `1` if bias-current thresholds are supported. |
| `transceiver_exporter_laser_bias_current_high_alarm_threshold_milliamperes` | mA | High-alarm threshold. |
| `transceiver_exporter_laser_bias_current_high_warning_threshold_milliamperes` | mA | High-warning threshold. |
| `transceiver_exporter_laser_bias_current_low_alarm_threshold_milliamperes` | mA | Low-alarm threshold. |
| `transceiver_exporter_laser_bias_current_low_warning_threshold_milliamperes` | mA | Low-warning threshold. |

### TX/RX optical power

`-collector.optical-power-in-dbm` (default `false`) selects which of the two families below is
emitted. **The two families are mutually exclusive** — enabling the flag replaces the `_milliwatts`
series with `_dbm` series, it does not add the latter alongside the former. The
`*_supports_thresholds_bool` metrics are emitted unconditionally either way.

| Metric (mW, default) | Metric (dBm, `-collector.optical-power-in-dbm`) | Unit | Meaning |
|---|---|---|---|
| `transceiver_exporter_laser_tx_power_supports_thresholds_bool` | *(same, unit-independent)* | boolean | `1` if TX power thresholds are supported. |
| `transceiver_exporter_laser_tx_power_milliwatts` | `transceiver_exporter_laser_tx_power_dbm` | mW / dBm | Laser transmit power. |
| `transceiver_exporter_laser_tx_power_high_alarm_threshold_milliwatts` | `transceiver_exporter_laser_tx_power_high_alarm_threshold_dbm` | mW / dBm | High-alarm threshold. |
| `transceiver_exporter_laser_tx_power_high_warning_threshold_milliwatts` | `transceiver_exporter_laser_tx_power_high_warning_threshold_dbm` | mW / dBm | High-warning threshold. |
| `transceiver_exporter_laser_tx_power_low_alarm_threshold_milliwatts` | `transceiver_exporter_laser_tx_power_low_alarm_threshold_dbm` | mW / dBm | Low-alarm threshold. |
| `transceiver_exporter_laser_tx_power_low_warning_threshold_milliwatts` | `transceiver_exporter_laser_tx_power_low_warning_threshold_dbm` | mW / dBm | Low-warning threshold. |
| `transceiver_exporter_laser_rx_power_supports_thresholds_bool` | *(same, unit-independent)* | boolean | `1` if RX power thresholds are supported. |
| `transceiver_exporter_laser_rx_power_milliwatts` | `transceiver_exporter_laser_rx_power_dbm` | mW / dBm | Laser receive power. |
| `transceiver_exporter_laser_rx_power_high_alarm_threshold_milliwatts` | `transceiver_exporter_laser_rx_power_high_alarm_threshold_dbm` | mW / dBm | High-alarm threshold. |
| `transceiver_exporter_laser_rx_power_high_warning_threshold_milliwatts` | `transceiver_exporter_laser_rx_power_high_warning_threshold_dbm` | mW / dBm | High-warning threshold. |
| `transceiver_exporter_laser_rx_power_low_alarm_threshold_milliwatts` | `transceiver_exporter_laser_rx_power_low_alarm_threshold_dbm` | mW / dBm | Low-alarm threshold. |
| `transceiver_exporter_laser_rx_power_low_warning_threshold_milliwatts` | `transceiver_exporter_laser_rx_power_low_warning_threshold_dbm` | mW / dBm | Low-warning threshold. |

The dBm conversion is `10 * log10(milliwatts)`. A raw reading of exactly `0` mW converts to
`-Inf`, which some scrapers or exporters not fully implementing IEEE 754 may mishandle — see
[Troubleshooting](troubleshooting.md).
