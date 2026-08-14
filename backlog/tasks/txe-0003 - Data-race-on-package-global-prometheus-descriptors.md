---
id: TXE-0003
title: Data race on package-global prometheus descriptors
status: Done
assignee: []
created_date: '2026-08-14 16:39'
labels:
  - migrated-from-roadmap
dependencies: []
type: bug
ordinal: 3000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
`NewCollector` mutated package-global `*prometheus.Desc` values per scrape, which is a genuine data race under concurrent `/metrics` requests.
<!-- SECTION:DESCRIPTION:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Fixed in `36800a4`: all ~60 descriptors are now built once in `init()` and never assigned to during a scrape. A concurrency test in `transceiver-collector/collector_test.go` runs under `-race` to keep it fixed. Any new metric must follow the same shape.
<!-- SECTION:FINAL_SUMMARY:END -->
