---
title: Manage KAS registrations
command:
  name: kas-registry
  aliases:
    - kasRegistry
    - key-access-registry
---

# Manage Key Access Server registrations within the platform

The Key Access Server (KAS) registry is a record of servers granting and maintaining public keys.
The registry contains critical information like each server's uri, its public key (which can be
either local or at a remote uri), and any metadata about the server. Key Access Servers grant keys
for specified Attributes and their Values via Attribute Key Access Grants and Attribute Value Key
Access Grants.
