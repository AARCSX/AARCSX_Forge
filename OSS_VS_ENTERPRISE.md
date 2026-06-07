# OSS vs Enterprise Boundary

## Overview
This document defines the clean boundaries between Open Source Software (OSS) and Enterprise features in AARCSX Forge to enable a future community/enterprise split.

## Current Structure
All code in the `internal/` directory is considered OSS/core functionality. Enterprise-specific features would be placed in:

```
internal/
├── enterprise/           # Enterprise-only features
│   ├── billing/          # Billing and subscription management
│   ├── ldap/             # LDAP/Active Directory integration
│   ├── sso/              # Advanced SSO (SAML, OIDC providers)
│   ├── audit/            # Advanced audit and compliance features
│   └── support/          # Priority support features
├── oss/                  # OSS-specific features (if any separation needed)
└── ...                   # Core shared functionality (remains in current structure)
```

## Boundary Rules
1. **Core Functionality**: All authentication, tenant management, storage, notifications, observability, and core API functionality remains in the current internal/* packages as OSS.

2. **Enterprise Features**: Features that are typically enterprise-only (like advanced compliance, premium support, advanced integrations) go in `internal/enterprise/`.

3. **Dependencies**: Enterprise packages must not be depended upon by core OSS packages. Dependencies flow from enterprise -> core, not vice versa.

4. **Configuration**: Enterprise features are enabled via configuration flags and feature flags.

5. **Documentation**: Enterprise features are documented separately from OSS features.

## Implementation Approach
When enterprise features are needed:
1. Create the feature under `internal/enterprise/{feature}/`
2. Ensure the feature can be disabled via configuration
3. Document the feature as enterprise-only
4. Update any dependency injection to conditionally include enterprise features

This structure allows for a clean separation where:
- The OSS distribution includes everything except `internal/enterprise/`
- An Enterprise distribution includes both core and enterprise features
- Users can upgrade from OSS to Enterprise by adding the enterprise package and configuring it