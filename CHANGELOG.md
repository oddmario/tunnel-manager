# Changelog

## v1.1.9

### Changed
* In case of `route_all_traffic_through_tunnel`: Tunnel Manager no longer replaces the default route; it adds a new route with the highest priority instead.

## v1.1.8

### Changed
* Don't set `CGO_ENABLED` to zero while compiling.
* Updated the `build.py` builder to generate more optimised binaries

## v1.1.7

### Added
* Docker image is now built and pushed to [GitHub Container Registry](https://ghcr.io/oddmario/tunnel-manager) ([#1](https://github.com/oddmario/tunnel-manager/pull/1))
* Added Docker Compose example file ([#1](https://github.com/oddmario/tunnel-manager/pull/1))
* Added a workflow for making releases. ([#1](https://github.com/oddmario/tunnel-manager/pull/1))
* Added a workflow for compiling commits and pull requests