---
title: Set a profile value
command:
  name: set-endpoint
  arguments:
    - name: profile
      description: Name of the profile to update
      required: true
      type: string
    - name: endpoint
      description: New endpoint URL for the profile
      required: true
      type: string
  flags:
    - name: tls-no-verify
      description: Disable TLS verification
      default: false
      type: bool
---

Update the endpoint URL for an existing profile.

This allows changing the platform endpoint without recreating the profile and losing authentication credentials.