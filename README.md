# faynoSync Go SDK

Production-oriented Go SDK for checking application updates with faynoSync.

This package is a small typed transport and developer experience layer. It does not implement update installation, platform normalization, metadata verification, caching, or business rules.

## Installation

```sh
go get github.com/ku9nov/faynosync-sdk-go
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"

	faynosync "github.com/ku9nov/faynosync-sdk-go"
)

func main() {
	client := faynosync.NewClient(faynosync.Config{
		BaseURL: "https://api.example.com",
	})

	resp, err := client.CheckForUpdates(context.Background(), faynosync.CheckOptions{
		Owner:    "admin",
		AppName:  "test",
		Version:  "0.0.0.5",
		Channel:  "nightly",
		Platform: "darwin",
		Arch:     "arm64",
	})
	if err != nil {
		log.Fatal(err)
	}

	if resp.UpdateAvailable {
		if resp.UpdateURL != "" {
			fmt.Printf("Update is available: %s\n", resp.UpdateURL)
		}
		for _, packageURL := range resp.PackageURLs {
			fmt.Printf("%s update is available: %s\n", packageURL.Package, packageURL.URL)
		}
	}
}
```

## Configuration

```go
client := faynosync.NewClient(faynosync.Config{
	BaseURL: "https://api.example.com",
	EdgeURL: "https://cdn.example.com",
	HTTPClient: &http.Client{
		Timeout: 10 * time.Second,
	},
})
```

`BaseURL` is required. It points to the faynoSync API.

`EdgeURL` is optional. When configured, the SDK tries a static edge JSON response before falling back to the API.

`HTTPClient` is optional. When omitted, the SDK creates an `http.Client` with a default timeout. Custom clients are useful for timeouts, proxies, custom transports, and connection pooling policies.

The client is safe for concurrent use.

## Update Checks

`CheckForUpdates` sends a context-aware request and returns a typed response:

```go
resp, err := client.CheckForUpdates(ctx, faynosync.CheckOptions{
	Owner:    "admin",
	AppName:  "test",
	Version:  "0.0.0.5",
	Channel:  "nightly",
	Platform: "darwin",
	Arch:     "arm64",
	DeviceID: "optional-device-id",
})
```

`DeviceID` is optional. When set, the SDK sends it as the `X-Device-ID` header.

## Base API Request

The BaseURL API request uses `GET /checkVersion`:

```text
GET /checkVersion?app_name=test&version=0.0.0.5&channel=nightly&platform=darwin&arch=arm64&owner=admin
X-Device-ID: optional
```

Query parameters are built from `CheckOptions` with typed fields. The SDK does not use untyped maps in its public API.

## Response Model

faynoSync may return a direct binary update URL:

```json
{
  "update_available": true,
  "update_url": "https://downloads.example.com/app"
}
```

It may also return package-specific URLs with dynamic field names:

```json
{
  "update_available": true,
  "update_url_deb": "https://downloads.example.com/app.deb",
  "update_url_rpm": "https://downloads.example.com/app.rpm",
  "changelog": "### Changelog\n\n- Added feature X",
  "critical": true,
  "is_intermediate_required": true,
  "possible_rollback": true
}
```

The SDK decodes these into a typed response:

```go
if resp.UpdateURL != "" {
	fmt.Println(resp.UpdateURL)
}

for _, packageURL := range resp.PackageURLs {
	fmt.Println(packageURL.Package, packageURL.URL)
}
```

## EdgeURL Fallback

When `EdgeURL` is configured, the SDK first tries a static JSON response:

```text
GET /responses/{owner}/{app_name}/{channel}/{platform}/{arch}/{version}.json
```

For example:

```text
GET /responses/admin/test/nightly/darwin/arm64/0.0.0.5.json
```

If the edge response succeeds with HTTP 200 and valid JSON, `UpdateResponse.Source` is `SourceEdge`.

The SDK falls back to the BaseURL API when the edge request has:

- a network error;
- a timeout;
- invalid JSON;
- HTTP 404;
- any other non-200 response.

If the fallback API succeeds, `UpdateResponse.Source` is `SourceAPI`.

## Platform, Channel, And Architecture Values

faynoSync supports fully custom platform, channel, and architecture values. This SDK never normalizes or remaps them.

The SDK will not change values such as:

- `macos` to `darwin`;
- `osx` to `darwin`;
- `stable` to `default`.

Whatever string you pass in `CheckOptions.Channel`, `CheckOptions.Platform`, and `CheckOptions.Arch` is the string sent to faynoSync.

## Optional System Helpers

The SDK provides optional helpers:

```go
platform := faynosync.SystemPlatform() // runtime.GOOS
arch := faynosync.SystemArch()         // runtime.GOARCH
```

These helpers are never used automatically. Use them only when Go runtime values match your faynoSync configuration.

## Error Handling

The SDK validates required fields and returns typed sentinel errors:

```go
resp, err := client.CheckForUpdates(ctx, opts)
if err != nil {
	switch {
	case errors.Is(err, faynosync.ErrMissingBaseURL):
		// Configure Config.BaseURL.
	case errors.Is(err, faynosync.ErrMissingOwner):
		// Set CheckOptions.Owner.
	case errors.Is(err, faynosync.ErrMissingAppName):
		// Set CheckOptions.AppName.
	case errors.Is(err, faynosync.ErrMissingVersion):
		// Set CheckOptions.Version.
	case errors.Is(err, faynosync.ErrRequestFailed):
		// Inspect the wrapped endpoint error.
	default:
		// Handle any other error.
	}
}
```

Request failures preserve underlying causes for `errors.Is` and `errors.As`:

```go
var endpointErr *faynosync.EndpointError
if errors.As(err, &endpointErr) {
	fmt.Println(endpointErr.URL)
	fmt.Println(endpointErr.StatusCode)
}
```

## Examples

Runnable examples are available in:

- `examples/basic`
- `examples/edge-fallback`
- `examples/custom-http-client`

## Security Scope

This SDK version performs update-check transport requests and typed response decoding only.

It does not verify TUF metadata, signatures, thresholds, expiration, rollback protection, or cache safety. Applications that need secure update metadata verification must perform that verification in the appropriate faynoSync component or a future SDK layer that explicitly implements it.

No signature, threshold, expiration, rollback, freeze, root-of-trust, or cache protection is weakened by this transport-only SDK.
