# Changelog

## [2.0.0](https://github.com/rknightion/transceiver-exporter/compare/v1.6.0...v2.0.0) (2026-09-05)


### ⚠ BREAKING CHANGES

* This release renames every exported metric from the `transceiver_` prefix to `transceiver_exporter_`; dashboards and alerts must be updated. Additionally, `linux/386` and `linux/arm`/`armv7` release binaries and container images are no longer produced — only `linux/amd64` and `linux/arm64` are supported.

### Features

* **ci:** add govulncheck to align with the rest of the 2otel family ([f67fe75](https://github.com/rknightion/transceiver-exporter/commit/f67fe75b214a7a5265845e0e437866e818d5e67f))
* mint release-please token from the OpenBao broker ([f8138f7](https://github.com/rknightion/transceiver-exporter/commit/f8138f743d88589e4513be150a64f7b66da9be69))


### Bug Fixes

* author is Rob Knight, not Rob Knighton ([67d5b69](https://github.com/rknightion/transceiver-exporter/commit/67d5b69d63b2854cfcf3a307b3d792dcff9a7e8d))
* **ci:** repin rknightion/.github refs to v1.9.7 so Renovate can track them ([2485ebd](https://github.com/rknightion/transceiver-exporter/commit/2485ebdd1049d9537232f170065a00205802f2a9))
* **collector:** use transceiver_exporter_ metric prefix and build descriptors once ([36800a4](https://github.com/rknightion/transceiver-exporter/commit/36800a45fc0fff7655f3cd597c205775cf4931bc))
* default version to "dev" instead of stale release string ([9790544](https://github.com/rknightion/transceiver-exporter/commit/9790544067d73a094194acfc0e54b8f264b4f17d))
* **deps:** update module github.com/prometheus/client_golang to v1.24.1 ([#8](https://github.com/rknightion/transceiver-exporter/issues/8)) ([758f6b5](https://github.com/rknightion/transceiver-exporter/commit/758f6b5514bf542babd04e60acfb9473fb47b0f5))
* **deps:** update module github.com/sirupsen/logrus to v1.10.0 ([#11](https://github.com/rknightion/transceiver-exporter/issues/11)) ([77227fb](https://github.com/rknightion/transceiver-exporter/commit/77227fb6df2ee40b2b895fefa5eda58edeb1c006))
* **deps:** update module github.com/sirupsen/logrus to v1.10.1 ([#12](https://github.com/rknightion/transceiver-exporter/issues/12)) ([ec5801d](https://github.com/rknightion/transceiver-exporter/commit/ec5801dfc82abc63da817cef8fed4c4a15939a43))
* **deps:** update module github.com/sirupsen/logrus to v1.10.2 ([#15](https://github.com/rknightion/transceiver-exporter/issues/15)) ([b0e53e1](https://github.com/rknightion/transceiver-exporter/commit/b0e53e1f9630232c7b0767a77f889a7f0f996e73))
* **deps:** update module github.com/wobcom/go-ethtool to v1.0.2 ([#27](https://github.com/rknightion/transceiver-exporter/issues/27)) ([bd4ac28](https://github.com/rknightion/transceiver-exporter/commit/bd4ac28cf3b6890ea8a293b30f1c5a45153c331e))


### Documentation

* add community-health files and v1 roadmap ([3854a50](https://github.com/rknightion/transceiver-exporter/commit/3854a50f52ab4dbe7f48450023a781f85a1999ca))
* **backlog:** complete Go 1.27 upgrade ([99de885](https://github.com/rknightion/transceiver-exporter/commit/99de88590aec5a92d3fc90320857dac06d97a333))
* **backlog:** sync fan-out protocol — CodeRabbit review gate ([7963110](https://github.com/rknightion/transceiver-exporter/commit/796311083563ebc15d9a81bfcdaae795fb67f014))
* **backlog:** sync fan-out protocol — success criteria vs write authority ([b417c72](https://github.com/rknightion/transceiver-exporter/commit/b417c72ff52a38790e25b2020dafc73c26879993))
* document all metrics, add install/Docker guide, refresh attribution ([bbf7ed3](https://github.com/rknightion/transceiver-exporter/commit/bbf7ed3f70068a9512078f5fe84136ad5a69ab44))
* improve search result descriptions ([3d04ed3](https://github.com/rknightion/transceiver-exporter/commit/3d04ed31a738626cc4b261b8e394d831ccc7c695))
* publish a documentation site and join the m7kni.io fleet ([25e26df](https://github.com/rknightion/transceiver-exporter/commit/25e26df70db3fe685389f2b71adc07d6a7bd83fd))
* put a copy-paste quickstart on the landing page ([c527f66](https://github.com/rknightion/transceiver-exporter/commit/c527f660bee9a8e884a8859704d001ccdc077cbf))
* re-import fan-out protocol (context-cost rules) ([0962df7](https://github.com/rknightion/transceiver-exporter/commit/0962df771a19d958c176b2b8233ecd333b450b2b))
* re-import the fan-out protocol at c1e6cb0 ([957fc50](https://github.com/rknightion/transceiver-exporter/commit/957fc50c526030ec720913c54bfb08e01d957203))
* re-render the fan-out protocol from agent-docs ([6a96261](https://github.com/rknightion/transceiver-exporter/commit/6a96261292cbc563742660dbf374606fba5d5466))
* re-render the fan-out protocol from agent-docs 711db6c ([6a29d16](https://github.com/rknightion/transceiver-exporter/commit/6a29d16bce61d706d13cd666a2eb824f8fb554f1))
* re-render the fan-out protocol from agent-docs b0d76d8 ([c9958a5](https://github.com/rknightion/transceiver-exporter/commit/c9958a5911d89741fdab106cab3e1d7f63853ec1))
* **readme:** lead with what the project is ([508017f](https://github.com/rknightion/transceiver-exporter/commit/508017fb86f1c87525ac047f982ba5b8d780efee))
* repoint moved repo references after the org consolidation ([060ed20](https://github.com/rknightion/transceiver-exporter/commit/060ed20713f7139830f80298fe4c523cba6ea8b5))
* sync agent-docs, a wave's launch message is a file not a chat block ([9e796ac](https://github.com/rknightion/transceiver-exporter/commit/9e796ac6f96e4d3205275e3d0424bf09f1c5897c))
* sync Astra routing and default wave reports to files ([428066f](https://github.com/rknightion/transceiver-exporter/commit/428066f9cc485fecbdb27083c96ab7bc3fea7f71))
* sync nineteen-worker Codex fan-out guidance ([cc0f88b](https://github.com/rknightion/transceiver-exporter/commit/cc0f88be9af634a23628a5ff3287dc285dee2d9d))
* sync wave-root stage authority and lab-Mac GUI gate ([8c9a88d](https://github.com/rknightion/transceiver-exporter/commit/8c9a88d957b440ce9ed10303cefc0cf31326d15e))
* **tracker:** align canonical fan-out protocol ([b571e5e](https://github.com/rknightion/transceiver-exporter/commit/b571e5ea9534c3e65ee3215366d0f897bfef490d))
* **tracker:** correct the canonical owner in the rendered header ([b209dfd](https://github.com/rknightion/transceiver-exporter/commit/b209dfd92630c8d0b8403be6f1354f76467ebbd1))
* **tracker:** re-import the fan-out protocol from canonical ([6574322](https://github.com/rknightion/transceiver-exporter/commit/65743225a5da262349a1aac454bc22936fb481bf))
* **tracker:** render agent documents from the canonical source ([52d93a6](https://github.com/rknightion/transceiver-exporter/commit/52d93a611151932adbc88d62861d22b5d495c2e4))


### Build & CI

* add auto-rc, arm-automerge and ghcr-cleanup ([1b1bd42](https://github.com/rknightion/transceiver-exporter/commit/1b1bd4239b3f1969bef949bb6ccddf765c7706e6))
* **auto-rc:** trigger on CI completion instead of push ([25a860d](https://github.com/rknightion/transceiver-exporter/commit/25a860dc1e93d47eb25fdc1a0f7fe7e08339a5af))
* bump the broker-token action pin ([cf476b6](https://github.com/rknightion/transceiver-exporter/commit/cf476b6ebf6181f99d4717979d015ce42e483b8f))
* docs-sync and grafana-sync targets moved to the m7kni org ([a5c6e9e](https://github.com/rknightion/transceiver-exporter/commit/a5c6e9eb90ee982591d43509b1529ecd6d6268f9))
* **release:** repin shared binaries workflow, grant attestations: write ([e7037d5](https://github.com/rknightion/transceiver-exporter/commit/e7037d5c552d14e7c605f350e2d956afea90859d))
* repin the release-automation reusables to v1.8.0 ([408b40f](https://github.com/rknightion/transceiver-exporter/commit/408b40f0de99ff7cc7afe512c8520570a043841f))
* repin the shared reusables to v1.18.1 ([19411d8](https://github.com/rknightion/transceiver-exporter/commit/19411d858c797ef4864116d90836fc273dbc606f))
* ship only linux/amd64 and linux/arm64 ([2b50c83](https://github.com/rknightion/transceiver-exporter/commit/2b50c83f3095de7ad142b1231ff6906c69003658))
* upgrade to Go 1.27 ([2d6f5fc](https://github.com/rknightion/transceiver-exporter/commit/2d6f5fcd972a185f0438be54517df0a955453750))

## [1.6.0](https://github.com/rknightion/transceiver-exporter/compare/v1.5.1...v1.6.0) (2026-07-03)


### Features

* use regex to include/exclude interfaces ([003b14d](https://github.com/rknightion/transceiver-exporter/commit/003b14d2f1d5aa127e1d44f36c947e0a9d92b9bd))


### Bug Fixes

* **interfaces-regex:** check for valid regex at start of scrape ([0a5c216](https://github.com/rknightion/transceiver-exporter/commit/0a5c21690378dbbe737bad7684747a2024018265))
* **interfaces-regex:** compile at startup ([e2fd2ee](https://github.com/rknightion/transceiver-exporter/commit/e2fd2eefac7c4681832034578e1f77cc562d1954))
* revert version bump ([b806da7](https://github.com/rknightion/transceiver-exporter/commit/b806da7dc72e7c0da2966cfd9b2df2bc8c158cfc))
* run container as root for ethtool EEPROM access ([8eeb7f3](https://github.com/rknightion/transceiver-exporter/commit/8eeb7f38028c9ab2117fffd39bbc7eb22d7dab99))


### Refactor

* make version variable to support dynamic injection during build ([ffe6b2b](https://github.com/rknightion/transceiver-exporter/commit/ffe6b2b03c5fb0fb6efe3e915faff676e0f61b9f))


### Documentation

* reorganize and expand Docker Compose configuration documentation ([fafeab3](https://github.com/rknightion/transceiver-exporter/commit/fafeab3079d8c4962204a82ec3633776c6d95e3b))


### Build & CI

* add 386 and ARM/v7 architecture builds to Docker publish workflow ([c43e551](https://github.com/rknightion/transceiver-exporter/commit/c43e551eb7a0bf1c3a8033f747b96477f8895e13))
* add Codacy coverage upload and repo-local path excludes ([92a1430](https://github.com/rknightion/transceiver-exporter/commit/92a1430ee46e3e6fdde2c064d99b26a2fd5c1d76))
* add Docker containerization support ([fdbdb23](https://github.com/rknightion/transceiver-exporter/commit/fdbdb23eb62711091ae7477f644e3dabb4f72b7b))
* add GitHub Actions workflow for multi-architecture Docker image publishing ([04838eb](https://github.com/rknightion/transceiver-exporter/commit/04838eb82d0fc0168b083be437d922a070fc98f3))
* add OpenSSF Scorecard via shared reusable workflow ([77b8e47](https://github.com/rknightion/transceiver-exporter/commit/77b8e4724b914255fc4dd2caf27a4babcc53e1b9))
* add Snyk monitor; credit original wobcom authors ([a54906e](https://github.com/rknightion/transceiver-exporter/commit/a54906e6c6447aaeaee0c9828fd055bd6f199d02))
* adopt shared container-publish + snyk reusables; de-fork module path ([53663a9](https://github.com/rknightion/transceiver-exporter/commit/53663a951f31c9fe86fd7af1b182ce9e40199c38))
* bump shared rknightion reusables to v1.3.1 ([d196d9a](https://github.com/rknightion/transceiver-exporter/commit/d196d9a6a38afa8c1a200e670c0f32e3f91c0acd))
* drop CodeQL pull_request trigger to trim Actions fan-out ([3cce558](https://github.com/rknightion/transceiver-exporter/commit/3cce558eee4c3aee548ba997f34e2ea9e6fe5bfa))
* Drop support for 386 architecture ([249e259](https://github.com/rknightion/transceiver-exporter/commit/249e259d0b1f22dcfe7b2ab2d90bc802250605c1))
* Initialize GH Actions ([9d40b48](https://github.com/rknightion/transceiver-exporter/commit/9d40b485760a93575a5b845f6f0b02a960f29ecc))
* keyless-sign release binaries (supply-chain parity) ([935547f](https://github.com/rknightion/transceiver-exporter/commit/935547f41781b9ce511645e10ddc01b7a3384685))
* migrate release binaries to GoReleaser via shared binaries reusable ([add6885](https://github.com/rknightion/transceiver-exporter/commit/add6885b14c63382a5cb53e1e3333e3781783ad3))
* pin shared rknightion reusables to v1.0.0 ([70ae3df](https://github.com/rknightion/transceiver-exporter/commit/70ae3df15407f8a807b9bf69d7831d89db74d0e1))
* refine Docker image tagging strategy for different build contexts ([ec27d6c](https://github.com/rknightion/transceiver-exporter/commit/ec27d6c3ba47883109c1b5e4975dec2c888ef0d6))
* Remove GitLab CI residue ([cc93401](https://github.com/rknightion/transceiver-exporter/commit/cc93401722b46c301cdc2e873fba7ca0c8d5cd09))
* remove notify-maintainer-on-new-issue workflow ([b14eaba](https://github.com/rknightion/transceiver-exporter/commit/b14eaba4671639aae5eb0f1c746ae13ebd9ea80b))
* Run lint checks on every push ([0c52d1e](https://github.com/rknightion/transceiver-exporter/commit/0c52d1e4188a1b2cf2f339ad612f66b9662e3a3a))

## 1.5.1 - 2025-03-03
### Changes
* Fixed a bug when setting include / exclude regex to nil
* Added more logging and minor refactors

## 1.5.0 - 2024-05-06
### Changes
* Added the option to include and exclude interfaces with regex
  * Thanks @rwxd for this contribution! 

## 1.4.1 - 2023-08-01
### Changes
* --version now returns the correct version

## 1.4.0 - 2023-07-10
### Changes
* Added the option to exclude admin down interfaces
  * `-exclude.interfaces-down`
  * Thanks @4xoc for this contribution! 

## 1.3.0 - 2023-06-14
### Changes
* Added the option to include specific interfaces
  * `-include.interfaces`
  * Thanks @SRv6d for this contribution! 

## 1.2.0 - 2023-06-07
### Changes
* Added arm binary to CI

## 1.1.0 - 2022-08-25
### Changes
* Added the option to export optical power in dBm instead of mW
  * `-collector.optical-power-in-dbm`
  * Thanks for @BarbarossaTM (Cloudflare) for contributing this feature.
* Updated dependencies
  * Switched from deprecated `prometheus/common/log` to `sirupsen/logrus`

### Notes
* For this release we moved the repository from GitLab to GitHub.

## 1.0.1 - 2020-07-14
### Changes
* Switched to GoLang compliant versioning scheme
* Fixed a bug where the scrape would fail due to reading bad data

## 1.0 -  2020-07-13
### Changes
* Initial release
