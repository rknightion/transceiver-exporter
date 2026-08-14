---
id: TXE-0001
title: Metric prefix mismatch between emitted metrics and docs
status: Done
assignee: []
created_date: '2026-08-14 16:39'
labels:
  - migrated-from-roadmap
  - breaking
dependencies: []
type: bug
ordinal: 1000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Emitted metrics used the `transceiver_` prefix while the documentation promised `transceiver_exporter_`, so every documented metric name was wrong against a live target.
<!-- SECTION:DESCRIPTION:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Fixed in `36800a4` by switching `const prefix` in `transceiver-collector/collector.go` to `transceiver_exporter_`. Breaking change for existing dashboards and alerts. The prefix constant is now the only source of the prefix, but `docs/metrics.md` still restates full metric names by hand — see the Wave operating model doc for the two-file rule this defect produced.
<!-- SECTION:FINAL_SUMMARY:END -->
