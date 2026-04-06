---
status: proposed
date: 2026-03-30
decision: Use a graph-first, create-only migration framework to namespace actions, subject mappings, subject condition sets, and registered resources
author: '@ehealy'
deciders: []
---

# Policy graph namespace migration

## Context and Problem Statement

The opentdf platform recently made actions, subject mappings (SM), subject condition sets (SCS), and registered resources (RR) namespace-scoped. Previously only attributes, obligations, and resource mappings were namespaced. A migration path is needed in otdfctl to move existing unnamespaced policy into the correct namespaces while respecting cross-object constraints.

An RR migration already exists (`migrations/registered-resources.go`) that uses a delete+recreate pattern. This ADR covers expanding migration to all newly namespaced types and changing the migration strategy to be safer and more flexible.

### Key constraints

- All policy referenced by another policy item must be in the same namespace. If an SM references an attribute value, an action, and an SCS, all must share the SM's namespace.
- Obligations and obligation triggers are already namespaced but their action references must be validated for namespace consistency.
- Standard actions (create, read, update, delete) must exist in every namespace. A separate platform DB migration will seed them into existing namespaces retroactively.

## Decision Drivers

- **Safety**: migration should not destroy data. If something goes wrong mid-run, rollback should be trivial.
- **Conflict visibility**: cross-namespace reference conflicts must be surfaced before any writes happen, regardless of which type is being migrated.
- **Incremental adoption**: operators should be able to migrate one type at a time or all at once.
- **Auditability**: it should be clear which items were migrated, what was created, and what is safe to clean up.

## Considered Options

1. **Formal graph library with node/edge types** — build a full graph data structure with typed nodes and edges, a constraint engine, and a planner that traverses the graph.
2. **Lightweight reference collection with conflict detection** — extend the existing pattern from the RR migration: fetch all policy, walk references, detect mismatches, produce an ordered operation plan. No formal graph types.
3. **Per-type independent migrations** — separate, standalone migration commands for each type with no shared analysis.

## Decision Outcome

Option 2: lightweight reference collection with conflict detection. This reuses the proven patterns from the existing RR migration and avoids the complexity of a formal graph library while still providing the graph-first analysis (all references are collected and validated before any writes).

### Migration strategy: create-only (no deletes)

Migration is purely additive. It creates namespaced copies of policy items but never deletes the old unnamespaced originals. After creating the namespaced copy, the old item's metadata is updated with `wasMigrated=true`. Cleanup is handled by separate `migrate prune` subcommands after the operator verifies correctness.

This replaces the current RR migration's delete+recreate pattern.

### Namespace determination rules

- **SM**: determined by its attribute value's namespace (already namespaced, so deterministic).
- **SCS**: determined by the SMs that reference it. If all referencing SMs target the same namespace, the SCS goes there. If multiple namespaces, the SCS is cloned per namespace. Orphaned SCS are skipped.
- **Actions**: global actions referenced from one namespace are created there. Referenced from multiple namespaces, created per namespace. If an action with the same name already exists in the target namespace, it is reused. Standard actions are skipped (assumed seeded by the platform).
- **RR**: determined by inspecting action-attribute-value (AAV) references, same as existing migration logic.

### Execution order

1. **Actions** — create namespaced actions (no dependencies)
2. **SCS** — clone into target namespaces (depends on knowing SM target NS, which is data-only)
3. **SM** — create new in target NS with rewritten action/SCS refs (depends on actions + SCS existing)
4. **RR** — create new in target NS with values (existing logic, updated to create-only)

### Command UX

Unified command runs all types in dependency order:

```
otdfctl migrate policy-graph                                          # dry-run, all types
otdfctl migrate policy-graph --commit --output=migration.json         # commit all, write manifest
otdfctl migrate policy-graph --commit --interactive                   # per-item confirmation
otdfctl migrate policy-graph --scope=actions                          # limit to actions only
otdfctl migrate policy-graph --scope=actions,sm                       # actions + SM (SCS implied)
```

Per-type subcommands each run the full planner but execute only their type:

```
otdfctl migrate actions --commit --output=actions.json
otdfctl migrate subject-condition-sets --commit --output=scs.json
otdfctl migrate subject-mappings --commit --output=sm.json
otdfctl migrate registered-resources --commit --output=rr.json
```

All commands use the same `BuildPolicyGraphPlan()` to analyze the full policy graph. Per-type subcommands display the full plan but only execute their scope. This ensures cross-type conflicts are always visible.

Prune subcommands clean up after migration:

```
otdfctl migrate prune actions --from=migration.json --commit
otdfctl migrate prune subject-condition-sets --from=migration.json --commit
otdfctl migrate prune subject-mappings --from=migration.json --commit
otdfctl migrate prune registered-resources --from=migration.json --commit
```

Prune uses the manifest (`--from`) for exact old->new ID matching, or falls back to the `wasMigrated=true` metadata label if no manifest is provided.

### Migration manifest

The `--output` flag writes a JSON manifest of all old->new ID mappings:

```json
{
  "timestamp": "2026-03-30T...",
  "actions": [
    {"old_id": "...", "new_id": "...", "name": "approve", "target_namespace": "https://example.com"}
  ],
  "subject_condition_sets": [
    {"old_id": "...", "new_id": "...", "target_namespace": "https://example.com"}
  ],
  "subject_mappings": [
    {"old_id": "...", "new_id": "...", "target_namespace": "https://example.com"}
  ],
  "registered_resources": [
    {"old_id": "...", "new_id": "...", "name": "my-resource", "target_namespace": "https://example.com"}
  ]
}
```

### Blocker types

Action blockers:
- `MISSING_ACTION_ID` — action reference has no ID
- `PARENT_NAMESPACE_REQUIRED` — parent object is unnamespaced (except legacy RR values)
- `ACTION_NAMESPACE_MISMATCH` — action already namespaced but in wrong NS
- `ATTRIBUTE_NAMESPACE_MISMATCH` — attribute value NS differs from parent NS

SM/SCS blockers:
- `SM_NAMESPACE_UNDETERMINED` — attribute value has no parseable namespace FQN
- `SM_ACTION_NOT_IN_PLAN` — SM references action not covered by action plan
- `SCS_NAMESPACE_CONFLICT` — SCS clone target cannot be determined

### Prerequisites

- Platform DB migration to seed standard actions (create, read, update, delete) into all existing namespaces. New namespaces are already seeded via `CreateNamespace()`, but existing namespaces need a retroactive seed. This is a separate platform task.

### PR breakdown

1. **Unified Planner (read-only)**: planning logic, display, expanded `MigrationHandler` interface, `migrate policy-graph` and per-type subcommands wired for dry-run only.
2. **Create-Only Execution**: action create + SCS clone + SM create execution, `--commit`, `--scope`, `--interactive`, `--output` manifest.
3. **RR Migration Update + Prune Commands**: update RR migration to create-only, prune logic for all types, `migrate prune <type>` subcommands.

### Consequences

- **Positive**:
  - No data loss risk during migration — purely additive, old items preserved.
  - `wasMigrated=true` label makes migrated items visible to operators and tooling.
  - Manifest file provides auditable record of what was migrated.
  - Full graph analysis runs regardless of scope, so cross-type conflicts are always surfaced.
  - Incremental migration is supported via per-type subcommands and `--scope`.
  - Prune is a separate deliberate step, reducing accidental cleanup.
- **Negative**:
  - During the window between migration and prune, both old and new copies exist. Operators must understand that the old unnamespaced items are stale.
  - Requires a platform DB migration as a prerequisite for standard action seeding.
  - The `MigrationHandler` interface grows to accommodate all policy types, which increases the mock surface in tests.
