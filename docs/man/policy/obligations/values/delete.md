---
title: Delete an obligation
command:
  name: delete
  flags:
    - name: id
      shorthand: i
      description: ID of the obligation
      required: true
---

Because obligations are a post-entitlement decision, they are safe to delete. Upon deletion, any derived obligations as a result of
mappings to platform policy attribute values will be removed, and any TDFs containing obligations within will remain accessible, just
without obligations returned along with an access request to drive PEP behavior.