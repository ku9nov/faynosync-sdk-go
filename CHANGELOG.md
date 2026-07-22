# Changelog

## v0.3.0

### Added

- Staged (canary) rollout support in `CheckForUpdates`: when `/checkVersion` returns a `rollout` object (`percent`, `seed`), the SDK computes a deterministic, sticky per-device bucket (`sha256(deviceID + ":" + seed)` → first 8 bytes big-endian uint64 `% 100`) and includes the install only when `bucket < percent`. Excluded installs get `UpdateAvailable: false` with `UpdateURL`/`PackageURLs` cleared so the gate cannot be bypassed. Exposes the decision via the new `RolloutInfo` (`Rollout` on `UpdateResponse`) and the `RolloutBucket` helper. Requires `DeviceID`; without it the install stays out of the rollout. Works in edge/CDN mode too

## v0.2.0

### Added

- "Updater" field to edge path.


## v0.1.0

Initial release of the faynoSync Go SDK (transport layer only).

### Added

- `Client` with `CheckForUpdates` against the Base API (`GET /checkVersion`)
- Optional `EdgeURL` static JSON lookup with automatic API fallback
- Typed `CheckOptions`, `UpdateResponse`, and package URL decoding
- Sentinel errors and `EndpointError` for request failures
- Optional `SystemPlatform` / `SystemArch` helpers
- Examples: basic, edge fallback, custom HTTP client
- Add telemetry beacon to the edge response
