# Policy Graph Migration Plan

## Context

The opentdf platform recently namespaced actions, subject mappings (SM), subject condition sets (SCS), and registered resources (RR). RR migration already exists in `migrations/registered-resources.go`. This plan expands migration to cover actions, SMs, and SCSs using a graph-first approach that detects cross-namespace conflicts before committing any changes.

**Key constraint**: all policy referenced by another policy item must be in the same namespace. If an SM references an attribute value, an action, and an SCS — all must share the SM's namespace.

**Obligations/triggers** are already namespaced and not migrated, but their action references are validated to ensure namespace consistency.

**Prerequisite**: a platform DB migration must seed standard actions (create, read, update, delete) into all existing namespaces. New namespaces already get seeded via `CreateNamespace()`, but existing namespaces created before action namespacing need a retroactive seed. This is a separate platform task.

## Approach

Build a lightweight "collect references -> detect conflicts -> plan operations" pattern rather than a formal graph library. A unified `PolicyGraphPlan` composes domain-specific analysis into a single ordered plan. This follows the same approach used in the existing RR migration (`registered-resources.go`) — fetch all relevant policy, walk references, detect mismatches, produce an operation plan.

### Create-only migration (no deletes)

Migration is purely additive — it creates namespaced copies of policy items but **never deletes** the old unnamespaced originals. This is safer than delete+recreate: if something goes wrong mid-migration, no data is lost. Cleanup of old unnamespaced items is handled by separate `migrate prune` subcommands run after verification.

After creating the namespaced copy, the migration updates the **old** (unnamespaced) item's metadata to add `wasMigrated=true`. This labels old items in-place so they're easily identifiable — both for operators inspecting policy and for the prune command.

This is a change from the current RR migration behavior (which deletes the old resource after creating the new one). The RR migration will be updated to follow this create-only pattern as well.

## Namespace Determination Rules

- **SM**: determined by its attribute value's namespace (attribute values are already namespaced, so this is deterministic)
- **SCS**: determined by the SMs that reference it. If all SMs target the same NS -> same NS. If multiple NS -> clone SCS per NS. Orphaned SCS (unreferenced) -> skip, user can prune separately.
- **Actions**: global actions referenced from one NS -> create in that NS. Referenced from multiple NS -> create per NS. If action with same name already exists in target NS -> reuse it. Standard actions (create/read/update/delete) are assumed to already exist per namespace via platform DB migration.

## Migration Strategy per Type

- **Actions**: Create namespaced actions in each required NS. Reuse if name already exists in target NS. Standard actions are skipped (assumed seeded by platform).
- **SCS**: Clone (create new copy with same subject sets content) into target NS.
- **SM**: Create new SM in target NS with correct action refs, attribute value ref, and SCS ref. Old unnamespaced SM is left in place.
- **RR**: Create new RR + values in target NS (updated from current delete+recreate pattern). Old unnamespaced RR is left in place.

## Execution Order

1. **Actions** — create namespaced actions (no dependencies)
2. **SCS** — clone into target namespaces (depends on knowing SM target NS, which is data-only)
3. **SM** — create new in target NS with rewritten action/SCS refs (depends on actions + SCS existing in target NS)
4. **RR** — create new in target NS with values (existing migration, updated to create-only pattern)

## Command UX

### Unified command (runs all types in dependency order)

```
otdfctl migrate policy-graph                                          # dry-run, all types
otdfctl migrate policy-graph --commit --output=migration.json         # commit all types, write manifest
otdfctl migrate policy-graph --commit --interactive                   # per-item confirmation
otdfctl migrate policy-graph --scope=actions                          # actions only
otdfctl migrate policy-graph --scope=actions,sm                       # actions + SM (+ SCS, implied by SM dependency)
```

`--output` writes a JSON manifest of all old->new ID mappings, which `migrate prune` can consume.

### Per-type subcommands (each runs the full planner but executes only its type)

```
otdfctl migrate actions                                               # dry-run actions
otdfctl migrate actions --commit --output=actions.json                # commit, write manifest
otdfctl migrate subject-condition-sets --commit --output=scs.json     # commit SCS clones
otdfctl migrate subject-mappings --commit --output=sm.json            # commit SM creates
otdfctl migrate registered-resources --commit --output=rr.json        # commit RR creates (no longer deletes old)
```

All commands — unified and per-type — use the same `BuildPolicyGraphPlan()` planner to analyze the full policy graph. Per-type subcommands display the full plan but only execute operations for their type. This ensures cross-type conflicts are always visible even when migrating one type at a time.

### Prune subcommands (cleanup after migration)

```
otdfctl migrate prune actions --from=migration.json                # dry-run: show what would be deleted
otdfctl migrate prune actions --from=migration.json --commit       # delete old unnamespaced actions
otdfctl migrate prune subject-condition-sets --from=migration.json --commit
otdfctl migrate prune subject-mappings --from=migration.json --commit
otdfctl migrate prune registered-resources --from=migration.json --commit
```

Prune reads the migration manifest (JSON file emitted by `--output`) to know exactly which old unnamespaced items map to which new namespaced copies. If `--from` is not provided, prune falls back to re-deriving matches (by name for actions/RR, by content for SCS, by attribute-value-id+actions+SCS for SM). Items that were never migrated are left untouched and flagged in the output.

## Files to Create

### `migrations/policy_graph_plan.go`

Unified planner that collects all policy references and produces a migration plan.

**Key types:**

```go
// Top-level plan containing all analysis results
type PolicyGraphPlan struct {
    ActionPlan    *ActionNamespacePlan
    SMPlans       []SubjectMappingMigrationPlan
    SCSPlans      []SCSMigrationPlan
    RRPlans       []RegisteredResourceMigrationPlan  // reuse existing type
    Blockers      []PolicyGraphBlocker
    Operations    []PolicyGraphOperation
    OrphanedSCS   []string                           // SCS IDs not referenced by any SM
}

// Per unnamespaced SM
type SubjectMappingMigrationPlan struct {
    SM                 *policy.SubjectMapping
    TargetNamespaceFQN string                // derived from attribute value
    TargetNamespaceID  string
    ActionRewrites     []ActionRewrite       // global action -> namespaced action
    SCSRewrite         *SCSRewrite           // old SCS -> cloned SCS in target NS (nil if already correct)
}

// Per SCS that needs to be cloned into a new namespace
type SCSMigrationPlan struct {
    SourceSCS          *policy.SubjectConditionSet
    TargetNamespaceFQN string
    TargetNamespaceID  string
    Operation          string                // "create_scs_clone" or "reuse_existing_scs"
    ExistingSCSID      string                // populated if reuse
}

// Action analysis results
type ActionNamespacePlan struct {
    TotalReferences int
    Matching        int
    RequiresRewrite int
    Blockers        []ActionNamespaceBlocker
    Operations      []ActionNamespaceOperation
}
```

**Key functions:**

- `BuildPolicyGraphPlan(ctx, handler)` — orchestrates full analysis:
  1. Fetch all SMs, SCSs, namespaces, actions, RRs, obligation triggers (paginated)
  2. Build action namespace plan: collect action refs from RR values, SMs, obligation triggers; detect mismatches; plan action creates + reference rewrites. Skip standard actions (they're seeded by the platform).
  3. For each unnamespaced SM: derive target NS from attribute value FQN using `identifier.Parse`
  4. For each SM's SCS: check if SCS is already in SM's target NS; if not, plan clone/reuse
  5. For each SM's actions: check if in SM's target NS; cross-reference with action plan
  6. Deduplicate SCS clones per `(source_scs_id, target_ns)`
  7. Check for existing actions with same name in target NS -> mark as reuse
  8. Build RR plans for unnamespaced resources (reuse existing `detectRequiredNamespace` logic)
  9. Emit blockers and ordered operations

- Paginated fetch helpers (same pattern as existing `fetchAllResourceValues`):
  - `fetchAllSubjectMappings(ctx, h)`
  - `fetchAllSubjectConditionSets(ctx, h)`
  - `fetchAllActions(ctx, h)`
  - `fetchAllObligationTriggers(ctx, h)`
  - `listAllRegisteredResources(ctx, h)`

- `extractNamespaceFQNFromAttributeValueFQN(valueFQN)` — parses namespace from attribute value FQN using `identifier.Parse[*identifier.FullyQualifiedAttribute]`

**Action namespace analysis logic:**

Collects `actionNamespaceReference` structs from three sources:
1. RR values -> actions via AAVs (parent NS from resource)
2. SMs -> actions (parent NS from SM's determined target namespace)
3. Obligation triggers -> actions (parent NS from obligation chain)

For each reference, checks:
- Action ID present? If not -> `MISSING_ACTION_ID` blocker
- Parent namespaced? If not -> `PARENT_NAMESPACE_REQUIRED` blocker (except legacy unnamespaced RR values)
- Attribute value NS matches parent NS? If not -> `ATTRIBUTE_NAMESPACE_MISMATCH` blocker
- Action global (no namespace)? -> plan `create_action` operation + rewrite
- Action namespace differs from parent? -> `ACTION_NAMESPACE_MISMATCH` blocker + plan create + rewrite

Deduplicates action creates by `(action_id, target_namespace_id)` key.

### `migrations/policy_graph_plan_test.go`

Unit tests with mock handler (extending mock pattern from `registered-resources_test.go`):

- SM with deterministic NS, all refs aligned -> no operations needed
- SM needing action rewrite (global action) -> action create + SM create operations
- SM needing SCS clone (SCS in different NS) -> SCS clone + SM create operations
- SM with unparseable attribute value NS -> `SM_NAMESPACE_UNDETERMINED` blocker
- SCS referenced by SMs in 2+ NS -> deduplicated clone operations
- Orphaned SCS -> reported in `OrphanedSCS`, not migrated
- Action name already exists in target NS -> `reuse_existing_action` operation
- Standard actions skipped (not planned for creation)
- Obligation trigger action NS mismatch -> blocker
- Already-namespaced SM/RR -> skipped (not included in plan)
- No unnamespaced policy -> empty plan, no operations

### `migrations/policy_graph_display.go`

Human-readable dry-run output:
- Action plan summary: total references, matching, rewrites needed, blockers
- SM migration summary: count to migrate, NS assignments, rewrites needed
- SCS migration summary: clones needed, reuses, orphans skipped
- RR migration summary: count to migrate, NS detection results
- Unified blocker list (capped at 25)
- Unified operation list in execution order

### `migrations/policy_graph_execute.go`

Commit-mode execution (create-only, no deletes):

- `ExecutePolicyGraphMigration(ctx, handler, prompter, plan, commit, interactive, scope)`:
  1. Validate no blockers (or only blockers outside scope)
  2. Confirm backup via prompter
  3. If actions in scope: create actions via `handler.CreateAction(name, targetNS, nil)`. For reuse cases, skip creation and record existing action ID.
  4. If SCS in scope: create SCS clones via `handler.CreateSubjectConditionSet(subjectSets, metadata, targetNS)`
  5. If SM in scope: create new SM in target NS via `handler.CreateNewSubjectMapping(attrValID, rewrittenActions, newSCSID, nil, metadata, targetNS)`. Then update old SM metadata to add `wasMigrated=true`. Old SM is NOT deleted.
  6. If RR in scope: create new RR + values in target NS (using existing `commitRegisteredResourceMigration` logic, modified to skip the delete step). Then update old RR metadata to add `wasMigrated=true`. Old RR is NOT deleted.
  7. For actions and SCS: after creating namespaced copies, update old items' metadata to add `wasMigrated=true`.
  8. Print summary (created/skipped/failed counts)
  9. If `--output` specified: write migration manifest JSON to file
- Interactive mode: confirm each operation with user via prompter before executing

**Migration manifest format** (`--output`):
```json
{
  "timestamp": "2026-03-26T...",
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

### `migrations/policy_graph_prune.go`

Prune logic for cleaning up unnamespaced items after migration:

- `PruneUnnamespaced(ctx, handler, policyType, manifestPath, commit)`:
  1. If `manifestPath` provided: load manifest JSON, use old->new ID mappings directly to identify what to delete. Verify the new item still exists before deleting the old one.
  2. If no manifest: list all items of the given type, filter to those with `wasMigrated=true` label and no namespace. This is the primary matching mechanism — no heuristic re-derivation needed since the migration labels old items.
  3. If match found: mark for deletion. If no `wasMigrated` label and no namespace: flag as "not yet migrated, skipping".
  4. If commit: delete matched items. If dry-run: display what would be deleted.

- `LoadMigrationManifest(path)` — reads and validates the manifest JSON file.

### `migrations/policy_graph_prune_test.go`

Tests for prune matching and deletion logic.

## Files to Modify

### `migrations/registered-resources.go`

1. Expand `MigrationHandler` interface with new methods:
```go
ListActions(ctx context.Context, limit, offset int32, namespace string) (*actions.ListActionsResponse, error)
CreateAction(ctx context.Context, name string, namespace string, metadata *common.MetadataMutable) (*policy.Action, error)
ListSubjectConditionSets(ctx context.Context, limit, offset int32, namespace string) (*subjectmapping.ListSubjectConditionSetsResponse, error)
CreateSubjectConditionSet(ctx context.Context, ss []*policy.SubjectSet, metadata *common.MetadataMutable, namespace string) (*policy.SubjectConditionSet, error)
CreateNewSubjectMapping(ctx context.Context, attrValID string, actions []*policy.Action, existingSCSId string, newScs *subjectmapping.SubjectConditionSetCreate, metadata *common.MetadataMutable, namespace string) (*policy.SubjectMapping, error)
UpdateSubjectMapping(ctx context.Context, id string, scsId string, actions []*policy.Action, metadata *common.MetadataMutable, behavior common.MetadataUpdateEnum) (*policy.SubjectMapping, error)
UpdateSubjectConditionSet(ctx context.Context, id string, ss []*policy.SubjectSet, metadata *common.MetadataMutable, behavior common.MetadataUpdateEnum) (*policy.SubjectConditionSet, error)
UpdateAction(ctx context.Context, id string, name string, metadata *common.MetadataMutable, behavior common.MetadataUpdateEnum) (*policy.Action, error)
UpdateRegisteredResource(ctx context.Context, id string, name string, metadata *common.MetadataMutable, behavior common.MetadataUpdateEnum) (*policy.RegisteredResource, error)
DeleteSubjectMapping(ctx context.Context, id string) (*policy.SubjectMapping, error)
DeleteRegisteredResource(ctx context.Context, id string) error  // already exists, needed for prune
```
All methods already exist on `handlers.Handler` — this just adds them to the interface.

2. Update `commitRegisteredResourceMigration` to skip the delete step (create-only). The delete logic moves to the prune command.

### `migrations/registered-resources_test.go`

Expand `MockMigrationHandler` with stub/tracking implementations for new interface methods.

### `cmd/migrate.go`

Add new subcommands:
- `migrate policy-graph` — unified command with `--scope` flag
- `migrate actions` — actions only
- `migrate subject-condition-sets` — SCS only
- `migrate subject-mappings` — SM only
- `migrate prune actions` — prune unnamespaced actions
- `migrate prune subject-condition-sets` — prune unnamespaced SCS
- `migrate prune subject-mappings` — prune unnamespaced SM
- `migrate prune registered-resources` — prune unnamespaced RR

Update existing `migrate registered-resources` to use new create-only behavior.

All share `--commit` and `--interactive` persistent flags from the parent `migrate` command.

### `docs/man/migrate/policy-graph.md` (new)

Man page for the unified command.

## Blocker Types

### Action blockers
- `MISSING_ACTION_ID` — action reference has no ID, cannot safely rewrite
- `PARENT_NAMESPACE_REQUIRED` — parent object is unnamespaced (except legacy RR values)
- `ACTION_NAMESPACE_MISMATCH` — action already namespaced but in wrong NS
- `ATTRIBUTE_NAMESPACE_MISMATCH` — attribute value NS differs from parent NS

### SM/SCS blockers
- `SM_NAMESPACE_UNDETERMINED` — attribute value has no parseable namespace FQN
- `SM_ACTION_NOT_IN_PLAN` — SM references action not covered by action plan
- `SCS_NAMESPACE_CONFLICT` — SCS clone target can't be determined

## PR Breakdown

### PR 1: Unified Planner (read-only)
- `migrations/policy_graph_plan.go` — all planning logic including action analysis
- `migrations/policy_graph_plan_test.go` — unit tests
- `migrations/policy_graph_display.go` — dry-run display
- Expand `MigrationHandler` interface in `registered-resources.go`
- Expand `MockMigrationHandler` in `registered-resources_test.go`
- Wire `migrate policy-graph` and per-type subcommands in `cmd/migrate.go` (dry-run only)

### PR 2: Create-Only Execution (actions, SCS, SM)
- `migrations/policy_graph_execute.go` — action create + SCS clone + SM create execution
- Tests for execution with name collision handling (reuse), SCS deduplication, SM creation
- `--commit` flag wired for all types
- `--scope` flag support
- `--interactive` mode

### PR 3: RR Migration Update + Prune Commands
- Update `commitRegisteredResourceMigration` to skip delete (create-only)
- `migrations/policy_graph_prune.go` — prune logic for all types
- `migrations/policy_graph_prune_test.go` — prune tests
- Wire `migrate prune <type>` subcommands in `cmd/migrate.go`
- Man pages

## Verification

1. `go test ./migrations/...` — unit tests pass
2. `otdfctl migrate policy-graph` against test platform — dry-run shows correct plan
3. `otdfctl migrate policy-graph --commit` — creates namespaced copies of actions, SCS, SM, RR
4. Verify old unnamespaced items still exist (not deleted)
5. `otdfctl policy subject-mappings list` — both old and new SMs visible
6. `otdfctl migrate policy-graph` again — reports nothing new to migrate (idempotent)
7. `otdfctl migrate prune subject-mappings` — shows unnamespaced SMs with namespaced copies
8. `otdfctl migrate prune subject-mappings --commit` — deletes only matched unnamespaced SMs
9. `otdfctl migrate prune registered-resources --commit` — deletes only matched unnamespaced RRs
