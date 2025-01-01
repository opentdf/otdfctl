---
title: Prune (delete all un-mapped Subject Condition Sets)

command:
  name: prune
  flags:
    - name: force
      description: Force prune without interactive confirmation (dangerous)
---

This command will delete all Subject Condition Sets that are not utilized within any Subject Mappings and are therefore 'stranded'.

For more information about subject condition sets, see the `subject-condition-sets` subcommand.

## Example

```shell
otdfctl policy subject-condition-set prune
```

```shell
  SUCCESS                                                                                                                                                                                                                                                                                                                                                                                                                                             
╭──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╮
│ID                                                                                                                                                                                                                        │
├──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│5ecb1088-9c66-4fad-aa50-1d79fc84a344                                                                                                                                                                                      │
│c3167e9e-1987-4200-a45b-35127d86785c                                                                                                                                                                                      │
│66fe121a-0d14-48d0-aa33-59a5b1934fc6                                                                                                                                                                                      │
│524401e1-0ed0-4f70-924f-8978174e224b                                                                                                                                                                                      │                                                                                                                                                                                     │
╰──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╯
                                      
```
