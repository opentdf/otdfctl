---
title: Authenticate to the platform with the client-credentials flow

command:
  name: client-credentials
  args: 
    - client-id
  arbitrary_args:
    - client-secret
---

Allows the user to login in via Client Credentials flow. The client credentials will be stored safely
in the OS keyring for future use.
