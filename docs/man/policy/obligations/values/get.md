---
title: Get an obligation
command:
  name: get
  flags:
    - name: id
      shorthand: i
      description: ID of the obligation
---

Retrieves the obligation, comprised of:

- ID
- value
- parent obligation (name and ID)
- FQN
- any assigned attribute value FQNs and IDs for derived obligations
- condition sets for mappings to entities that satisfy the obligation
