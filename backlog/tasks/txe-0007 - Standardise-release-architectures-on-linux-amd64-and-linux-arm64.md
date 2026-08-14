---
id: TXE-0007
title: Standardise release architectures on linux/amd64 and linux/arm64
status: Done
assignee: []
created_date: '2026-08-14 16:39'
labels:
  - migrated-from-roadmap
  - breaking
dependencies: []
type: chore
ordinal: 7000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Release binaries and container images were built for different architecture sets, so a documented architecture could exist as a binary but not an image, or the reverse.
<!-- SECTION:DESCRIPTION:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Settled in `2b50c83` (`build!:`): `386` and `arm`/`armv7` dropped from `.goreleaser.yaml` so binaries match the container-publish matrix at `linux/amd64` + `linux/arm64`. Breaking for anyone consuming the dropped binaries or images. The `!` in that subject is what makes the pending release-please PR a major bump.
<!-- SECTION:FINAL_SUMMARY:END -->
