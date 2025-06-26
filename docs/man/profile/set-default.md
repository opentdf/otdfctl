---
title: Set a profile as default
command:
  name: set-default
  arguments:
    - name: profile
      description: Name of the profile to set as default
      required: true
      type: string
---

Set the specified profile as the default profile.

The default profile is used when no --profile flag is specified in commands.