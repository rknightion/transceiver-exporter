---
title: FAQ
description: Answers about supported NIC drivers, SFP and QSFP modules, permissions, collection behavior, and Prometheus metrics.
---

# FAQ

## What NIC drivers and transceiver types are supported?

Any driver that implements the Linux `ethtool` module-info/module-EEPROM operations, and any
pluggable optic form factor whose EEPROM layout `wobcom/go-ethtool` can parse (SFP, SFP+, SFP28,
QSFP, QSFP+, QSFP28 and similar). There is no exporter-side allowlist of drivers or vendors — see
[Architecture](architecture.md) and [Troubleshooting](troubleshooting.md#an-interfaces-driver-doesnt-support-ethtool-module-eeprom).

## Does it work on non-Linux hosts?

No. `ethtool`'s `SIOCETHTOOL` ioctl interface is Linux-specific. Release binaries are built for
`linux/amd64` and `linux/arm64` only — see [Installation](installation.md).

## Why does it need `CAP_NET_ADMIN` just to read data?

`SIOCETHTOOL` ioctls (which cover both driver info and module EEPROM reads) require
`CAP_NET_ADMIN` on this exporter's code path, even though the exporter itself only performs reads
and never reconfigures an interface. See [Permissions](permissions.md) for how to grant only that
capability rather than running the whole process as root.

## Why does the container run as root?

The shipped image runs as root (UID 0) with all capabilities dropped except `NET_ADMIN`, rather
than as a distroless non-root user with the capability added — an earlier revision tried the
non-root approach and it broke EEPROM access, so it was reverted. See
[Permissions](permissions.md#docker-runs-as-root-not-a-distroless-non-root-user).

## Can I run more than one instance against the same host?

There's no reason to — each instance would enumerate and read the same interfaces independently,
duplicating series with no coordination between them (Prometheus would see two targets reporting
the same `interface` label values). Run exactly one instance per host, on the host network.

## Why is `CAP_NET_ADMIN` sufficient, but the container still uses `--network host`?

They cover different things. `CAP_NET_ADMIN` grants the *capability* to issue networking-related
ioctls; `--network host` (or `hostNetwork: true` on Kubernetes) puts the container in the *network
namespace* that actually contains the host's physical interfaces. Without the host network
namespace, `net.Interfaces()` only sees the container's own virtual interfaces, no matter what
capability the process holds. See [Getting Started](getting-started.md).

## Does it support authentication on `/metrics`?

No — see [Security](security.md). Restrict access at the network layer (firewall, reverse proxy,
or a management-network bind).

## Milliwatts or dBm — which should I use?

Milliwatts (the default) matches what the module reports natively, with no conversion or edge
cases. dBm (`-collector.optical-power-in-dbm`) is more familiar for fiber-optic power budgets but
converts a `0` mW reading to `-Inf` — see
[Troubleshooting](troubleshooting.md#dbm-conversion-produces-inf) before enabling it if your
storage pipeline is picky about non-finite floats.

## How is this related to `wobcom/transceiver-exporter`?

It's an independently maintained continuation of that project and its underlying
`wobcom/go-ethtool` library — not a GitHub fork. See [the project section on the home
page](index.md#project) for the full attribution and licensing note.
