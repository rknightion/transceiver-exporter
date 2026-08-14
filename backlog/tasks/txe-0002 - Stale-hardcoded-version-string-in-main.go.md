---
id: TXE-0002
title: Stale hardcoded version string in main.go
status: Done
assignee: []
created_date: '2026-08-14 16:39'
labels:
  - migrated-from-roadmap
dependencies: []
type: bug
ordinal: 2000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
The `main.go` version fallback was the literal `1.5.1`, which outlived the release it named and misreported the version of any build not made by the official pipeline.
<!-- SECTION:DESCRIPTION:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Fixed in `9790544`: fallback changed to `dev`; official builds inject the real version via ldflags. `main_test.go` was added in the same commit. Release numbers must not reappear in Go source — release-please owns versioning.
<!-- SECTION:FINAL_SUMMARY:END -->
