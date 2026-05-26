package faynosync

import (
	"encoding/json"
	"sort"
	"strings"
)

// CheckOptions contains the typed parameters used to check for updates.
//
// Channel, Platform, and Arch are intentionally user-controlled values. The SDK
// does not detect, normalize, remap, or default them.
type CheckOptions struct {
	Owner   string
	AppName string
	Version string

	Channel  string
	Platform string
	Arch     string

	// DeviceID optionally enables server-side telemetry when supported by the API.
	// When empty, the X-Device-ID header is omitted.
	DeviceID string
}

// UpdateResponse contains the typed faynoSync update check response.
type UpdateResponse struct {
	UpdateAvailable        bool   `json:"update_available"`
	UpdateURL              string `json:"update_url,omitempty"`
	Changelog              string `json:"changelog,omitempty"`
	Critical               bool   `json:"critical,omitempty"`
	IsIntermediateRequired bool   `json:"is_intermediate_required,omitempty"`
	PossibleRollback       bool   `json:"possible_rollback,omitempty"`

	// PackageURLs contains package-specific URLs decoded from fields such as
	// update_url_deb, update_url_rpm, or any future update_url_<package> key.
	PackageURLs []PackageUpdateURL `json:"-"`

	// Source identifies whether the response came from the edge or API fallback.
	Source UpdateSource `json:"-"`
}

// PackageUpdateURL contains one package-specific update URL.
type PackageUpdateURL struct {
	Package string
	URL     string
}

// UnmarshalJSON decodes fixed response fields and dynamic update_url_<package>
// fields into a typed representation.
func (r *UpdateResponse) UnmarshalJSON(data []byte) error {
	type responseAlias UpdateResponse

	var fixed responseAlias
	if err := json.Unmarshal(data, &fixed); err != nil {
		return err
	}

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	packageURLs := make([]PackageUpdateURL, 0)
	for key, raw := range fields {
		if !strings.HasPrefix(key, "update_url_") {
			continue
		}

		var updateURL string
		if err := json.Unmarshal(raw, &updateURL); err != nil {
			return err
		}

		packageURLs = append(packageURLs, PackageUpdateURL{
			Package: strings.TrimPrefix(key, "update_url_"),
			URL:     updateURL,
		})
	}

	sort.Slice(packageURLs, func(i, j int) bool {
		return packageURLs[i].Package < packageURLs[j].Package
	})

	*r = UpdateResponse(fixed)
	r.PackageURLs = packageURLs

	return nil
}

// UpdateSource identifies where an update response was loaded from.
type UpdateSource int

const (
	// SourceUnknown indicates that the response source is unknown.
	SourceUnknown UpdateSource = iota

	// SourceEdge indicates that the response came from the configured EdgeURL.
	SourceEdge

	// SourceAPI indicates that the response came from the configured BaseURL API.
	SourceAPI
)
