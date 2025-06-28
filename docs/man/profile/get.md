---
title: Get a profile value
command:
  name: get
  arguments:
    - name: profile
      description: Name of the profile to retrieve
      required: true
      type: string
---

Get detailed information about a specific profile.

Shows the profile name, endpoint, default status, and authentication type (with credentials masked for security).