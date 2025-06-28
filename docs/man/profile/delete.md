---
title: Delete a profile
command:
  name: delete
  arguments:
    - name: profile
      description: Name of the profile to delete
      required: true
      type: string
---

Delete a profile.

The default profile cannot be deleted. Set another profile as default before deleting the current default profile.