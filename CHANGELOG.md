# Changelog

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
