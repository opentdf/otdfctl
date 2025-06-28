---
title: Create a new profile
command:
  name: create
  aliases:
    - add
  arguments:
    - name: profile
      description: Name of the profile to create
      required: true
      type: string
    - name: endpoint
      description: Platform endpoint URL
      required: true
      type: string
  flags:
    - name: set-default
      description: Set the profile as default
      default: false
      type: bool
    - name: tls-no-verify
      description: Disable TLS verification
      default: false
      type: bool
---

Create a new profile with the specified name and endpoint.

A profile stores connection settings for an OpenTDF Platform instance, including the endpoint URL and authentication credentials.