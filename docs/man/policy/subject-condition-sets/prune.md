---
title: Prune (delete all unmapped Subject Condition Sets)

command:
  name: prune
  flags:
    - name: force
      description: Force prune without interactive confirmation (dangerous)
---

This command will delete all Subject Condition Sets that are not utilized within any Subject Mappings and are therefore 'stranded'.

For more information about subject condition sets, see the `subject-condition-sets` subcommand.
