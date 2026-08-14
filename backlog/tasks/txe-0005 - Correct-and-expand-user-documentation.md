---
id: TXE-0005
title: Correct and expand user documentation
status: Done
assignee: []
created_date: '2026-08-14 16:39'
labels:
  - migrated-from-roadmap
dependencies: []
type: docs
ordinal: 5000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Documentation carried wrong metric names, omitted whole metric families, had no install path, and credited the wrong maintainers after the de-fork.
<!-- SECTION:DESCRIPTION:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Landed in `bbf7ed3`: corrected every metric name, documented the `*_dbm` family, `bus_info` and `laser_supports_monitoring_bool`, added install / Docker / compose sections, and refreshed maintainer and author attribution crediting the original wobcom authors. Later superseded in scope by `25e26df`, which moved the reference into the published `docs/` site built by the m7kni.io hub.
<!-- SECTION:FINAL_SUMMARY:END -->
