# policy registered-resources

The `policy registered-resources` command provides CRUD operations for managing registered resources in the system. This command allows you to create, retrieve, update, delete, and list registered resources.

## Subcommands

### create

Create a new registered resource.

**Usage:**
```
otdfctl policy registered-resources create --name <name> [--label <label>...]
```

**Flags:**
- `--name` (required): The name of the registered resource.
- `--label`: Metadata labels to associate with the resource.

### get

Retrieve details of a specific registered resource by ID.

**Usage:**
```
otdfctl policy registered-resources get --id <id>
```

**Flags:**
- `--id` (required): The ID of the registered resource.

### list

List all registered resources with pagination.

**Usage:**
```
otdfctl policy registered-resources list --limit <limit> --offset <offset>
```

**Flags:**
- `--limit` (required): The number of resources to retrieve.
- `--offset` (required): The starting point for the list.

### update

Update an existing registered resource.

**Usage:**
```
otdfctl policy registered-resources update --id <id> [--name <name>] [--label <label>...]
```

**Flags:**
- `--id` (required): The ID of the registered resource to update.
- `--name`: The new name for the resource.
- `--label`: Metadata labels to update.

### delete

Delete a registered resource by ID.

**Usage:**
```
otdfctl policy registered-resources delete --id <id>
```

**Flags:**
- `--id` (required): The ID of the registered resource to delete.