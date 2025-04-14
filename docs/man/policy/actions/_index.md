---
title: Manage Actions
command:
  name: actions
  aliases:
    - action
---

Actions are a set of `standard` and `custom` verbs at the core of an Access Decision or an
Obligation. Essentially, Actions answer what an Entity can _do_ to a Resource?

Standard Actions in Policy are comprise the below, and only their metadata labels are mutable:
- create
- read
- update
- delete

Custom Actions known to Policy are admin-defined.