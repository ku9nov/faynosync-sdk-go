# Changelog

## v0.1.0

Initial release of the faynoSync Go SDK (transport layer only).

### Added

- `Client` with `CheckForUpdates` against the Base API (`GET /checkVersion`)
- Optional `EdgeURL` static JSON lookup with automatic API fallback
- Typed `CheckOptions`, `UpdateResponse`, and package URL decoding
- Sentinel errors and `EndpointError` for request failures
- Optional `SystemPlatform` / `SystemArch` helpers
- Examples: basic, edge fallback, custom HTTP client
