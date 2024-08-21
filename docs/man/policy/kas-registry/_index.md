---
title: Manage KAS registrations
command:
  name: kas-registry
  aliases:
    - kasr
    - kas-registries
---

The Key Access Server (KAS) registry is a record of KASs safeguarding access and maintaining public keys.
The registry contains critical information like each server's uri, its public key (which can be
either local or at a remote uri), and any metadata about the server. Key Access Servers grant keys
for specified Namespaces, Attributes, and their Values via Attribute Key Access Grants and Attribute Value Key
Access Grants.

For more information about grants and how KASs are utilized once registered, see the manual for the
`kas-grants` command.
