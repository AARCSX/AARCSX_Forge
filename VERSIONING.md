# API Versioning Strategy

## Overview
AARCSX Forge uses URI versioning for its public APIs. The version is included in the API path, e.g., `/api/v1/resource`.

## Current Version
The current stable version is `v1`, accessible under the `/api/v1` path.

## Versioning Policy
- **Backward Compatibility**: Within a major version (e.g., v1), we strive to maintain backward compatibility. Breaking changes are only introduced in a new major version (e.g., v2).
- **Deprecation Policy**: When a version is deprecated, it will be supported for a minimum of 6 months after the release of the next major version. Deprecation announcements will be made at least 3 months in advance.
- **Version Support**: We typically support the current stable version and the previous major version (n-1) for critical bug fixes and security patches.

## Implementation
- All API routes are defined under a versioned group in the API server (see `cmd/api/main.go`).
- The version is part of the URL path, making it explicit and cache-friendly.
- Internal interfaces and contracts are not versioned; they evolve with the codebase.

## Future Versions
When a new major version is introduced (e.g., v2):
- A new route group will be created: `/api/v2/...`
- The previous version (v1) will continue to be maintained according to the deprecation policy.
- Documentation will be updated to reflect the new version and migration guidelines.