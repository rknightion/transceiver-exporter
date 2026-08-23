---
id: TXE-0008
title: Upgrade Go toolchain to 1.27
status: In Progress
assignee: []
created_date: '2026-08-23 19:06'
updated_date: '2026-08-23 19:49'
labels: []
dependencies: []
ordinal: 8000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Adopt Go 1.27 consistently across the application, nested modules, build images, CI configuration, setup automation, and version-specific contributor documentation.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 All active Go module and toolchain pins require Go 1.27
- [x] #2 Build images, CI jobs, setup automation, and current documentation agree with the Go 1.27 requirement
- [x] #3 The repository green-bar validation passes under Go 1.27
<!-- AC:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 go build ./...
- [x] #2 go vet ./...
- [x] #3 go test -race ./...
- [x] #4 golangci-lint run
<!-- DOD:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Inventory every active Go version pin, including nested modules and container or CI toolchains. 2. Update the pins and version-specific documentation to Go 1.27.0 without changing historical records or fixtures. 3. Run the repository-defined validation gate, review the diff, commit to main, push, and confirm hosted CI.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Local Go 1.27.0 evidence: go build ./..., go vet ./..., go test -race ./..., and golangci-lint run all passed with 0 lint issues. Final diff check passed. CodeRabbit was skipped because the change is declarative toolchain, CI, container, and documentation configuration only.
<!-- SECTION:NOTES:END -->
