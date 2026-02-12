# Shell Key Management Integration Plan

## Current State

The shell currently supports:
- `/namespaces` - browse and manage namespaces, attributes, values
- `/keys` - browse existing KAS keys (read-only listing)
- `/namespaces/<ns>/key-assignments` - assign/remove keys from resources
- `registered-resources/` - placeholder

## Proposed Hierarchy

```
/
├── namespaces/                          # Existing
│   └── <namespace>/
│       ├── attribute-definitions/
│       │   └── <attribute>/
│       │       ├── attribute-values/
│       │       │   └── <value>/
│       │       │       └── key-assignments/
│       │       └── key-assignments/
│       └── key-assignments/
│
├── kas-registry/                        # NEW - Key Access Servers
│   └── <kas-name-or-id>/
│       └── keys/                        # Keys belonging to this KAS
│           └── <key-id>/
│
├── keys/                                # Existing - All keys (flat view)
│   └── <key-id>/                        # Key details + mappings
│
├── providers/                           # NEW - Key provider configs
│   └── <provider-name>/
│
└── registered-resources/                # Existing placeholder
```

## Implementation Phases

### Phase 1: KAS Registry Navigation

Add `/kas-registry` as a top-level directory.

**Path Types to Add:**
- `kas-registry` - list of all KAS servers
- `kas` - specific KAS server (shows keys/ subdir)

**Commands:**
| Location | Command | Action |
|----------|---------|--------|
| `/kas-registry` | `ls` | List all KAS servers |
| `/kas-registry` | `create` | Create new KAS (wizard) |
| `/kas-registry` | `cd <kas>` | Navigate into KAS |
| `/kas-registry/<kas>` | `ls` | Show `keys/` |
| `/kas-registry/<kas>` | `get` | Show KAS details |
| `/kas-registry/<kas>` | `rm` | Delete KAS |
| `/kas-registry/<kas>/keys` | `ls` | List keys for this KAS |
| `/kas-registry/<kas>/keys` | `create` | Create new key (wizard) |
| `/kas-registry/<kas>/keys/<key>` | `get` | Show key details |
| `/kas-registry/<kas>/keys/<key>` | `rm` | Delete key (unsafe, with warning) |

**New Wizard: KasRegistryWizard**
```
Create Key Access Server
─────────────────────────
URI: [https://kas.example.com________]
Name (optional): [my-kas________________]

Press Enter to continue, Esc to cancel
```

### Phase 2: Provider Configuration

Add `/providers` directory.

**Path Types:**
- `providers` - list of provider configs
- `provider` - specific provider config

**Commands:**
| Location | Command | Action |
|----------|---------|--------|
| `/providers` | `ls` | List all providers |
| `/providers` | `create` | Create provider (wizard) |
| `/providers` | `cd <provider>` | Navigate to provider |
| `/providers/<provider>` | `get` | Show provider details |
| `/providers/<provider>` | `rm` | Delete provider |

**New Wizard: ProviderWizard**
```
Create Key Provider Configuration
─────────────────────────────────
Name: [my-vault-provider_____________]
Manager: [openbao______________________]
Config (JSON): [{"url":"https://vault.example.com"}]

Press Enter to continue, Esc to cancel
```

### Phase 3: KAS Key Creation (Complex)

The key creation wizard needs to handle 4 different modes with different required fields.

**New Wizard: KasKeyWizard (Multi-step)**

```
Step 1: Select Key Mode
───────────────────────
▸ Local - KAS generates key, you provide wrapping key
  Provider - External provider stores wrapping key
  Remote - Key stored entirely at external provider
  Public Key Only - Import public key (no private key)

↑/↓ to select, Enter to continue, Esc to cancel
```

```
Step 2: Key Configuration
─────────────────────────
Key ID: [my-encryption-key___________]
Algorithm: ▸ rsa:2048
             rsa:4096
             ec:secp256r1
             ec:secp384r1
             ec:secp521r1

↑/↓ to select, Enter to continue, Esc to cancel
```

```
Step 3a (Local mode):
─────────────────────
Wrapping Key ID: [local-wrapping-key________]
Wrapping Key (hex): [a8c4824daafcfa38ed0d13...]

Press Enter to continue, Esc to cancel
```

```
Step 3b (Provider mode):
────────────────────────
Provider Config: ▸ my-vault-provider
                   aws-kms-config
Public Key PEM (base64): [LS0tLS1CRUdJTi...]
Private Key PEM (base64): [LS0tLS1CRUdJTi...]
Wrapping Key ID: [vault-key-id_____________]

Press Enter to continue, Esc to cancel
```

```
Step 3c (Remote mode):
──────────────────────
Provider Config: [select from list]
Public Key PEM (base64): [LS0tLS1CRUdJTi...]
Wrapping Key ID: [remote-key-id____________]

Press Enter to continue, Esc to cancel
```

```
Step 3d (Public Key Only mode):
───────────────────────────────
Public Key PEM (base64): [LS0tLS1CRUdJTi...]

Press Enter to continue, Esc to cancel
```

```
Step 4: Confirm
───────────────
Mode: Local
Key ID: my-encryption-key
Algorithm: rsa:2048
KAS: https://kas.example.com
Wrapping Key ID: local-wrapping-key

Press Enter to create, Esc to cancel
```

### Phase 4: Additional Key Operations

**From `/kas-registry/<kas>/keys/<key>` or `/keys/<key>`:**
- `rotate` - rotate the key (creates new key, marks old as rotated)
- `mappings` - show what namespaces/attributes/values use this key

**From root:**
- `base-key` - show current base key
- `set-base-key` - set a key as the base key

---

## File Changes Required

### tui/shell.go

1. **Add new path types** to `PathSegment.Type`:
   - `kas-registry`, `kas`, `providers`, `provider`

2. **Update `lsRoot()`**:
   ```go
   return "namespaces/\nkas-registry/\nkeys/\nproviders/\nregistered-resources/"
   ```

3. **Add navigation functions**:
   - `cdIntoKasRegistry(name)` - handle `cd kas-registry` and `cd <kas-name>`
   - `cdIntoKas(name)` - navigate from kas into `keys/`
   - `cdIntoProviders(name)` - handle `cd providers` and `cd <provider-name>`

4. **Add ls functions**:
   - `lsKasRegistry()` - list all KAS servers
   - `lsKas()` - show `keys/` directory
   - `lsKasKeysUnderKas()` - list keys for current KAS
   - `lsProviders()` - list all providers
   - `lsProvider()` - show provider details

5. **Update `startCreateWizard()`**:
   - Handle `kas-registry` context → `KasRegistryWizard`
   - Handle `kas` or `kas-keys` context → `KasKeyWizard`
   - Handle `providers` context → `ProviderWizard`

6. **Update `startDeleteWizard()`**:
   - Handle `kas` type → delete KAS (with confirmation)
   - Handle `provider` type → delete provider (with confirmation)
   - Handle key deletion → extra unsafe warning

7. **Update `getAvailableItems()`** for tab completion

### tui/wizard.go

1. **Add `KasRegistryWizard`** (~100 lines)
   - Simple 2-field wizard: uri, name
   - Uses existing `Wizard` pattern

2. **Add `ProviderWizard`** (~120 lines)
   - 3-field wizard: name, manager, config
   - Config is free-form JSON text

3. **Add `KasKeyWizard`** (~400 lines)
   - Multi-step wizard with mode selection
   - Dynamic fields based on mode
   - Provider selection dropdown (loads from API)
   - Algorithm selection dropdown

### pkg/handlers (if needed)

Check that all required handlers exist:
- `CreateKasRegistryEntry()` ✓
- `DeleteKasRegistryEntry()` ✓
- `CreateKasKey()` ✓
- `CreateProviderConfig()` ✓
- `DeleteProviderConfig()` ✓
- `ListProviderConfigs()` ✓

---

## Estimated Effort

| Component | Complexity | Est. Lines |
|-----------|------------|------------|
| KAS Registry navigation | Medium | ~150 |
| KasRegistryWizard | Low | ~100 |
| Providers navigation | Medium | ~100 |
| ProviderWizard | Low | ~120 |
| KasKeyWizard | High | ~400 |
| Delete handling updates | Low | ~50 |
| Tab completion updates | Low | ~30 |
| **Total** | | **~950** |

---

## Priority Order

1. **KAS Registry + KasRegistryWizard** - Foundation, enables key management
2. **Providers + ProviderWizard** - Needed for provider/remote key modes
3. **KasKeyWizard** - Most complex, but core functionality
4. **Key operations (rotate, mappings)** - Nice to have
5. **Base key management** - Nice to have

---

## Open Questions

1. **JSON input for provider config**:
   - Option A: Simple text input (user types/pastes JSON)
   - Option B: File path input
   - **Recommendation**: Text input, keep it simple

2. **PEM input for keys**:
   - Option A: Base64 text input (current CLI approach)
   - Option B: File path input
   - **Recommendation**: Base64 text input for consistency with CLI

3. **Should `/keys` remain?**
   - Yes - useful for seeing all keys across all KAS servers
   - `/kas-registry/<kas>/keys` shows keys for specific KAS

4. **Unsafe key delete**:
   - Support with extra warning ("THIS CANNOT BE UNDONE" + type confirmation)
