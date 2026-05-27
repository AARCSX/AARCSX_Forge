# Shared Package Governance

`internal/shared` is reserved for non-domain primitives only:

- utility helpers
- generic constants
- low-level reusable abstractions
- pure formatting/parsing helpers

Never place:

- business logic
- tenant workflows
- authorization policies
- cross-module orchestration

If code has domain behavior, it belongs in its owning module package.
